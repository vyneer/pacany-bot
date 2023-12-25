package help

import (
	"github.com/vyneer/pacany-bot/tg/commands/implementation"
)

func init() {
	implementation.CreateInteractableCommand(New)
}
