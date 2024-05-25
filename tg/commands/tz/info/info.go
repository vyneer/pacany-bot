package info

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/vyneer/pacany-bot/db"
	"github.com/vyneer/pacany-bot/tg/commands/implementation"
)

const (
	name              string = "info"
	parentName        string = "tz"
	help              string = "Get timezone list"
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
		Capitalize: false,
	}

	tzs, err := a.DB.GetTimezones(ctx, a.ChatID)
	if err != nil {
		slog.Warn("unable to get timezones", "err", err)
		resp.Text = err.Error()
		return []implementation.CommandResponse{
			resp,
		}
	}

	t := time.Now()
	timezoneMap := map[int][]db.Timezone{}
	for _, v := range tzs {
		if v.Username != a.User.UserName {
			tz, _ := time.LoadLocation(v.Timezone)
			_, offset := t.In(tz).Zone()
			timezoneMap[offset] = append(timezoneMap[offset], v)
		}
	}

	var timezonesPretty []string
	for k, timezoneSlice := range timezoneMap {
		_, offset := t.Zone()
		if k == offset {
			continue
		}

		var names []string
		var tz *time.Location
		for i, v := range timezoneSlice {
			if i == 0 {
				tz, _ = time.LoadLocation(v.Timezone)
			}

			if v.Description != "" {
				names = append(names, v.Description)
			} else {
				names = append(names, v.Username)
			}
		}

		if len(names) == 0 {
			continue
		}

		msg := strings.Join(names, ", ")

		timezonesPretty = append(timezonesPretty, fmt.Sprintf("%s - %s", msg, t.In(tz).Format("02/01 15:04")))
	}

	if len(timezonesPretty) == 0 {
		resp.Text = "No timezones in this group chat"
		return []implementation.CommandResponse{
			resp,
		}
	}

	resp.Text = strings.Join(timezonesPretty, "\n")

	return []implementation.CommandResponse{
		resp,
	}
}
