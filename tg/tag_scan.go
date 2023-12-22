package tg

import (
	"context"
	"log/slog"
	"slices"
	"strings"
)

func (b *Bot) Scan(ctx context.Context, chatID int64, username, text string) (string, bool) {
	tags, err := b.tagDB.GetTags(ctx, chatID)
	if err != nil {
		slog.Warn("unable to get tags", "err", err)
		return err.Error(), true
	}

	fields := strings.Fields(text)

	for _, v := range tags {
		if slices.Contains[[]string](fields, v.Name) {
			if filtered, ok := b.filterMentions(v.Mentions, username); ok {
				return filtered, false
			}
			return "You're the only person using this tag", true
		}
	}

	return "", false
}
