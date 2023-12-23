package info

import (
	"context"
	"fmt"
	"log/slog"
	"slices"
	"strings"

	"github.com/vyneer/tg-tagbot/db"
	"github.com/vyneer/tg-tagbot/tg/commands/implementation"
	tag_errors "github.com/vyneer/tg-tagbot/tg/commands/tag/internal/errors"
	"github.com/vyneer/tg-tagbot/tg/commands/tag/internal/util"
)

const (
	name        string = "info"
	parentName  string = "tag"
	help        string = "/tag info <tag_name> - Get tag user count and user list"
	helpOrder   int    = 5
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

	if len(a.Args) != 1 {
		resp.Text = c.help
		return resp
	}

	tags, err := a.DB.GetTags(ctx, a.ChatID)
	if err != nil {
		slog.Warn("unable to get tags", "err", err)
		resp.Text = err.Error()
		return resp
	}

	name := a.Args[0]
	if !util.IsValidTagName(name) {
		resp.Text = tag_errors.ErrInvalidTag.Error()
		return resp
	}

	i := slices.IndexFunc[[]db.Tag](tags, func(t db.Tag) bool {
		return t.Name == name
	})
	if i == -1 {
		resp.Text = db.ErrTagDoesntExist.Error()
		return resp
	}

	var info []string
	fields := strings.Fields(tags[i].Mentions)
	info = append(info, fmt.Sprintf("Tag name: %s", tags[i].Name))
	info = append(info, fmt.Sprintf("User count: %d", len(fields)))
	info = append(info, "User list:")
	for _, v := range fields {
		info = append(info, fmt.Sprintf("- %s", strings.TrimPrefix(v, "@")))
	}

	resp.Text = strings.Join(info, "\n")

	return resp
}
