package tag

import (
	"github.com/vyneer/pacany-bot/config"
	"github.com/vyneer/pacany-bot/tg/commands/implementation"
	"github.com/vyneer/pacany-bot/tg/commands/tag/help"
	"github.com/vyneer/pacany-bot/tg/commands/tag/internal/errors"
	"github.com/vyneer/pacany-bot/tg/commands/tag/internal/list"
	"github.com/vyneer/pacany-bot/tg/commands/tag/internal/util"
)

const name string = "tag"

type tag struct{}

func (t *tag) Name() string {
	return "tag"
}

func (t *tag) Description() string {
	return "tagging"
}

func (t *tag) IsDisableable() bool {
	return true
}

func (t *tag) Initialize() []implementation.Command {
	implementation.EnableParentCommand(name)

	cmds := []implementation.Command{
		help.New(),
	}

	return append(cmds, list.Commands...)
}

func (t *tag) Configure(cfg *config.Config) {
	util.SetTagPrefix(cfg.AllowedTagPrefixSymbols)
	errors.SetErrInvalidTag(cfg.AllowedTagPrefixSymbols)
}

func init() {
	implementation.CreateParentCommand(&tag{})
}
