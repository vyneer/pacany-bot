package new

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/vyneer/pacany-bot/tg/commands/implementation"
	tag_errors "github.com/vyneer/pacany-bot/tg/commands/tag/internal/errors"
	"github.com/vyneer/pacany-bot/tg/commands/tag/internal/util"
)

const (
	name              string = "new"
	parentName        string = "tag"
	help              string = "Add a new tag"
	arguments         string = "<tag_name> [description] <username>..."
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

	if len(a.Args) < 2 {
		resp.Text, _ = c.GetHelp()
		return []implementation.CommandResponse{
			resp,
		}
	}

	name := a.Args[0]
	if !util.IsValidTagName(name) {
		resp.Text = tag_errors.ErrInvalidTag.Error()
		return []implementation.CommandResponse{
			resp,
		}
	}

	lastDesc := 1
	descriptionSplit := []string{}
	for i, v := range a.Args[1:] {
		if util.IsValidUserName(v) {
			lastDesc += i
			break
		}
		descriptionSplit = append(descriptionSplit, v)
	}
	description := strings.Join(descriptionSplit, " ")

	mentions := util.FilterInvalidUsernames(a.Args[lastDesc:])
	if len(mentions) == 0 {
		resp.Text = tag_errors.ErrNoValidUsers.Error()
		return []implementation.CommandResponse{
			resp,
		}
	}

	err := a.DB.NewTag(ctx, a.ChatID, name, description, mentions...)
	if err != nil {
		slog.Warn("unable to create new tag", "err", err)
		resp.Text = err.Error()
		return []implementation.CommandResponse{
			resp,
		}
	}

	resp.Text = fmt.Sprintf("Added tag \"%s\"", name)

	return []implementation.CommandResponse{
		resp,
	}
}
