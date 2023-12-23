package main

import (
	"log/slog"
	"os"

	"github.com/vyneer/tg-tagbot/config"
	"github.com/vyneer/tg-tagbot/db"
	"github.com/vyneer/tg-tagbot/tg"
	_ "github.com/vyneer/tg-tagbot/tg/commands/help"
	_ "github.com/vyneer/tg-tagbot/tg/commands/tag"
)

func main() {
	lvl := &slog.LevelVar{}

	h := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
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
