package tz

import (
	"github.com/vyneer/pacany-bot/tg/commands/implementation"
	"github.com/vyneer/pacany-bot/tg/commands/tz/help"
	"github.com/vyneer/pacany-bot/tg/commands/tz/internal/list"
)

const name string = "tz"

type tz struct{}

func (t *tz) Name() string {
	return "tz"
}

func (t *tz) Description() string {
	return "timezones"
}

func (t *tz) IsDisableable() bool {
	return true
}

func (t *tz) Initialize() []implementation.Command {
	implementation.EnableParentCommand(name)

	cmds := []implementation.Command{
		help.New(),
	}

	return append(cmds, list.Commands...)
}

func init() {
	implementation.CreateParentCommand(&tz{})
}
