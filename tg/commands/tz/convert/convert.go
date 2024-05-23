package convert

import (
	"context"
	"fmt"
	"log/slog"
	"slices"
	"strings"
	"time"

	"github.com/vyneer/pacany-bot/db"
	"github.com/vyneer/pacany-bot/tg/commands/implementation"
	tz_errors "github.com/vyneer/pacany-bot/tg/commands/tz/internal/errors"
)

const (
	name              string = "convert"
	parentName        string = "tz"
	help              string = "Convert specified time to every timezone on the list"
	arguments         string = "<time>"
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

	if len(a.Args) < 1 {
		resp.Text, _ = c.GetHelp()
		return resp
	}

	tzs, err := a.DB.GetTimezones(ctx, a.ChatID)
	if err != nil {
		slog.Warn("unable to get timezones", "err", err)
		resp.Text = err.Error()
		return resp
	}

	i := slices.IndexFunc[[]db.Timezone](tzs, func(t db.Timezone) bool {
		return t.UserID == a.User.ID
	})
	if i == -1 {
		resp.Text = tz_errors.ErrTimezoneNotSet.Error()
		return resp
	}

	timeString := strings.Join(a.Args, " ")

	var t time.Time
	tz, _ := time.LoadLocation(tzs[i].Timezone)
	tFull, err1 := time.ParseInLocation("2006-01-02 15:04:05", timeString, tz)
	tSeconds, err2 := time.Parse("15:04:05", timeString)
	tMinutes, err3 := time.Parse("15:04", timeString)

	switch {
	case err1 == nil:
		t = tFull.UTC()
	case err2 == nil:
		t = time.Now().UTC()
		t = time.Date(t.Year(), t.Month(), t.Day(), tSeconds.Hour(), tSeconds.Minute(), tSeconds.Second(), 0, tz)
	case err3 == nil:
		t = time.Now().UTC()
		t = time.Date(t.Year(), t.Month(), t.Day(), tMinutes.Hour(), tMinutes.Minute(), 0, 0, tz)
	default:
		resp.Text = tz_errors.ErrUnableToParse.Error()
		return resp
	}

	var timezonesPretty []string
	for _, v := range tzs {
		if v.UserID != a.User.ID {
			tz, _ := time.LoadLocation(v.Timezone)
			timezonesPretty = append(timezonesPretty, fmt.Sprintf("%s - %s - %s", v.Name, v.Description, t.In(tz).Format("2006-01-02 15:04:05 -07:00")))
		}
	}

	if len(timezonesPretty) == 0 {
		resp.Text = "No other timezones in this group chat"
		return resp
	}

	resp.Text = strings.Join(timezonesPretty, "\n")

	return resp
}
