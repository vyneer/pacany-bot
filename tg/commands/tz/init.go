package tz

import (
	"github.com/vyneer/pacany-bot/tg/commands/implementation"
	"github.com/vyneer/pacany-bot/tg/commands/tz/add"
	"github.com/vyneer/pacany-bot/tg/commands/tz/clear"
	"github.com/vyneer/pacany-bot/tg/commands/tz/convert"
	"github.com/vyneer/pacany-bot/tg/commands/tz/help"
	"github.com/vyneer/pacany-bot/tg/commands/tz/info"
	"github.com/vyneer/pacany-bot/tg/commands/tz/remove"
	"github.com/vyneer/pacany-bot/tg/commands/tz/set"
)

const name string = "tz"

func init() {
	implementation.CreateParentCommand(implementation.ParentCommand{
		Name:        name,
		Description: "timezones",
		Initialize:  initialize,
	})
}

func initialize() {
	implementation.EnableParentCommand(name)

	implementation.CreateInteractableCommand(help.New)
	implementation.CreateInteractableCommand(set.New)
	implementation.CreateInteractableCommand(clear.New)
	implementation.CreateInteractableCommand(info.New)
	implementation.CreateInteractableCommand(convert.New)
	implementation.CreateInteractableCommand(add.New)
	implementation.CreateInteractableCommand(remove.New)
}
