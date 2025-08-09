package tag

import (
	"github.com/vyneer/pacany-bot/config"
	"github.com/vyneer/pacany-bot/tg/commands/implementation"
	"github.com/vyneer/pacany-bot/tg/commands/tag/help"
	"github.com/vyneer/pacany-bot/tg/commands/tag/internal/errors"
	"github.com/vyneer/pacany-bot/tg/commands/tag/internal/list"
	"github.com/vyneer/pacany-bot/tg/commands/tag/internal/util"
)

const (
	name        string = "tag"
	description string = "tagging"
)

type Parent struct{}

func NewTag() *Parent {
	return &Parent{}
}

func (t *Parent) Name() string {
	return name
}

func (t *Parent) Description() string {
	return description
}

func (t *Parent) IsDisableable() bool {
	return true
}

func (t *Parent) Initialize(cfg *config.Config) []implementation.Command {
	util.SetTagPrefix(cfg.AllowedTagPrefixSymbols)
	errors.SetErrInvalidTag(cfg.AllowedTagPrefixSymbols)

	cmds := []implementation.Command{
		help.New(),
	}

	return append(cmds, list.Commands...)
}
