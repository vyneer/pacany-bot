package help

import (
	"github.com/vyneer/pacani-bot/tg/commands/implementation"
)

func init() {
	implementation.CreateInteractableCommand(New)
}
