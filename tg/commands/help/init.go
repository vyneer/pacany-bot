package help

import (
	"github.com/vyneer/tg-tagbot/tg/commands/implementation"
)

func init() {
	implementation.CreateInteractableCommand(New)
}
