package scan

import (
	"context"
	"fmt"
	"log/slog"
	"slices"
	"strings"

	"github.com/vyneer/pacany-bot/tg/commands/implementation"
	"github.com/vyneer/pacany-bot/tg/commands/tag/internal/util"
)

const (
	name              string = "scan"
	parentName        string = "tag"
	help              string = ""
	shape             string = ""
	showInCommandList bool   = false
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

func (c *Command) Run(ctx context.Context, a implementation.CommandArgs) implementation.CommandResponse {
	tags, err := a.DB.GetTags(ctx, a.ChatID)
	if err != nil {
		slog.Warn("unable to get tags", "err", err)
		return implementation.CommandResponse{
			Text:       err.Error(),
			Reply:      true,
			Capitalize: true,
		}
	}

	username := a.Args[0]
	text := a.Args[1]

	fields := strings.Fields(text)

	tagCount := 0
	var allMentions []string
	description := ""
	for _, v := range tags {
		if i := slices.Index[[]string](fields, v.Name); i != -1 {
			if tagCount == 0 {
				description = v.Description
			}
			tagCount++
			allMentions = append(allMentions, strings.Fields(v.Mentions)...)
			fields = append(fields[:i], fields[i+1:]...)
		}
	}

	if len(allMentions) > 0 {
		if filtered, ok := util.FilterMentions(allMentions, username); ok {
			msg := filtered
			if tagCount == 1 {
				msg = strings.Join([]string{description, filtered}, " ")
			}
			if len(fields) > 0 {
				msg = strings.Join([]string{strings.Join(fields, " "), filtered}, " ")
			}

			return implementation.CommandResponse{
				Text:       msg,
				Reply:      false,
				Capitalize: false,
			}
		}
		return implementation.CommandResponse{
			Text:       "You're the only person using this tag",
			Reply:      true,
			Capitalize: true,
		}
	}

	return implementation.CommandResponse{
		Text:       "",
		Reply:      false,
		Capitalize: false,
	}
}
