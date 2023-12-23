package help

import (
	"context"
	"fmt"

	"github.com/vyneer/tg-tagbot/tg/commands/implementation"
)

const (
	name              string = "help"
	parentName        string = ""
	help              string = "Print help"
	helpOrder         int    = -1
	shape             string = "/help"
	descriptionOrder  int    = 0
	showInCommandList bool   = true
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

func (c *Command) GetHelp() (string, int) {
	return fmt.Sprintf("%s - %s", shape, help), helpOrder
}

func (c *Command) GetDescription() (string, int) {
	if !showInCommandList {
		return "", descriptionOrder
	}
	return fmt.Sprintf("%s - %s", help, shape), descriptionOrder
}

func (c *Command) Run(_ context.Context, _ implementation.CommandArgs) implementation.CommandResponse {
	return implementation.CommandResponse{
		Text:  "Multi-purpose Telegram bot, check other help commands for more details.\n\nCurrently supported functions:\n- tagging (/taghelp)",
		Reply: true,
	}
}
