package help

import (
	"github.com/vyneer/pacany-bot/tg/commands/implementation"
)

type helpCmd struct{}

func (h *helpCmd) Name() string {
	return "help"
}

func (h *helpCmd) Description() string {
	return ""
}

func (h *helpCmd) IsDisableable() bool {
	return false
}

func (h *helpCmd) Initialize() []implementation.Command {
	implementation.EnableParentCommand("help")

	return []implementation.Command{
		New(),
	}
}

func init() {
	implementation.CreateParentCommand(&helpCmd{})
}
