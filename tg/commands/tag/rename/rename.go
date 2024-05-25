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
	arguments         string = "<tag_old_name> <tag_new_name>"
	showInCommandList bool   = true
	showInHelp        bool   = true
	adminOnly         bool   = true
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
		Capitalize: true,
	}

	if len(a.Args) != 2 {
		resp.Text, _ = c.GetHelp()
		return []implementation.CommandResponse{
			resp,
		}
	}

	oldName := a.Args[0]
	if !util.IsValidTagName(oldName) {
		resp.Text = tag_errors.ErrInvalidTag.Error()
		return []implementation.CommandResponse{
			resp,
		}
	}

	newName := a.Args[1]
	if !util.IsValidTagName(newName) {
		resp.Text = tag_errors.ErrInvalidTag.Error()
		return []implementation.CommandResponse{
			resp,
		}
	}

	if oldName == newName {
		resp.Text = "Identical name provided"
		return []implementation.CommandResponse{
			resp,
		}
	}

	err := a.DB.RenameTag(ctx, a.ChatID, oldName, newName)
	if err != nil {
		slog.Warn("unable to rename tag", "err", err)
		resp.Text = err.Error()
		return []implementation.CommandResponse{
			resp,
		}
	}

	resp.Text = fmt.Sprintf("Renamed tag \"%s\" to \"%s\"", oldName, newName)

	return []implementation.CommandResponse{
		resp,
	}
}
