package remove

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/vyneer/tg-tagbot/tg/commands/implementation"
	tag_errors "github.com/vyneer/tg-tagbot/tg/commands/tag/internal/errors"
	"github.com/vyneer/tg-tagbot/tg/commands/tag/internal/util"
)

const (
	name        string = "remove"
	parentName  string = "tag"
	help        string = "/tag remove <tag_name> - Remove the specified tag"
	helpOrder   int    = 2
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

	name := a.Args[0]
	if !util.IsValidTagName(name) {
		resp.Text = tag_errors.ErrInvalidTag.Error()
		return resp
	}

	err := a.DB.RemoveTag(ctx, a.ChatID, name)
	if err != nil {
		slog.Warn("unable to remove tag", "err", err)
		resp.Text = err.Error()
		return resp
	}

	resp.Text = fmt.Sprintf("Removed tag \"%s\"", name)

	return resp
}
