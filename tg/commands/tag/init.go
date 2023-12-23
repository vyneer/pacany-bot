package tag

import (
	"github.com/vyneer/tg-tagbot/tg/commands/implementation"
	"github.com/vyneer/tg-tagbot/tg/commands/tag/adduser"
	"github.com/vyneer/tg-tagbot/tg/commands/tag/info"
	"github.com/vyneer/tg-tagbot/tg/commands/tag/list"
	"github.com/vyneer/tg-tagbot/tg/commands/tag/new"
	"github.com/vyneer/tg-tagbot/tg/commands/tag/remove"
	"github.com/vyneer/tg-tagbot/tg/commands/tag/removeuser"
	"github.com/vyneer/tg-tagbot/tg/commands/tag/scan"
)

func init() {
	implementation.CreateCommand(New)
	implementation.CreateCommand(new.New)
	implementation.CreateCommand(remove.New)
	implementation.CreateCommand(adduser.New)
	implementation.CreateCommand(removeuser.New)
	implementation.CreateCommand(info.New)
	implementation.CreateCommand(list.New)
	implementation.CreateCommand(scan.New)
}
