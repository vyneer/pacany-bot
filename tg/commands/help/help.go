package help

import (
	"context"
	"strings"

	"github.com/vyneer/tg-tagbot/tg/commands/implementation"
)

const (
	name        string = "help"
	parentName  string = ""
	help        string = "/help - Print this message"
	helpOrder   int    = 0
	description string = "Prints help"
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
	helpMap := map[int]string{}

	for _, v := range implementation.Map {
		helpString, order := v.GetHelp()
		if order != -1 {
			helpMap[order] = helpString
		}
	}

	var helpSlice []string
	for i := 0; i < len(helpMap); i++ {
		helpSlice = append(helpSlice, helpMap[i])
	}

	return implementation.CommandResponse{
		Text:  strings.Join(helpSlice, "\n\n"),
		Reply: true,
	}
}
