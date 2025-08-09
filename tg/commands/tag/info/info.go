package info

import (
	"context"
	"fmt"
	"log/slog"
	"slices"
	"strings"

	"github.com/vyneer/pacany-bot/db"
	"github.com/vyneer/pacany-bot/tg/commands/implementation"
	tag_errors "github.com/vyneer/pacany-bot/tg/commands/tag/internal/errors"
	"github.com/vyneer/pacany-bot/tg/commands/tag/internal/util"
)

const (
	name              string = "info"
	parentName        string = "tag"
	help              string = "Get tag list, or, if a tag is specified, its description, user count and user list"
	arguments         string = "[tag_name]"
	showInCommandList bool   = true
	showInHelp        bool   = true
	adminOnly         bool   = false
)

type Command struct{}

func New() implementation.InteractableCommand {
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

	if len(a.Args) > 1 {
		resp.Text, _ = c.GetHelp()
		return []implementation.CommandResponse{
			resp,
		}
	}

	tags, err := a.DB.GetTags(ctx, a.ChatID)
	if err != nil {
		slog.Warn("unable to get tags", "err", err)
		resp.Text = err.Error()
		return []implementation.CommandResponse{
			resp,
		}
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
			return []implementation.CommandResponse{
				resp,
			}
		}

		resp.Text = strings.Join(tagNames, "\n")
	} else {
		name := a.Args[0]
		if !util.IsValidTagName(name) {
			resp.Text = tag_errors.ErrInvalidTag.Error()
			return []implementation.CommandResponse{
				resp,
			}
		}

		i := slices.IndexFunc(tags, func(t db.Tag) bool {
			return t.Name == name
		})
		if i == -1 {
			resp.Text = db.ErrTagDoesntExist.Error()
			return []implementation.CommandResponse{
				resp,
			}
		}

		var info []string
		fields := strings.Fields(tags[i].Mentions)
		info = append(info, fmt.Sprintf("Tag name: %s", tags[i].Name))
		if tags[i].Description != "" {
			info = append(info, fmt.Sprintf("Tag description: %s", tags[i].Description))
		}
		info = append(info, fmt.Sprintf("User count: %d", len(fields)))
		info = append(info, "User list:")
		for _, v := range fields {
			info = append(info, fmt.Sprintf("- %s", strings.TrimPrefix(v, "@")))
		}

		resp.Text = strings.Join(info, "\n")
	}

	return []implementation.CommandResponse{
		resp,
	}
}
