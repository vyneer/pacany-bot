package tz

import (
	"github.com/vyneer/pacany-bot/tg/commands/implementation"
	"github.com/vyneer/pacany-bot/tg/commands/tz/clear"
	"github.com/vyneer/pacany-bot/tg/commands/tz/convert"
	"github.com/vyneer/pacany-bot/tg/commands/tz/help"
	"github.com/vyneer/pacany-bot/tg/commands/tz/info"
	"github.com/vyneer/pacany-bot/tg/commands/tz/set"
)

func init() {
	implementation.CreateInteractableCommand(help.New)
	implementation.CreateInteractableCommand(set.New)
	implementation.CreateInteractableCommand(clear.New)
	implementation.CreateInteractableCommand(info.New)
	implementation.CreateInteractableCommand(convert.New)
}
