package scan

import (
	"context"
	"log/slog"
	"slices"
	"strings"

	"github.com/vyneer/tg-tagbot/tg/commands/implementation"
	"github.com/vyneer/tg-tagbot/tg/commands/tag/internal/util"
)

const (
	name        string = "scan"
	parentName  string = "tag"
	help        string = ""
	helpOrder   int    = -1
	description string = ""
)

type Command struct {
	name        string
	parentName  string
	help        string
	helpOrder   int
	description string
}

func New() implementation.Command {
	return &Command{
		name:        name,
		parentName:  parentName,
		help:        help,
		helpOrder:   helpOrder,
		description: description,
	}
}

func (c *Command) GetName() string {
	return c.name
}

func (c *Command) GetParentName() string {
	return c.parentName
}

func (c *Command) GetHelp() (string, int) {
	return c.help, c.helpOrder
}

func (c *Command) GetDescription() string {
	return c.description
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

	for _, v := range tags {
		if slices.Contains[[]string](fields, v.Name) {
			if filtered, ok := util.FilterMentions(v.Mentions, username); ok {
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
	}

	return implementation.CommandResponse{
		Text:  "",
		Reply: false,
	}
}
