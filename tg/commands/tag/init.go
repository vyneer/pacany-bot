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
		tagnew.New(),
		del.New(),
		tagname.New(),
		desc.New(),
		descdel.New(),
		add.New(),
		kick.New(),
		info.New(),
		scan.New(),
	}

	return append([]implementation.Command{help.New(cmds)}, cmds...)
}
