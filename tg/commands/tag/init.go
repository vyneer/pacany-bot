package tag

import (
	"github.com/vyneer/pacany-bot/config"
	"github.com/vyneer/pacany-bot/tg/commands/implementation"
	"github.com/vyneer/pacany-bot/tg/commands/tag/add"
	"github.com/vyneer/pacany-bot/tg/commands/tag/del"
	"github.com/vyneer/pacany-bot/tg/commands/tag/desc"
	"github.com/vyneer/pacany-bot/tg/commands/tag/descdel"
	"github.com/vyneer/pacany-bot/tg/commands/tag/help"
	"github.com/vyneer/pacany-bot/tg/commands/tag/info"
	"github.com/vyneer/pacany-bot/tg/commands/tag/internal/errors"
	"github.com/vyneer/pacany-bot/tg/commands/tag/internal/util"
	"github.com/vyneer/pacany-bot/tg/commands/tag/kick"
	tagname "github.com/vyneer/pacany-bot/tg/commands/tag/name"
	tagnew "github.com/vyneer/pacany-bot/tg/commands/tag/new"
	"github.com/vyneer/pacany-bot/tg/commands/tag/scan"
)

const name string = "tag"

type tag struct{}

func (t *tag) Name() string {
	return "tag"
}

func (t *tag) Description() string {
	return "tagging"
}

func (t *tag) Initialize() {
	implementation.EnableParentCommand(name)

	implementation.CreateInteractableCommand(help.New)
	implementation.CreateInteractableCommand(tagnew.New)
	implementation.CreateInteractableCommand(del.New)
	implementation.CreateInteractableCommand(tagname.New)
	implementation.CreateInteractableCommand(descdel.New)
	implementation.CreateInteractableCommand(desc.New)
	implementation.CreateInteractableCommand(add.New)
	implementation.CreateInteractableCommand(kick.New)
	implementation.CreateInteractableCommand(info.New)
	implementation.CreateAutomaticCommand(scan.New)
}

func (t *tag) Configure(cfg *config.Config) {
	util.SetTagPrefix(cfg.AllowedTagPrefixSymbols)
	errors.SetErrInvalidTag(cfg.AllowedTagPrefixSymbols)
}

func init() {
	implementation.CreateParentCommand(&tag{})
}
