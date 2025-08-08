package tg

import (
	"context"
	"log"
	"log/slog"
	"slices"
	"strings"
	"time"
	"unicode"

	"github.com/eko/gocache/lib/v4/cache"
	gocache_store "github.com/eko/gocache/store/go_cache/v4"
	tgbotapi "github.com/go-telegram/bot"
	tgbotapiModels "github.com/go-telegram/bot/models"
	"github.com/google/shlex"
	gocache "github.com/patrickmn/go-cache"
	"github.com/vyneer/pacany-bot/config"
	"github.com/vyneer/pacany-bot/db"
	"github.com/vyneer/pacany-bot/tg/commands/implementation"
)

type Bot struct {
	API        *tgbotapi.Bot
	adminCache *cache.Cache[[]tgbotapiModels.ChatMember]
	db         *db.DB
}

func New(c *config.Config, tagDB *db.DB) (*Bot, error) {
	gocacheClient := gocache.New(10*time.Minute, 15*time.Minute)
	gocacheStore := gocache_store.NewGoCache(gocacheClient)
	cacheManager := cache.New[[]tgbotapiModels.ChatMember](gocacheStore)

	b := Bot{
		db:         tagDB,
		adminCache: cacheManager,
	}

	bot, err := tgbotapi.New(c.Token, tgbotapi.WithAllowedUpdates(tgbotapi.AllowedUpdates{
		"message", "chat_member",
	}), tgbotapi.WithDefaultHandler(handle(&b)))
	if err != nil {
		return nil, err
	}
	// bot.Debug = c.Debug == 2
	botMe, err := bot.GetMe(context.Background())
	if err != nil {
		return nil, err
	}

	slog.Debug("authorized on bot", "account", botMe.Username)

	botCmdSlice := []tgbotapiModels.BotCommand{}
	botCmdAdminSlice := []tgbotapiModels.BotCommand{}
	for _, v := range implementation.GetInteractableOrder() {
		if desc, show := v.GetDescription(); show {
			if !v.IsAdminOnly() {
				botCmdSlice = append(botCmdSlice, tgbotapiModels.BotCommand{
					Command:     v.GetParentName() + v.GetName(),
					Description: desc,
				})
			}

			botCmdAdminSlice = append(botCmdAdminSlice, tgbotapiModels.BotCommand{
				Command:     v.GetParentName() + v.GetName(),
				Description: desc,
			})
		}
	}

	if _, err := bot.SetMyCommands(context.Background(), &tgbotapi.SetMyCommandsParams{
		Commands: botCmdSlice,
		Scope:    &tgbotapiModels.BotCommandScopeAllGroupChats{},
	}); err != nil {
		return nil, err
	}

	if _, err := bot.SetMyCommands(context.Background(), &tgbotapi.SetMyCommandsParams{
		Commands: botCmdAdminSlice,
		Scope:    &tgbotapiModels.BotCommandScopeAllChatAdministrators{},
	}); err != nil {
		return nil, err
	}

	b.API = bot

	return &b, nil
}

