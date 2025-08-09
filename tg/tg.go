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
	bot         *tgbotapi.Bot
	botUsername string
	adminCache  *cache.Cache[[]tgbotapiModels.ChatMember]
	db          *db.DB
}

func New(c *config.Config, tagDB *db.DB) (*Bot, error) {
	gocacheClient := gocache.New(10*time.Minute, 15*time.Minute)
	gocacheStore := gocache_store.NewGoCache(gocacheClient)
	cacheManager := cache.New[[]tgbotapiModels.ChatMember](gocacheStore)

	newBot, err := tgbotapi.New(c.Token, tgbotapi.WithAllowedUpdates(tgbotapi.AllowedUpdates{
		"message", "chat_member",
	}))
	if err != nil {
		return nil, err
	}

	botMe, err := newBot.GetMe(context.Background())
	if err != nil {
		return nil, err
	}

	slog.Debug("authorized on bot", "account", botMe.Username)

	b := Bot{
		bot:         newBot,
		botUsername: botMe.Username,
		adminCache:  cacheManager,
		db:          tagDB,
	}

	return &b, nil
}

func (b *Bot) RegisterCommands(commands []implementation.Command) error {
	tgbotapi.WithDefaultHandler(func(ctx context.Context, bot *tgbotapi.Bot, update *tgbotapiModels.Update) {
		// this doesnt actually do what i expected it to do
		// i wanted updates on member status changes (as in, promoted to admin from member, etc)
		// but this only tracks when users get added to chat or get removed from chat as far as i can tell
		// i'll keep it for now but idk if it's necessary
		if update.ChatMember != nil {
			_, err := b.setAdmins(ctx, bot, update.ChatMember.Chat)
			if err != nil {
				slog.Warn("unable to update admin list", "err", err)
				return
			}
		}
	})(b.bot)

	botCmdSlice := []tgbotapiModels.BotCommand{}
	botCmdAdminSlice := []tgbotapiModels.BotCommand{}

	for _, command := range commands {
		switch v := command.(type) {
		case implementation.InteractableCommand:
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

			shortCommandPrefix := v.GetParentName() + v.GetName()
			longCommandPrefix := v.GetParentName() + v.GetName() + "@" + b.botUsername

			b.bot.RegisterHandlerMatchFunc(func(update *tgbotapiModels.Update) bool {
				if len(update.Message.Entities) == 0 {
					return false
				}

				firstEntity := update.Message.Entities[0]

				if firstEntity.Type == tgbotapiModels.MessageEntityTypeBotCommand {
					if firstEntity.Offset != 0 {
						return false
					}

					cmdString := update.Message.Text[firstEntity.Offset+1 : firstEntity.Offset+firstEntity.Length]
					if cmdString == shortCommandPrefix || cmdString == longCommandPrefix {
						return true
					}
				}

				return false
			}, b.interactableCommandHandler(v))
		case implementation.AutomaticCommand:
			b.bot.RegisterHandlerRegexp(tgbotapi.HandlerTypeMessageText, v.GetMatcher(), b.automaticCommandHandler(v))
		}
	}

	if _, err := b.bot.SetMyCommands(context.Background(), &tgbotapi.SetMyCommandsParams{
		Commands: botCmdSlice,
		Scope:    &tgbotapiModels.BotCommandScopeAllGroupChats{},
	}); err != nil {
		return err
	}

	if _, err := b.bot.SetMyCommands(context.Background(), &tgbotapi.SetMyCommandsParams{
		Commands: botCmdAdminSlice,
		Scope:    &tgbotapiModels.BotCommandScopeAllChatAdministrators{},
	}); err != nil {
		return err
	}

	return nil
}

func (b *Bot) Start(ctx context.Context) {
	b.bot.Start(ctx)
}

func (b *Bot) getAdmins(ctx context.Context, bot *tgbotapi.Bot, c tgbotapiModels.Chat) ([]tgbotapiModels.ChatMember, error) {
	slog.Debug("getting admin list", "chatID", c.ID)

	admins, err := b.adminCache.Get(ctx, c.ID)
	if err == nil {
		return admins, nil
	}

	admins, err = b.setAdmins(ctx, bot, c)
	if err != nil {
		return nil, err
	}

	return admins, nil
}

