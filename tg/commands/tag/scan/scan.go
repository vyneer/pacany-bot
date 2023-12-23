package scan

import (
	"context"
	"fmt"
	"log/slog"
	"slices"
	"strings"

	"github.com/vyneer/tg-tagbot/tg/commands/implementation"
	"github.com/vyneer/tg-tagbot/tg/commands/tag/internal/util"
)

const (
	name              string = "scan"
	parentName        string = "tag"
	help              string = ""
	helpOrder         int    = -1
	shape             string = ""
	showInCommandList bool   = false
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

func (c *Command) GetHelp() (string, int) {
	return fmt.Sprintf("%s - %s", shape, help), helpOrder
}

func (c *Command) GetDescription() string {
	if !showInCommandList {
		return ""
	}
	return fmt.Sprintf("%s - %s", help, shape)
}

func (c *Command) Run(ctx context.Context, a implementation.CommandArgs) implementation.CommandResponse {
	tags, err := a.DB.GetTags(ctx, a.ChatID)
	if err != nil {
		slog.Warn("unable to get tags", "err", err)
		return implementation.CommandResponse{
			Text:  err.Error(),
			Reply: true,
		}
	}

	username := a.Args[0]
	text := a.Args[1]

	fields := strings.Fields(text)

	var allMentions []string
	for _, v := range tags {
		if slices.Contains[[]string](fields, v.Name) {
			allMentions = append(allMentions, strings.Fields(v.Mentions)...)
		}
	}

	if len(allMentions) > 0 {
		if filtered, ok := util.FilterMentions(allMentions, username); ok {
			return implementation.CommandResponse{
				Text:  filtered,
				Reply: false,
			}
		}
		return implementation.CommandResponse{
			Text:  "You're the only person using this tag",
			Reply: true,
		}
	}

	return implementation.CommandResponse{
		Text:  "",
		Reply: false,
	}
}
