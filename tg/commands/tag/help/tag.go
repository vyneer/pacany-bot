package help

import (
	"context"
	"fmt"
	"strings"

	"github.com/vyneer/pacany-bot/tg/commands/implementation"
)

const (
	name              string = "help"
	parentName        string = "tag"
	help              string = "Print help for the tag commands"
	shape             string = "/taghelp"
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
	helpSlice := []string{}

	for _, v := range implementation.InteractableOrder {
		if v.GetParentName() != parentName {
			continue
		}
		if helpString, show := v.GetHelp(); show {
			helpSlice = append(helpSlice, helpString)
		}
	}

	return implementation.CommandResponse{
		Text:       strings.Join(helpSlice, "\n\n"),
		Reply:      true,
		Capitalize: true,
	}
}
