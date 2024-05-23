package info

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/vyneer/pacany-bot/tg/commands/implementation"
)

const (
	name              string = "info"
	parentName        string = "tz"
	help              string = "Get timezone list"
	arguments         string = ""
	showInCommandList bool   = true
	showInHelp        bool   = true
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

func (c *Command) Run(ctx context.Context, a implementation.CommandArgs) implementation.CommandResponse {
	resp := implementation.CommandResponse{
		Reply:      true,
		Capitalize: true,
	}

	tzs, err := a.DB.GetTimezones(ctx, a.ChatID)
	if err != nil {
		slog.Warn("unable to get timezones", "err", err)
		resp.Text = err.Error()
		return resp
	}

	var timezonesPretty []string
	for _, v := range tzs {
		tz, _ := time.LoadLocation(v.Timezone)
		timezonesPretty = append(timezonesPretty, fmt.Sprintf("%s - %s - %s", v.Name, v.Description, time.Now().In(tz).Format("2006-01-02 15:04:05 -07:00")))
	}

	if len(timezonesPretty) == 0 {
		resp.Text = "No timezones in this group chat"
		return resp
	}

	resp.Text = strings.Join(timezonesPretty, "\n")

	return resp
}
