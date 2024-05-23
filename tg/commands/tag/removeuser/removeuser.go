package removeuser

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/vyneer/pacany-bot/db"
	"github.com/vyneer/pacany-bot/tg/commands/implementation"
	tag_errors "github.com/vyneer/pacany-bot/tg/commands/tag/internal/errors"
	"github.com/vyneer/pacany-bot/tg/commands/tag/internal/util"
)

const (
	name              string = "removeuser"
	parentName        string = "tag"
	help              string = "Remove specified users from an existing tag"
	arguments         string = "<tag_name> <username>..."
	showInCommandList bool   = true
	showInHelp        bool   = true
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

func (c *Command) Run(ctx context.Context, a implementation.CommandArgs) implementation.CommandResponse {
	resp := implementation.CommandResponse{
		Reply:      true,
		Capitalize: true,
	}

	if len(a.Args) < 2 {
		resp.Text, _ = c.GetHelp()
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

	err := a.DB.RemoveMentionsFromTag(ctx, a.ChatID, name, mentions...)
	if err != nil {
		if errors.Is(err, db.ErrEmptyTag) {
			err := a.DB.RemoveTag(ctx, a.ChatID, name)
			if err != nil {
				slog.Warn("unable to remove tag", "err", err)
				resp.Text = err.Error()
				return resp
			}
			resp.Text = fmt.Sprintf("Removed tag \"%s\"", name)
			return resp
		}
		slog.Warn("unable to add mentions to tag", "err", err)
		resp.Text = err.Error()
		return resp
	}

	resp.Text = fmt.Sprintf("Removed user%s from tag \"%s\"", func() string {
		if len(mentions) != 1 {
			return "s"
		}
		return ""
	}(), name)

	return resp
}
