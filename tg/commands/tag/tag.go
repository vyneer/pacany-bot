package tag

import (
	"context"

	"github.com/vyneer/tg-tagbot/tg/commands/implementation"
)

const (
	name        string = "tag"
	parentName  string = ""
	help        string = ""
	helpOrder   int    = -1
	description string = "Manage tags"
)

type Command struct {
	name        string
	parentName  string
	help        string
	helpOrder   int
	description string
}

func New() implementation.Command {
	return &Command{
		name:        name,
		parentName:  parentName,
		help:        help,
		helpOrder:   helpOrder,
		description: description,
	}
}

func (c *Command) GetName() string {
	return c.name
}

func (c *Command) GetParentName() string {
	return c.parentName
}

func (c *Command) GetHelp() (string, int) {
	return c.help, c.helpOrder
}

func (c *Command) GetDescription() string {
	return c.description
}

func (c *Command) Run(_ context.Context, _ implementation.CommandArgs) implementation.CommandResponse {
	return implementation.CommandResponse{
		Text:  "/tag <new|remove|add-user|remove-user|info|list> ...\n\nFor more information use /help",
		Reply: true,
	}
}
