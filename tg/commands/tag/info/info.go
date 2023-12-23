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
	name              string = "info"
	parentName        string = "tag"
	help              string = "Get tag list, or, if a tag is specified, its user count and user list"
	helpOrder         int    = 5
	shape             string = "/taginfo or /taginfo <tag_name>"
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

	if len(a.Args) > 1 {
		resp.Text, _ = c.GetHelp()
		return resp
	}

	tags, err := a.DB.GetTags(ctx, a.ChatID)
	if err != nil {
		slog.Warn("unable to get tags", "err", err)
		resp.Text = err.Error()
		return resp
	}

	if len(a.Args) == 0 {
		var tagNames []string
		for _, v := range tags {
			l := len(strings.Fields(v.Mentions))
			tagNames = append(tagNames, fmt.Sprintf("%s - %d user%s", v.Name, l, func() string {
				if l != 1 {
					return "s"
				}
				return ""
			}()))
		}

		if len(tagNames) == 0 {
			resp.Text = "No tags in this group chat"
			return resp
		}

		resp.Text = strings.Join(tagNames, "\n")
	} else {
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
	}

	return resp
}
