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
	arguments         string = ""
	showInCommandList bool   = true
	showInHelp        bool   = false
	adminOnly         bool   = false
)

type Command struct {
	commandList []implementation.Command
}

func New(commandList []implementation.Command) implementation.InteractableCommand {
	return &Command{
		commandList: commandList,
	}
}

func (c *Command) GetName() string {
	return name
}

func (c *Command) GetParentName() string {
	return parentName
}

func (c *Command) GetHelp() (string, bool) {
	if arguments == "" {
		return fmt.Sprintf("/%s%s - %s", parentName, name, help), showInHelp
	}
	return fmt.Sprintf("/%s%s %s - %s", parentName, name, arguments, help), showInHelp
}

func (c *Command) GetDescription() (string, bool) {
	if arguments == "" {
		return help, showInCommandList
	}
	return fmt.Sprintf("%s - %s", arguments, help), showInCommandList
}

func (c *Command) IsAdminOnly() bool {
	return adminOnly
}

func (c *Command) Run(_ context.Context, a implementation.CommandArgs) []implementation.CommandResponse {
	helpSlice := []string{}

	for _, cmd := range c.commandList {
		if v, ok := cmd.(implementation.InteractableCommand); ok {
			if !a.IsAdmin && v.IsAdminOnly() {
				continue
			}
			if v.GetParentName() != parentName {
				continue
			}
			if helpString, show := v.GetHelp(); show {
				helpSlice = append(helpSlice, helpString)
			}
		}
	}

	return []implementation.CommandResponse{
		{
			Text:       strings.Join(helpSlice, "\n\n"),
			Reply:      true,
			Capitalize: true,
		},
	}
}
