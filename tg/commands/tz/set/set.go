package set

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/vyneer/pacany-bot/geonames"
	"github.com/vyneer/pacany-bot/tg/commands/implementation"
	tz_errors "github.com/vyneer/pacany-bot/tg/commands/tz/internal/errors"
)

const (
	name              string = "set"
	parentName        string = "tz"
	help              string = "Set your timezone"
	arguments         string = "<timezone> [description]"
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

	if len(a.Args) < 1 {
		resp.Text, _ = c.GetHelp()
		return []implementation.CommandResponse{
			resp,
		}
	}

	timezone := a.Args[0]
	tz, err := time.LoadLocation(timezone)
	if err != nil {
		tz, err = geonames.CityToTimezone.Get(timezone)
		if err != nil {
			resp.Text = tz_errors.ErrInvalidTimezone.Error()
			return []implementation.CommandResponse{
				resp,
			}
		}
	}

	descriptionSplit := []string{}
	descriptionSplit = append(descriptionSplit, a.Args[1:]...)

	description := a.User.UserName
	if len(descriptionSplit) > 0 {
		description = strings.Join(descriptionSplit, " ")
	}

	err = a.DB.NewTimezone(ctx, a.ChatID, a.User.UserName, tz.String(), description)
	if err != nil {
		slog.Warn("unable to set timezone", "err", err)
		resp.Text = err.Error()
		return []implementation.CommandResponse{
			resp,
		}
	}

	resp.Text = fmt.Sprintf("Set your timezone to \"%s\"", tz.String())

	return []implementation.CommandResponse{
		resp,
	}
}
