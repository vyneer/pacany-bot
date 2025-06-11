package tag

import (
	"github.com/vyneer/pacany-bot/tg/commands/implementation"
	"github.com/vyneer/pacany-bot/tg/commands/tag/add"
	"github.com/vyneer/pacany-bot/tg/commands/tag/del"
	"github.com/vyneer/pacany-bot/tg/commands/tag/desc"
	"github.com/vyneer/pacany-bot/tg/commands/tag/descdel"
	"github.com/vyneer/pacany-bot/tg/commands/tag/help"
	"github.com/vyneer/pacany-bot/tg/commands/tag/info"
	"github.com/vyneer/pacany-bot/tg/commands/tag/kick"
	tagname "github.com/vyneer/pacany-bot/tg/commands/tag/name"
	tagnew "github.com/vyneer/pacany-bot/tg/commands/tag/new"
	"github.com/vyneer/pacany-bot/tg/commands/tag/scan"
)

const name string = "tag"

func init() {
	implementation.CreateParentCommand(implementation.ParentCommand{
		Name:        name,
		Description: "tagging",
		Initialize:  initialize,
	})
}

func initialize() {
	implementation.EnableParentCommand(name)

	implementation.CreateInteractableCommand(help.New)
	implementation.CreateInteractableCommand(tagnew.New)
	implementation.CreateInteractableCommand(del.New)
	implementation.CreateInteractableCommand(tagname.New)
	implementation.CreateInteractableCommand(descdel.New)
	implementation.CreateInteractableCommand(desc.New)
	implementation.CreateInteractableCommand(add.New)
	implementation.CreateInteractableCommand(kick.New)
	implementation.CreateInteractableCommand(info.New)
	implementation.CreateAutomaticCommand(scan.New)
}
