package rename

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/vyneer/pacany-bot/tg/commands/implementation"
	tag_errors "github.com/vyneer/pacany-bot/tg/commands/tag/internal/errors"
	"github.com/vyneer/pacany-bot/tg/commands/tag/internal/util"
)

const (
	name              string = "rename"
	parentName        string = "tag"
	help              string = "Rename the specified tag"
	helpOrder         int    = 2
	shape             string = "/tagrename <tag_old_name> <tag_new_name>"
	descriptionOrder  int    = 4
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

func (c *Command) GetDescription() (string, int) {
	if !showInCommandList {
		return "", descriptionOrder
	}
	return fmt.Sprintf("%s - %s", help, shape), descriptionOrder
}

func (c *Command) Run(ctx context.Context, a implementation.CommandArgs) implementation.CommandResponse {
	resp := implementation.CommandResponse{
		Reply:      true,
		Capitalize: true,
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

	if oldName == newName {
		resp.Text = "Identical name provided"
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
