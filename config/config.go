package config

import (
	"errors"
	"os"
	"strings"
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

	if p, ok := os.LookupEnv("DB_PATH"); ok {
		c.DBPath = p
	} else {
		c.DBPath = "pacani-bot.sqlite"
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
