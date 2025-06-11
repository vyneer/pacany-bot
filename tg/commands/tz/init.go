package tz

import (
	"github.com/vyneer/pacany-bot/tg/commands/implementation"
	"github.com/vyneer/pacany-bot/tg/commands/tz/add"
	tzClear "github.com/vyneer/pacany-bot/tg/commands/tz/clear"
	"github.com/vyneer/pacany-bot/tg/commands/tz/convert"
	"github.com/vyneer/pacany-bot/tg/commands/tz/help"
	"github.com/vyneer/pacany-bot/tg/commands/tz/info"
	"github.com/vyneer/pacany-bot/tg/commands/tz/remove"
	"github.com/vyneer/pacany-bot/tg/commands/tz/set"
)

const name string = "tz"

type tz struct{}

func (t *tz) Name() string {
	return "tz"
}

func (t *tz) Description() string {
	return "timezones"
}

func (t *tz) Initialize() {
	implementation.EnableParentCommand(name)

	implementation.CreateInteractableCommand(help.New)
	implementation.CreateInteractableCommand(set.New)
	implementation.CreateInteractableCommand(tzClear.New)
	implementation.CreateInteractableCommand(info.New)
	implementation.CreateInteractableCommand(convert.New)
	implementation.CreateInteractableCommand(add.New)
	implementation.CreateInteractableCommand(remove.New)
}

func init() {
	implementation.CreateParentCommand(&tz{})
}
