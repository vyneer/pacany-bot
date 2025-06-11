package tag

import (
	"github.com/vyneer/pacany-bot/config"
	"github.com/vyneer/pacany-bot/tg/commands/implementation"
	"github.com/vyneer/pacany-bot/tg/commands/tag/add"
	"github.com/vyneer/pacany-bot/tg/commands/tag/adduser"
	"github.com/vyneer/pacany-bot/tg/commands/tag/changedesc"
	"github.com/vyneer/pacany-bot/tg/commands/tag/help"
	"github.com/vyneer/pacany-bot/tg/commands/tag/info"
	"github.com/vyneer/pacany-bot/tg/commands/tag/internal/errors"
	"github.com/vyneer/pacany-bot/tg/commands/tag/internal/util"
	"github.com/vyneer/pacany-bot/tg/commands/tag/remove"
	"github.com/vyneer/pacany-bot/tg/commands/tag/removedesc"
	"github.com/vyneer/pacany-bot/tg/commands/tag/removeuser"
	"github.com/vyneer/pacany-bot/tg/commands/tag/rename"
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
	implementation.CreateInteractableCommand(add.New)
	implementation.CreateInteractableCommand(remove.New)
	implementation.CreateInteractableCommand(rename.New)
	implementation.CreateInteractableCommand(removedesc.New)
	implementation.CreateInteractableCommand(changedesc.New)
	implementation.CreateInteractableCommand(adduser.New)
	implementation.CreateInteractableCommand(removeuser.New)
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
