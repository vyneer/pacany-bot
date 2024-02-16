package help

import (
	"context"
	"fmt"

	"github.com/vyneer/pacany-bot/tg/commands/implementation"
)

var (
	Version   = "dev"
	Commit    = "deadbeef"
	Timestamp = "0"
)

const (
	name              string = "help"
	parentName        string = ""
	help              string = "Print help"
	shape             string = "/help"
	showInCommandList bool   = true
	showInHelp        bool   = false
)

type Command struct{}

func New() implementation.Command {
	return &Command{}
}

func (c *Command) GetName() string {
	return name
}

func (c *Command) GetParentName() string {
	return parentName
}

func (c *Command) GetHelp() (string, bool) {
	return fmt.Sprintf("%s - %s", shape, help), showInHelp
}

func (c *Command) GetDescription() (string, bool) {
	return fmt.Sprintf("%s - %s", help, shape), showInCommandList
}

func (c *Command) Run(_ context.Context, _ implementation.CommandArgs) implementation.CommandResponse {
	return implementation.CommandResponse{
		Text:       fmt.Sprintf("Multi-purpose Telegram bot, check other help commands for more details.\n\nCurrently supported functions:\n- tagging (/taghelp)\n\nVersion: v%s-%s-%s", Version, Commit, Timestamp),
		Reply:      true,
		Capitalize: true,
	}
}
