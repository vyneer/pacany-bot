package remove

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/vyneer/pacany-bot/tg/commands/implementation"
	tz_errors "github.com/vyneer/pacany-bot/tg/commands/tz/internal/errors"
	"github.com/vyneer/pacany-bot/tg/commands/tz/internal/util"
)

const (
	name              string = "remove"
	parentName        string = "tz"
	help              string = "Remove specified user's timezone"
	arguments         string = "<username>"
	showInCommandList bool   = true
	showInHelp        bool   = true
	adminOnly         bool   = true
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

func (c *Command) Run(ctx context.Context, a implementation.CommandArgs) []implementation.CommandResponse {
	resp := implementation.CommandResponse{
		Reply:      true,
		Capitalize: true,
	}

	if len(a.Args) != 1 {
		resp.Text, _ = c.GetHelp()
		return []implementation.CommandResponse{
			resp,
		}
	}

	username := a.Args[0]
	if !util.IsValidUserName(username) {
		resp.Text = tz_errors.ErrInvalidUsername.Error()
		return []implementation.CommandResponse{
			resp,
		}
	}
	username = strings.TrimPrefix(username, "@")

	err := a.DB.RemoveTimezone(ctx, a.ChatID, username)
	if err != nil {
		slog.Warn("unable to remove timezone", "err", err)
		resp.Text = err.Error()
		return []implementation.CommandResponse{
			resp,
		}
	}

	resp.Text = fmt.Sprintf("Removed %s's timezone", username)

	return []implementation.CommandResponse{
		resp,
	}
}
