package tg

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"regexp"
	"slices"
	"strings"
	"time"
	"unicode"

	"github.com/eko/gocache/lib/v4/cache"
	gocache_store "github.com/eko/gocache/store/go_cache/v4"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	gocache "github.com/patrickmn/go-cache"
	"github.com/vyneer/tg-tagbot/config"
	"github.com/vyneer/tg-tagbot/db"
)

var (
	prefixRegex = regexp.MustCompile(`(?i)^[@%#!&]{1}`)

	ErrInvalidTag   = errors.New("invalid tag: only @%#!& can be used")
	ErrNoValidUsers = errors.New("no valid users provided")
)

type Bot struct {
	api        *tgbotapi.BotAPI
	adminCache *cache.Cache[[]tgbotapi.ChatMember]
	tagDB      *db.DB
}

func New(c *config.Config, tagDB *db.DB) (Bot, error) {
	_ = tgbotapi.SetLogger(log.Default())
	bot, err := tgbotapi.NewBotAPI(c.Token)
	if err != nil {
		return Bot{}, err
	}
	bot.Debug = c.Debug == 2

	slog.Debug("authorized on bot", "account", bot.Self.UserName)

	if _, err := bot.Request(tgbotapi.NewSetMyCommandsWithScope(
		tgbotapi.NewBotCommandScopeAllChatAdministrators(),
		tgbotapi.BotCommand{
			Command:     "help",
			Description: "Prints help",
		},
		tgbotapi.BotCommand{
			Command:     "tag",
			Description: "Manage tags",
		},
	)); err != nil {
		return Bot{}, err
	}

	gocacheClient := gocache.New(10*time.Minute, 15*time.Minute)
	gocacheStore := gocache_store.NewGoCache(gocacheClient)
	cacheManager := cache.New[[]tgbotapi.ChatMember](gocacheStore)

	return Bot{
		api:        bot,
		tagDB:      tagDB,
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
		case update.Message != nil && update.Message.Chat.IsGroup():
			chatID := update.Message.Chat.ID
			text := update.Message.Text
			username := update.Message.From.UserName

			response := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			var reply bool

			if !update.Message.IsCommand() {
				response.Text, reply = b.Scan(ctx, chatID, username, text)
			} else {
				admins, err := b.getAdmins(ctx, update.Message.Chat.ChatConfig())
				if err != nil {
					slog.Warn("unable to get admin list", "err", err)
					continue
				}

				if slices.ContainsFunc[[]tgbotapi.ChatMember](admins, func(cm tgbotapi.ChatMember) bool {
					return cm.User.ID == update.Message.From.ID
				}) {
					response.Text, reply = b.command(ctx, chatID, update.Message.Command(), update.Message.CommandArguments())
				}
			}

			response.Text = b.capitalize(response.Text)
			if reply {
				response.ReplyToMessageID = update.Message.MessageID
			}

			if len(response.Text) != 0 {
				if _, err := b.api.Send(response); err != nil {
					log.Panic(err)
				}
			}
		}
	}

	return nil
}

func (b *Bot) command(ctx context.Context, chatID int64, command string, args string) (string, bool) {
	slog.Debug("running command", "chatID", chatID, "command", command)

	switch command {
	case "help":
		return b.Help()
	case "tag":
		subcommandSplit := strings.Fields(args)
		if len(subcommandSplit) < 1 || subcommandSplit[0] == "" {
			return b.SmallHelp()
		}
		subcommand := subcommandSplit[0]

		slog.Debug("running subcommand", "chatID", chatID, "command", command, "subcommand", subcommand, "args", args)

		switch subcommand {
		case "new":
			return b.NewTag(ctx, chatID, subcommandSplit...)
		case "remove":
			return b.RemoveTag(ctx, chatID, subcommandSplit...)
		case "add-user":
			return b.AddUsers(ctx, chatID, subcommandSplit...)
		case "remove-user":
			return b.RemoveUsers(ctx, chatID, subcommandSplit...)
		case "info":
			return b.Info(ctx, chatID, subcommandSplit...)
		case "list":
			return b.List(ctx, chatID)
		}
	}

	return "", false
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

func (b *Bot) isValidUserName(username string) bool {
	return strings.HasPrefix(username, "@") && len(username) > 1
}

func (b *Bot) filterInvalidUsernames(usernames []string) []string {
	var validUsernames []string

	for _, v := range usernames {
		if b.isValidUserName(v) {
			validUsernames = append(validUsernames, v)
		}
	}

	return validUsernames
}

func (b *Bot) filterMentions(mentions string, ignore string) (string, bool) {
	var filteredMentions []string

	for _, v := range strings.Fields(mentions) {
		if strings.TrimPrefix(v, "@") != ignore {
			filteredMentions = append(filteredMentions, v)
		}
	}

	return strings.Join(filteredMentions, " "), len(filteredMentions) > 0
}

func (b *Bot) isValidTagName(name string) bool {
	return prefixRegex.MatchString(name) && len(name) > 1
}

func (b *Bot) capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	r := []rune(s)
	return string(append([]rune{unicode.ToUpper(r[0])}, r[1:]...))
}
