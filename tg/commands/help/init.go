package help

import (
	"github.com/vyneer/pacany-bot/config"
	"github.com/vyneer/pacany-bot/tg/commands/implementation"
)

type Parent struct{}

func NewHelp() *Parent {
	return &Parent{}
}

func (h *Parent) Name() string {
	return "help"
}

func (h *Parent) Description() string {
	return ""
}

func (h *Parent) IsDisableable() bool {
	return false
}

func (h *Parent) Initialize(cfg *config.Config) []implementation.Command {
	return []implementation.Command{
		New(cfg.GetParentCommandList()),
	}
}
