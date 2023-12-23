package list

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/vyneer/tg-tagbot/tg/commands/implementation"
)

const (
	name        string = "list"
	parentName  string = "tag"
	help        string = "/tag list - List all tags and their associated user count"
	helpOrder   int    = 6
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
	resp := implementation.CommandResponse{
		Reply: true,
	}

	tags, err := a.DB.GetTags(ctx, a.ChatID)
	if err != nil {
		slog.Warn("unable to get tags", "err", err)
		resp.Text = err.Error()
		return resp
	}

	var tagNames []string
	for _, v := range tags {
		l := len(strings.Fields(v.Mentions))
		tagNames = append(tagNames, fmt.Sprintf("%s - %d user%s", v.Name, l, func() string {
			if l != 1 {
				return "s"
			}
			return ""
		}()))
	}

	if len(tagNames) == 0 {
		resp.Text = "No tags in this group chat"
		return resp
	}

	resp.Text = strings.Join(tagNames, "\n")

	return resp
}
