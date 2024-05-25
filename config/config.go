package config

import (
	"errors"
	"log/slog"
	"os"
	"strings"

	"github.com/vyneer/pacany-bot/tg/commands/implementation"
)

var ErrNoToken = errors.New("no token provided")

type Config struct {
	Token  string
	DBPath string
	Debug  int
}

func New() (Config, error) {
	c := Config{}

	if t, ok := os.LookupEnv("TELEGRAM_TOKEN"); ok {
		c.Token = t
	} else {
		return c, ErrNoToken
	}

	if com, ok := os.LookupEnv("COMMANDS"); ok {
		split := strings.Split(com, ",")
		for _, parentName := range split {
			if parentCommand, ok := implementation.GetParentCommand(parentName); ok {
				parentCommand.Initialize()
				slog.Info("initialized command", "name", parentName)
			}
		}
	} else {
		for parentName, parentCommand := range implementation.GetAllParentCommands() {
			parentCommand.Initialize()
			slog.Info("initialized command", "name", parentName)
		}
	}

	if p, ok := os.LookupEnv("DB_PATH"); ok {
		c.DBPath = p
	} else {
		c.DBPath = "pacany-bot.sqlite"
	}

	if d, ok := os.LookupEnv("DEBUG"); ok {
		switch strings.ToLower(d) {
		case "debug", "1":
			c.Debug = 1
		case "trace", "2":
			c.Debug = 2
		}
	}

	return c, nil
}
