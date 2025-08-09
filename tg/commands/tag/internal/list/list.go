package list

import (
	"github.com/vyneer/pacany-bot/tg/commands/implementation"
	"github.com/vyneer/pacany-bot/tg/commands/tag/add"
	"github.com/vyneer/pacany-bot/tg/commands/tag/del"
	"github.com/vyneer/pacany-bot/tg/commands/tag/desc"
	"github.com/vyneer/pacany-bot/tg/commands/tag/descdel"
	"github.com/vyneer/pacany-bot/tg/commands/tag/info"
	"github.com/vyneer/pacany-bot/tg/commands/tag/kick"
	tagname "github.com/vyneer/pacany-bot/tg/commands/tag/name"
	tagnew "github.com/vyneer/pacany-bot/tg/commands/tag/new"
	"github.com/vyneer/pacany-bot/tg/commands/tag/scan"
)

var Commands = []implementation.Command{
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
