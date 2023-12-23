package tag

import (
	"github.com/vyneer/tg-tagbot/tg/commands/implementation"
	"github.com/vyneer/tg-tagbot/tg/commands/tag/add"
	"github.com/vyneer/tg-tagbot/tg/commands/tag/adduser"
	"github.com/vyneer/tg-tagbot/tg/commands/tag/info"
	"github.com/vyneer/tg-tagbot/tg/commands/tag/remove"
	"github.com/vyneer/tg-tagbot/tg/commands/tag/removeuser"
	"github.com/vyneer/tg-tagbot/tg/commands/tag/rename"
	"github.com/vyneer/tg-tagbot/tg/commands/tag/scan"
)

func init() {
	implementation.CreateInteractableCommand(add.New)
	implementation.CreateInteractableCommand(remove.New)
	implementation.CreateInteractableCommand(rename.New)
	implementation.CreateInteractableCommand(adduser.New)
	implementation.CreateInteractableCommand(removeuser.New)
	implementation.CreateInteractableCommand(info.New)
	implementation.CreateAutomaticCommand(scan.New)
}
