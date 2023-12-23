package new

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/vyneer/tg-tagbot/tg/commands/implementation"
	tag_errors "github.com/vyneer/tg-tagbot/tg/commands/tag/internal/errors"
	"github.com/vyneer/tg-tagbot/tg/commands/tag/internal/util"
)

const (
	name        string = "new"
	parentName  string = "tag"
	help        string = "/tag new <tag_name> <username_1> <username_2> ... <username_n> - Add a new tag"
	helpOrder   int    = 1
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

	if len(a.Args) < 2 {
		resp.Text = c.help
		return resp
	}

	name := a.Args[0]
	if !util.IsValidTagName(name) {
		resp.Text = tag_errors.ErrInvalidTag.Error()
		return resp
	}
	mentions := util.FilterInvalidUsernames(a.Args[1:])
	if len(mentions) == 0 {
		resp.Text = tag_errors.ErrNoValidUsers.Error()
		return resp
	}

	err := a.DB.NewTag(ctx, a.ChatID, name, mentions...)
	if err != nil {
		slog.Warn("unable to create new tag", "err", err)
		resp.Text = err.Error()
		return resp
	}

	resp.Text = fmt.Sprintf("Added tag \"%s\"", name)

	return resp
}
