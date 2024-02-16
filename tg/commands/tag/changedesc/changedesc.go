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
	help              string = "Change the description of a tag"
	shape             string = "/tagchangedesc <tag_name> <tag_new_description>"
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
	return fmt.Sprintf("%s - %s", shape, help), showInHelp
}

func (c *Command) GetDescription() (string, bool) {
	return fmt.Sprintf("%s - %s", help, shape), showInCommandList
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
		resp.Text = err.Error()
		return resp
	}

	resp.Text = fmt.Sprintf("Changed tags description to \"%s\"", description)

	return resp
}