func handle(b *Bot) tgbotapi.HandlerFunc {
	return func(ctx context.Context, bot *tgbotapi.Bot, update *tgbotapiModels.Update) {
		switch {
		case update.ChatMember != nil:
			_, err := b.setAdmins(ctx, update.ChatMember.Chat)
			if err != nil {
				slog.Warn("unable to update admin list", "err", err)
				return
			}
		case update.Message != nil && (update.Message.Chat.Type == "group" || update.Message.Chat.Type == "supergroup"):
			chatID := update.Message.Chat.ID
			text := update.Message.Text
			username := update.Message.From.Username

			var commandResponses []implementation.CommandResponse

			if !isCommand(update.Message) {
				commandResponses = implementation.GetAutomaticCommand("tagscan").Run(ctx, implementation.CommandArgs{
					DB:     b.db,
					ChatID: chatID,
					User:   update.Message.From,
					Args: []string{
						username,
						text,
					},
				})
			} else {
				admins, err := b.getAdmins(ctx, update.Message.Chat)
				if err != nil {
					slog.Warn("unable to get admin list", "err", err)
					return
				}

				if slices.ContainsFunc(admins, func(cm tgbotapiModels.ChatMember) bool {
					return cm.Member.User.ID == update.Message.From.ID
				}) {
					commandResponses = b.command(ctx, true, chatID, update.Message.From, command(update.Message), commandArguments(update.Message))
				} else {
					commandResponses = b.command(ctx, false, chatID, update.Message.From, command(update.Message), commandArguments(update.Message))
				}
			}

			for _, r := range commandResponses {
				response := &tgbotapi.SendMessageParams{}
				response.ChatID = update.Message.Chat.ID
				response.Text = r.Text
				if update.Message.Chat.IsForum {
					response.MessageThreadID = update.Message.MessageThreadID
				}
				if r.Capitalize {
					response.Text = b.capitalize(response.Text)
				}
				if r.Reply {
					response.ReplyParameters = &tgbotapiModels.ReplyParameters{
						ChatID:    update.Message.Chat.ID,
						MessageID: update.Message.ID,
					}
				}

				if len(response.Text) != 0 {
					if _, err := bot.SendMessage(ctx, response); err != nil {
						log.Panic(err)
					}
				}
			}
		}
	}
}

func (b *Bot) command(ctx context.Context, admin bool, chatID int64, user *tgbotapiModels.User, command string, args string) []implementation.CommandResponse {
	argsSplit, err := shlex.Split(args)
	if err != nil {
		slog.Warn("unable to shlex args string", "err", err)
		argsSplit = strings.Fields(args)
	}

	cmd := implementation.GetInteractableCommand(command)
	if cmd == nil {
		return []implementation.CommandResponse{}
	}

	if cmd.IsAdminOnly() && !admin {
		return []implementation.CommandResponse{}
	}

	slog.Debug("running command", "chatID", chatID, "command", command)

	return cmd.Run(ctx, implementation.CommandArgs{
		DB:      b.db,
		ChatID:  chatID,
		User:    user,
		IsAdmin: admin,
		Args:    argsSplit,
	})
}

func (b *Bot) getAdmins(ctx context.Context, c tgbotapiModels.Chat) ([]tgbotapiModels.ChatMember, error) {
	slog.Debug("getting admin list", "chatID", c.ID)

	admins, err := b.adminCache.Get(ctx, c.ID)
	if err == nil {
		return admins, nil
	}

	admins, err = b.setAdmins(ctx, c)
	if err != nil {
		return nil, err
	}

	return admins, nil
}

func (b *Bot) setAdmins(ctx context.Context, c tgbotapiModels.Chat) ([]tgbotapiModels.ChatMember, error) {
	slog.Debug("setting admin list", "chatID", c.ID)

	admins, err := b.API.GetChatAdministrators(ctx, &tgbotapi.GetChatAdministratorsParams{
		ChatID: c.ID,
	})
	if err != nil {
		return nil, err
	}

	if err := b.adminCache.Set(ctx, c.ID, admins); err != nil {
		return nil, err
	}

	return admins, nil
}

func (b *Bot) capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	r := []rune(s)
	return string(append([]rune{unicode.ToUpper(r[0])}, r[1:]...))
}

func isCommand(m *tgbotapiModels.Message) bool {
	if len(m.Entities) == 0 {
		return false
	}

	entity := m.Entities[0]
	return entity.Offset == 0 && entity.Type == "bot_command"
}

func command(m *tgbotapiModels.Message) string {
	if !isCommand(m) {
		return ""
	}

	entity := m.Entities[0]
	command := m.Text[1:entity.Length]

	if i := strings.Index(command, "@"); i != -1 {
		command = command[:i]
	}

	return command
}

func commandArguments(m *tgbotapiModels.Message) string {
	if !isCommand(m) {
		return ""
	}

	entity := m.Entities[0]

	if len(m.Text) == entity.Length {
		return ""
	}

	return m.Text[entity.Length+1:]
}
