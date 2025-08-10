package scan

import (
	"context"
	"log/slog"
	"regexp"
	"slices"
	"strings"

	"github.com/vyneer/pacany-bot/tg/commands/implementation"
	"github.com/vyneer/pacany-bot/tg/commands/tag/internal/util"
)

const (
	identifier string = "tag-scan"
)

type Command struct{}

func New() implementation.AutomaticCommand {
	return &Command{}
}

func (c *Command) GetIdentifier() string {
	return identifier
}

func (c *Command) GetMatcher() *regexp.Regexp {
	return util.GetMatcherRegex()
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
