package rename

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/vyneer/tg-tagbot/tg/commands/implementation"
	tag_errors "github.com/vyneer/tg-tagbot/tg/commands/tag/internal/errors"
	"github.com/vyneer/tg-tagbot/tg/commands/tag/internal/util"
)

const (
	name              string = "rename"
	parentName        string = "tag"
	help              string = "Rename the specified tag"
	helpOrder         int    = 2
	shape             string = "/tagrename <tag_old_name> <tag_new_name>"
	showInCommandList bool   = true
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
	resp := implementation.CommandResponse{
		Reply: true,
	}

	if len(a.Args) != 2 {
		resp.Text, _ = c.GetHelp()
		return resp
	}

	oldName := a.Args[0]
	if !util.IsValidTagName(oldName) {
		resp.Text = tag_errors.ErrInvalidTag.Error()
		return resp
	}

	newName := a.Args[1]
	if !util.IsValidTagName(newName) {
		resp.Text = tag_errors.ErrInvalidTag.Error()
		return resp
	}

	err := a.DB.RenameTag(ctx, a.ChatID, oldName, newName)
	if err != nil {
		slog.Warn("unable to rename tag", "err", err)
		resp.Text = err.Error()
		return resp
	}

	resp.Text = fmt.Sprintf("Renamed tag \"%s\" to \"%s\"", oldName, newName)

	return resp
}
