package config

import (
	"errors"
	"log/slog"
	"os"
	"slices"
	"strings"

	"github.com/vyneer/pacany-bot/tg/commands/implementation"
)

var ErrNoToken = errors.New("no token provided")

var allowedSymbols = []rune{'@', '%', '#', '!', '&'}

type Config struct {
	Token                   string
	DBPath                  string
	Debug                   int
	Geonames                bool
	AllowedTagPrefixSymbols string

	commandList []implementation.Command
}

type ConfigurableCommand interface {
	Configure(cfg *Config)
}

func New() (Config, error) {
	c := Config{
		Geonames:                true,
		AllowedTagPrefixSymbols: "@%#!&",

		commandList: []implementation.Command{},
	}

	if t, ok := os.LookupEnv("TELEGRAM_TOKEN"); ok {
		c.Token = t
	} else {
		return c, ErrNoToken
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

	if gn, ok := os.LookupEnv("GEONAMES"); ok {
		if strings.ToLower(gn) == "false" || strings.ToLower(gn) == "0" {
			c.Geonames = false
		}
	}

	if ts, ok := os.LookupEnv("ALLOWED_TAG_PREFIX_SYMBOLS"); ok {
		tsRunes := []rune(ts)
		slices.Sort(tsRunes)
		tsRunes = slices.Compact(tsRunes)

		tsBuf := ""
		for _, r := range tsRunes {
			if slices.Contains(allowedSymbols, r) {
				tsBuf += string(r)
			}
		}

		if len(tsBuf) > 0 {
			c.AllowedTagPrefixSymbols = tsBuf
		}
	}

	cmdsSplit := []string{}
	if com, ok := os.LookupEnv("COMMANDS"); ok {
		cmdsSplit = strings.Split(com, ",")
	}

	for parentName, parentCommand := range implementation.GetAllParentCommands() {
		if slices.Contains(cmdsSplit, parentName) || len(cmdsSplit) == 0 || !parentCommand.IsDisableable() {
			c.commandList = append(c.commandList, parentCommand.Initialize()...)
			slog.Info("initialized command", "name", parentName)
			if configurable, ok := parentCommand.(ConfigurableCommand); ok {
				configurable.Configure(&c)
				slog.Info("configured command", "name", parentName)
			}
		}
	}

	return c, nil
}

func (c *Config) GetCommandList() []implementation.Command {
	return c.commandList
}
