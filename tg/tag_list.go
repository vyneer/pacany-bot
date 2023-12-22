package tg

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
)

const listHelpMessage string = "/tag list - List all tags and their associated user count"

func (b *Bot) List(ctx context.Context, chatID int64) string {
	tags, err := b.tagDB.GetTags(ctx, chatID)
	if err != nil {
		slog.Warn("unable to get tags", "err", err)
		return err.Error()
	}

	var tagNames []string
	for _, v := range tags {
		tagNames = append(tagNames, fmt.Sprintf("%s - %d users", v.Name, len(strings.Fields(v.Mentions))))
	}

	if len(tagNames) == 0 {
		return "No tags in this group chat"
	}

	return strings.Join(tagNames, "\n")
}
