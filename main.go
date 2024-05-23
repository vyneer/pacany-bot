package main

import (
	"log/slog"
	"os"

	"github.com/vyneer/pacany-bot/config"
	"github.com/vyneer/pacany-bot/db"
	"github.com/vyneer/pacany-bot/tg"
	_ "github.com/vyneer/pacany-bot/tg/commands/help"
	_ "github.com/vyneer/pacany-bot/tg/commands/tag"
	_ "github.com/vyneer/pacany-bot/tg/commands/tz"
)

func main() {
	lvl := &slog.LevelVar{}

	h := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		AddSource: true,
		Level:     lvl,
	})

	slog.SetDefault(slog.New(h))

	c, err := config.New()
	if err != nil {
		slog.Error("config error", "err", err)
		os.Exit(1)
	}

	if c.Debug > 0 {
		lvl.Set(slog.LevelDebug)
	}

	tagDB, err := db.New(&c)
	if err != nil {
		slog.Error("db setup error", "err", err)
		os.Exit(1)
	}

	bot, err := tg.New(&c, &tagDB)
	if err != nil {
		slog.Error("tg bot setup error", "err", err)
		os.Exit(1)
	}

	if err := bot.Run(); err != nil {
		slog.Error("tg bot run error", "err", err)
		os.Exit(1)
	}
}
