package tag

import (
	"github.com/vyneer/pacani-bot/tg/commands/implementation"
	"github.com/vyneer/pacani-bot/tg/commands/tag/add"
	"github.com/vyneer/pacani-bot/tg/commands/tag/adduser"
	"github.com/vyneer/pacani-bot/tg/commands/tag/help"
	"github.com/vyneer/pacani-bot/tg/commands/tag/info"
	"github.com/vyneer/pacani-bot/tg/commands/tag/remove"
	"github.com/vyneer/pacani-bot/tg/commands/tag/removeuser"
	"github.com/vyneer/pacani-bot/tg/commands/tag/rename"
	"github.com/vyneer/pacani-bot/tg/commands/tag/scan"
)

func init() {
	implementation.CreateInteractableCommand(help.New)
	implementation.CreateInteractableCommand(add.New)
	implementation.CreateInteractableCommand(remove.New)
	implementation.CreateInteractableCommand(rename.New)
	implementation.CreateInteractableCommand(adduser.New)
	implementation.CreateInteractableCommand(removeuser.New)
	implementation.CreateInteractableCommand(info.New)
	implementation.CreateAutomaticCommand(scan.New)
}
