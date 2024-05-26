package convert

import (
	"cmp"
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

	if len(a.Args) < 1 {
		resp.Text, _ = c.GetHelp()
		return []implementation.CommandResponse{
			resp,
		}
	}

	tzs, err := a.DB.GetTimezones(ctx, a.ChatID)
	if err != nil {
		slog.Warn("unable to get timezones", "err", err)
		resp.Text = err.Error()
		return []implementation.CommandResponse{
			resp,
		}
	}

	i := slices.IndexFunc(tzs, func(t db.Timezone) bool {
		return t.Username == a.User.UserName
	})
	if i == -1 {
		resp.Text = tz_errors.ErrTimezoneNotSet.Error()
		return []implementation.CommandResponse{
			resp,
		}
	}

	timeString := strings.Join(a.Args, " ")

	var t time.Time
	tzString := tzs[i].Timezone
	tz, _ := time.LoadLocation(tzString)
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
		return []implementation.CommandResponse{
			resp,
		}
	}

	timezoneMap := map[int][]db.Timezone{}
	for _, v := range tzs {
		if v.Username != a.User.UserName {
			tz, _ := time.LoadLocation(v.Timezone)
			_, offset := t.In(tz).Zone()
			timezoneMap[offset] = append(timezoneMap[offset], v)
		}
	}

	if len(timezoneMap) == 0 {
		resp.Text = "No timezones in this group chat"
		return []implementation.CommandResponse{
			resp,
		}
	}

	sortedByTz := make([]int, 0, len(timezoneMap))
	sortedByLength := make([]int, 0, len(timezoneMap))
	for k := range timezoneMap {
		sortedByTz = append(sortedByTz, k)
		sortedByLength = append(sortedByLength, k)
	}

	slices.Sort(sortedByTz)
	slices.SortFunc(sortedByLength, func(a int, b int) int {
		return cmp.Compare(len(timezoneMap[b]), len(timezoneMap[a]))
	})

	biggestTz := -1
	if len(sortedByLength) > 1 {
		first := timezoneMap[sortedByLength[0]]
		second := timezoneMap[sortedByLength[1]]
		if len(first) > len(second) && len(first) > 3 {
			biggestTz = sortedByLength[0]
		}
	}

	var timezonesPretty []string
	for _, k := range sortedByTz {
		timezoneSlice := timezoneMap[k]

		var tz *time.Location
		if k == biggestTz {
			continue
		}

		var names []string
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

	if biggestTz != -1 {
		biggest := timezoneMap[biggestTz]
		tz, _ := time.LoadLocation(biggest[0].Timezone)
		timezonesPretty = append(timezonesPretty, []string{"", fmt.Sprintf("Остальные (%d чел.) - %s", len(biggest), t.In(tz).Format("02/01 15:04"))}...)
	}

	if len(timezonesPretty) == 0 {
		resp.Text = "No other timezones in this group chat"
		return []implementation.CommandResponse{
			resp,
		}
	}

	resp.Text = strings.Join(timezonesPretty, "\n")

	return []implementation.CommandResponse{
		resp,
	}
}
