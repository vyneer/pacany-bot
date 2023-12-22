package tg

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
)

const listHelpMessage string = "/tag list - List all tags and their associated user count"

func (b *Bot) List(ctx context.Context, chatID int64) (string, bool) {
	tags, err := b.tagDB.GetTags(ctx, chatID)
	if err != nil {
		slog.Warn("unable to get tags", "err", err)
		return err.Error(), true
	}

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
		return "No tags in this group chat", true
	}

	return strings.Join(tagNames, "\n"), true
}
