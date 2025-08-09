package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"

	"github.com/vyneer/pacany-bot/config"
	"github.com/vyneer/pacany-bot/db"
	"github.com/vyneer/pacany-bot/geonames"
	"github.com/vyneer/pacany-bot/tg"
	"github.com/vyneer/pacany-bot/tg/commands/help"
	"github.com/vyneer/pacany-bot/tg/commands/tag"
	"github.com/vyneer/pacany-bot/tg/commands/tz"
)

func main() {
	lvl := &slog.LevelVar{}

	h := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		AddSource: true,
		Level:     lvl,
	})

	slog.SetDefault(slog.New(h))

	c, err := config.New([]config.ParentCommand{
		help.NewHelp(),
		tag.NewTag(),
		tz.NewTZ(),
	})
	if err != nil {
		slog.Error("config error", "err", err)
		os.Exit(1)
	}

	if c.Debug > 0 {
		lvl.Set(slog.LevelDebug)
	}

	if c.Geonames {
		if err := geonames.New(); err != nil {
			slog.Error("geonames error", "err", err)
			os.Exit(1)
		}
	}

	tagDB, err := db.New(c.DBPath)
	if err != nil {
		slog.Error("db setup error", "err", err)
		os.Exit(1)
	}

	bot, err := tg.New(&c, &tagDB)
	if err != nil {
		slog.Error("tg bot setup error", "err", err)
		os.Exit(1)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	if err := bot.RegisterCommands(c.GetCommandList()); err != nil {
		slog.Error("tg bot command registration error", "err", err)
		return
	}

	bot.Start(ctx)
}
