package help

import (
	"github.com/vyneer/tg-tagbot/tg/commands/help/tag"
	"github.com/vyneer/tg-tagbot/tg/commands/implementation"
)

func init() {
	implementation.CreateInteractableCommand(New)
	implementation.CreateInteractableCommand(tag.New)
}
