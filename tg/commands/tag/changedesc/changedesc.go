package changedesc

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
	name              string = "changedesc"
	parentName        string = "tag"
	help              string = "Change the description of a specified tag"
	arguments         string = "<tag_name> <tag_new_description>"
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

	descriptionSplit := []string{}
	for _, v := range a.Args[1:] {
		if util.IsValidUserName(v) {
			break
		}
		descriptionSplit = append(descriptionSplit, v)
	}
	description := strings.Join(descriptionSplit, " ")

	err := a.DB.ChangeDescriptionOfTag(ctx, a.ChatID, name, description)
	if err != nil {
		slog.Warn("unable to change tag description", "err", err)
		return []implementation.CommandResponse{
			resp,
		}
	}

	resp.Text = fmt.Sprintf("Changed tags description to \"%s\"", description)

	return []implementation.CommandResponse{
		resp,
	}
}
