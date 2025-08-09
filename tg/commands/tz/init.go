package tz

import (
	"github.com/vyneer/pacany-bot/config"
	"github.com/vyneer/pacany-bot/tg/commands/implementation"
	"github.com/vyneer/pacany-bot/tg/commands/tz/help"
	"github.com/vyneer/pacany-bot/tg/commands/tz/internal/list"
)

const (
	name        string = "tz"
	description string = "timezones"
)

type Parent struct{}

func NewTZ() *Parent {
	return &Parent{}
}

func (t *Parent) Name() string {
	return name
}

func (t *Parent) Description() string {
	return description
}

func (t *Parent) IsDisableable() bool {
	return true
}

func (t *Parent) Initialize(_ *config.Config) []implementation.Command {
	cmds := []implementation.Command{
		help.New(),
	}

	return append(cmds, list.Commands...)
}
