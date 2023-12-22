package tg

import (
	"context"
	"log/slog"
	"slices"
	"strings"
)

func (b *Bot) Scan(ctx context.Context, chatID int64, text string) string {
	tags, err := b.tagDB.GetTags(ctx, chatID)
	if err != nil {
		slog.Warn("unable to get tags", "err", err)
		return err.Error()
	}

	fields := strings.Fields(text)

	for _, v := range tags {
		if slices.Contains[[]string](fields, v.Name) {
			return v.Mentions
		}
	}

	return ""
}
