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
	gocache "github.com/patrickmn/go-cache"
	"github.com/vyneer/pacany-bot/config"
	"github.com/vyneer/pacany-bot/db"
	"github.com/vyneer/pacany-bot/tg/commands/implementation"
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

	botCmdsMap := map[int]tgbotapi.BotCommand{}
	for _, v := range implementation.Interactable {
		desc, order := v.GetDescription()
		if order != -1 {
			botCmdsMap[order] = tgbotapi.BotCommand{
				Command:     v.GetParentName() + v.GetName(),
				Description: desc,
			}
		}
	}

	keys := make([]int, 0, len(botCmdsMap))
	for k := range botCmdsMap {
		keys = append(keys, k)
	}
	slices.Sort(keys)

	botCmdSlice := make([]tgbotapi.BotCommand, 0, len(botCmdsMap))
	for _, i := range keys {
		botCmdSlice = append(botCmdSlice, botCmdsMap[i])
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
		case update.Message != nil && (update.Message.Chat.IsGroup() || update.Message.Chat.IsSuperGroup()):
			chatID := update.Message.Chat.ID
			text := update.Message.Text
			username := update.Message.From.UserName

			var commandResponse implementation.CommandResponse

			if !update.Message.IsCommand() {
				commandResponse = implementation.GetAutomaticCommand("tagscan").Run(ctx, implementation.CommandArgs{
					DB:     b.tagDB,
					ChatID: chatID,
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

				if slices.ContainsFunc[[]tgbotapi.ChatMember](admins, func(cm tgbotapi.ChatMember) bool {
					return cm.User.ID == update.Message.From.ID
				}) {
					commandResponse = b.command(ctx, chatID, update.Message.Command(), update.Message.CommandArguments())
				}
			}

			response := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			response.Text = commandResponse.Text
			if commandResponse.Capitalize {
				response.Text = b.capitalize(response.Text)
			}
			if commandResponse.Reply {
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

func (b *Bot) command(ctx context.Context, chatID int64, command string, args string) implementation.CommandResponse {
	argsSplit := strings.Fields(args)

	cmd := implementation.GetInteractableCommand(command)
	if cmd == nil {
		return implementation.CommandResponse{
			Text:       "",
			Reply:      false,
			Capitalize: true,
		}
	}

	slog.Debug("running command", "chatID", chatID, "command", command)

	return cmd.Run(ctx, implementation.CommandArgs{
		DB:     b.tagDB,
		ChatID: chatID,
		Args:   argsSplit,
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
