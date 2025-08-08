package clear

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/vyneer/pacany-bot/tg/commands/implementation"
)

const (
	name              string = "clear"
	parentName        string = "tz"
	help              string = "Clear your timezone"
	arguments         string = ""
	showInCommandList bool   = true
	showInHelp        bool   = true
	adminOnly         bool   = false
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

func (c *Command) Run(ctx context.Context, a implementation.CommandArgs) []implementation.CommandResponse {
	resp := implementation.CommandResponse{
		Reply:      true,
		Capitalize: true,
	}

	err := a.DB.RemoveTimezone(ctx, a.ChatID, a.User.Username)
	if err != nil {
		slog.Warn("unable to remove timezone", "err", err)
		resp.Text = err.Error()
		return []implementation.CommandResponse{
			resp,
		}
	}

	resp.Text = "Cleared your timezone"

	return []implementation.CommandResponse{
		resp,
	}
}
