package tz

import (
	"github.com/vyneer/pacany-bot/config"
	"github.com/vyneer/pacany-bot/tg/commands/implementation"
	"github.com/vyneer/pacany-bot/tg/commands/tz/add"
	tzClear "github.com/vyneer/pacany-bot/tg/commands/tz/clear"
	"github.com/vyneer/pacany-bot/tg/commands/tz/convert"
	"github.com/vyneer/pacany-bot/tg/commands/tz/help"
	"github.com/vyneer/pacany-bot/tg/commands/tz/info"
	"github.com/vyneer/pacany-bot/tg/commands/tz/remove"
	"github.com/vyneer/pacany-bot/tg/commands/tz/set"
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
		set.New(),
		tzClear.New(),
		info.New(),
		convert.New(),
		add.New(),
		remove.New(),
	}

	return append([]implementation.Command{help.New(cmds)}, cmds...)
}