func (b *Bot) setAdmins(ctx context.Context, bot *tgbotapi.Bot, c tgbotapiModels.Chat) ([]tgbotapiModels.ChatMember, error) {
	slog.Debug("setting admin list", "chatID", c.ID)

	admins, err := bot.GetChatAdministrators(ctx, &tgbotapi.GetChatAdministratorsParams{
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

func (b *Bot) interactableCommandHandler(command implementation.InteractableCommand) tgbotapi.HandlerFunc {
	return func(ctx context.Context, bot *tgbotapi.Bot, update *tgbotapiModels.Update) {
		if isGroupChat(update) {
			chatID := update.Message.Chat.ID

			// this feels ugly, maybe worth making admin cache a completely separate structure
			admins, err := b.getAdmins(ctx, bot, update.Message.Chat)
			if err != nil {
				slog.Warn("unable to get admin list", "err", err)
				return
			}

			isAdmin := slices.ContainsFunc(admins, func(cm tgbotapiModels.ChatMember) bool {
				return cm.Member.User.ID == chatID
			})

			if command.IsAdminOnly() && !isAdmin {
				return
			}

			commandArgs := commandArguments(update.Message)

			argsSplit, err := shlex.Split(commandArgs)
			if err != nil {
				slog.Warn("unable to shlex args string", "err", err)
				argsSplit = strings.Fields(commandArgs)
			}

			slog.Debug("running command", "chatID", chatID, "command", command.GetParentName()+command.GetName())

			commandResponses := command.Run(ctx, implementation.CommandArgs{
				DB:      b.db,
				ChatID:  chatID,
				User:    update.Message.From,
				IsAdmin: isAdmin,
				Args:    argsSplit,
			})

			msgsToSend := transformCommandResponsesIntoMessages(update, commandResponses)
			for _, m := range msgsToSend {
				if _, err := bot.SendMessage(ctx, &m); err != nil {
					log.Panic(err)
				}
			}
		}
	}
}

func (b *Bot) automaticCommandHandler(command implementation.AutomaticCommand) tgbotapi.HandlerFunc {
	return func(ctx context.Context, bot *tgbotapi.Bot, update *tgbotapiModels.Update) {
		if isGroupChat(update) {
			slog.Debug("running automatic command", "chatID", update.Message.Chat.ID, "command", command.GetIdentifier())

			commandResponses := command.Run(ctx, implementation.CommandArgs{
				DB:     b.db,
				ChatID: update.Message.Chat.ID,
				User:   update.Message.From,
				Args: []string{
					update.Message.From.Username,
					update.Message.Text,
				},
			})

			msgsToSend := transformCommandResponsesIntoMessages(update, commandResponses)
			for _, m := range msgsToSend {
				if _, err := bot.SendMessage(ctx, &m); err != nil {
					log.Panic(err)
				}
			}
		}
	}
}

func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	r := []rune(s)
	return string(append([]rune{unicode.ToUpper(r[0])}, r[1:]...))
}

func transformCommandResponsesIntoMessages(update *tgbotapiModels.Update, resps []implementation.CommandResponse) []tgbotapi.SendMessageParams {
	paramsSlice := []tgbotapi.SendMessageParams{}

	for _, r := range resps {
		response := tgbotapi.SendMessageParams{}
		response.ChatID = update.Message.Chat.ID
		response.Text = r.Text
		if update.Message.Chat.IsForum {
			response.MessageThreadID = update.Message.MessageThreadID
		}
		if r.Capitalize {
			response.Text = capitalize(response.Text)
		}
		if r.Reply {
			response.ReplyParameters = &tgbotapiModels.ReplyParameters{
				ChatID:    update.Message.Chat.ID,
				MessageID: update.Message.ID,
			}
		}

		if len(response.Text) != 0 {
			paramsSlice = append(paramsSlice, response)
		}
	}

	return paramsSlice
}

func commandArguments(m *tgbotapiModels.Message) string {
	entity := m.Entities[0]

	if len(m.Text) == entity.Length {
		return ""
	}

	return m.Text[entity.Length+1:]
}

func isGroupChat(update *tgbotapiModels.Update) bool {
	return update.Message != nil && (update.Message.Chat.Type == "group" || update.Message.Chat.Type == "supergroup")
}
