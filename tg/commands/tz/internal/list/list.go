package list

import (
	"github.com/vyneer/pacany-bot/tg/commands/implementation"
	"github.com/vyneer/pacany-bot/tg/commands/tz/add"
	tzClear "github.com/vyneer/pacany-bot/tg/commands/tz/clear"
	"github.com/vyneer/pacany-bot/tg/commands/tz/convert"
	"github.com/vyneer/pacany-bot/tg/commands/tz/info"
	"github.com/vyneer/pacany-bot/tg/commands/tz/remove"
	"github.com/vyneer/pacany-bot/tg/commands/tz/set"
)

var Commands = []implementation.Command{
	set.New(),
	tzClear.New(),
	info.New(),
	convert.New(),
	add.New(),
	remove.New(),
}
