package help

import (
	"context"
	"fmt"
	"strings"

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
	arguments         string = ""
	showInCommandList bool   = true
	showInHelp        bool   = false
	adminOnly         bool   = false
)

type Command struct{}

func New() implementation.InteractableCommand {
	return &Command{}
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

func (c *Command) Run(_ context.Context, _ implementation.CommandArgs) []implementation.CommandResponse {
	text := []string{
		"Multi-purpose Telegram bot, check other help commands for more details.",
	}

	parents := implementation.GetEnabledParentCommands()
	if len(parents) > 0 {
		subText := []string{
			"Currently enabled functions:",
		}
		for _, parentCommand := range parents {
			subText = append(subText, fmt.Sprintf("- %s (/%shelp)", parentCommand.Description(), parentCommand.Name()))
		}
		text = append(text, strings.Join(subText, "\n"))
	} else {
		text = append(text, "No functions are currently enabled.")
	}

	text = append(text, fmt.Sprintf("Version: v%s-%s-%s", Version, Commit, Timestamp))

	return []implementation.CommandResponse{
		{
			Text:       strings.Join(text, "\n\n"),
			Reply:      true,
			Capitalize: true,
		},
	}
}
