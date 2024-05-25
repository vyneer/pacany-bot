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
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/shlex"
	gocache "github.com/patrickmn/go-cache"
	"github.com/vyneer/pacany-bot/config"
	"github.com/vyneer/pacany-bot/db"
	"github.com/vyneer/pacany-bot/tg/commands/implementation"
)

type Bot struct {
	api        *tgbotapi.BotAPI
	adminCache *cache.Cache[[]tgbotapi.ChatMember]
	db         *db.DB
}

func New(c *config.Config, tagDB *db.DB) (Bot, error) {
	_ = tgbotapi.SetLogger(log.Default())
	bot, err := tgbotapi.NewBotAPI(c.Token)
	if err != nil {
		return Bot{}, err
	}
	bot.Debug = c.Debug == 2

	slog.Debug("authorized on bot", "account", bot.Self.UserName)

	botCmdSlice := []tgbotapi.BotCommand{}
	for _, v := range implementation.GetInteractableOrder() {
		if desc, show := v.GetDescription(); show {
			botCmdSlice = append(botCmdSlice, tgbotapi.BotCommand{
				Command:     v.GetParentName() + v.GetName(),
				Description: desc,
			})
		}
	}

	if _, err := bot.Request(tgbotapi.NewSetMyCommandsWithScope(
		tgbotapi.NewBotCommandScopeAllChatAdministrators(),
		botCmdSlice...,
	)); err != nil {
		return Bot{}, err
	}

	gocacheClient := gocache.New(10*time.Minute, 15*time.Minute)
	gocacheStore := gocache_store.NewGoCache(gocacheClient)
	cacheManager := cache.New[[]tgbotapi.ChatMember](gocacheStore)

	return Bot{
		api:        bot,
		db:         tagDB,
		adminCache: cacheManager,
	}, nil
}

func (b *Bot) Run() error {
	u := tgbotapi.NewUpdate(-1)
	u.Timeout = 60
	u.AllowedUpdates = []string{"message", "chat_member"}
	updates := b.api.GetUpdatesChan(u)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	for update := range updates {
		switch {
		case update.ChatMember != nil:
			_, err := b.setAdmins(ctx, update.ChatMember.Chat.ChatConfig())
			if err != nil {
				slog.Warn("unable to update admin list", "err", err)
				continue
			}
		case update.Message != nil && (update.Message.Chat.IsGroup() || update.Message.Chat.IsSuperGroup()):
			chatID := update.Message.Chat.ID
			text := update.Message.Text
			username := update.Message.From.UserName

			var commandResponses []implementation.CommandResponse

			if !update.Message.IsCommand() {
				commandResponses = implementation.GetAutomaticCommand("tagscan").Run(ctx, implementation.CommandArgs{
					DB:     b.db,
					ChatID: chatID,
					User:   update.SentFrom(),
					Args: []string{
						username,
						text,
					},
				})
			} else {
				admins, err := b.getAdmins(ctx, update.Message.Chat.ChatConfig())
				if err != nil {
					slog.Warn("unable to get admin list", "err", err)
					continue
				}

				if slices.ContainsFunc(admins, func(cm tgbotapi.ChatMember) bool {
					return cm.User.ID == update.Message.From.ID
				}) {
					commandResponses = b.command(ctx, true, chatID, update.SentFrom(), update.Message.Command(), update.Message.CommandArguments())
				} else {
					commandResponses = b.command(ctx, false, chatID, update.SentFrom(), update.Message.Command(), update.Message.CommandArguments())
				}
			}

			for _, r := range commandResponses {
				response := tgbotapi.NewMessage(update.Message.Chat.ID, "")
				response.Text = r.Text
				if r.Capitalize {
					response.Text = b.capitalize(response.Text)
				}
				if r.Reply {
					response.ReplyToMessageID = update.Message.MessageID
				}

				if len(response.Text) != 0 {
					if _, err := b.api.Send(response); err != nil {
						log.Panic(err)
					}
				}
			}
		}
	}

	return nil
}

func (b *Bot) command(ctx context.Context, admin bool, chatID int64, user *tgbotapi.User, command string, args string) []implementation.CommandResponse {
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

func (b *Bot) getAdmins(ctx context.Context, c tgbotapi.ChatConfig) ([]tgbotapi.ChatMember, error) {
	slog.Debug("getting admin list", "chatID", c.ChatID)

	admins, err := b.adminCache.Get(ctx, c.ChatID)
	if err == nil {
		return admins, nil
	}

	admins, err = b.setAdmins(ctx, c)
	if err != nil {
		return nil, err
	}

	return admins, nil
}

func (b *Bot) setAdmins(ctx context.Context, c tgbotapi.ChatConfig) ([]tgbotapi.ChatMember, error) {
	slog.Debug("setting admin list", "chatID", c.ChatID)

	admins, err := b.api.GetChatAdministrators(tgbotapi.ChatAdministratorsConfig{
		ChatConfig: c,
	})
	if err != nil {
		return nil, err
	}

	if err := b.adminCache.Set(ctx, c.ChatID, admins); err != nil {
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
