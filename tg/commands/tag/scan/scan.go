package scan

import (
	"context"
	"fmt"
	"log/slog"
	"slices"
	"strings"

	"github.com/vyneer/pacany-bot/tg/commands/implementation"
	"github.com/vyneer/pacany-bot/tg/commands/tag/internal/util"
)

const (
	name              string = "scan"
	parentName        string = "tag"
	help              string = ""
	arguments         string = ""
	showInCommandList bool   = false
	showInHelp        bool   = false
	adminOnly         bool   = false
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
	tags, err := a.DB.GetTags(ctx, a.ChatID)
	if err != nil {
		slog.Warn("unable to get tags", "err", err)
		return []implementation.CommandResponse{
			{
				Text:       err.Error(),
				Reply:      true,
				Capitalize: true,
			},
		}
	}

	username := a.Args[0]
	text := a.Args[1]

	fields := strings.Fields(text)

	tagCount := 0
	var allMentions []string
	description := ""
	for _, v := range tags {
		if i := slices.Index(fields, v.Name); i != -1 {
			if tagCount == 0 {
				description = v.Description
			}
			tagCount++
			allMentions = append(allMentions, strings.Fields(v.Mentions)...)
			fields = append(fields[:i], fields[i+1:]...)
		}
	}

	if len(allMentions) > 0 {
		if filtered, ok := util.FilterMentions(allMentions, username); ok {
			msgs := []implementation.CommandResponse{
				{
					Text:       strings.Join(filtered, " "),
					Reply:      false,
					Capitalize: false,
				},
			}

			if tagCount == 1 {
				msgs = splitMentionsIntoResponses(filtered, description)
			}

			if len(fields) > 0 {
				msgs = splitMentionsIntoResponses(filtered, strings.Join(fields, " "))
			}

			return msgs
		}

		return []implementation.CommandResponse{
			{
				Text:       "You're the only person using this tag",
				Reply:      true,
				Capitalize: true,
			},
		}
	}

	return []implementation.CommandResponse{}
}

func splitMentionsIntoResponses(mentions []string, desc string) []implementation.CommandResponse {
	var msgs []implementation.CommandResponse

	for i := 0; i < len(mentions); i += 5 {
		if i+5 > len(mentions) {
			msgs = append(msgs, implementation.CommandResponse{
				Text:       strings.Join([]string{desc, strings.Join(mentions[i:], " ")}, " "),
				Reply:      false,
				Capitalize: false,
			})
			continue
		}

		msgs = append(msgs, implementation.CommandResponse{
			Text:       strings.Join([]string{desc, strings.Join(mentions[i:i+5], " ")}, " "),
			Reply:      false,
			Capitalize: false,
		})
	}

	return msgs
}
