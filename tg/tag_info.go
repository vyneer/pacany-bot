package tg

import (
	"context"
	"fmt"
	"log/slog"
	"slices"
	"strings"

	"github.com/vyneer/tg-tagbot/db"
)

const infoHelpMessage string = "/tag info <tag_name> - Get tag user count and user list"

func (b *Bot) Info(ctx context.Context, chatID int64, args ...string) (string, bool) {
	if len(args) < 2 {
		return infoHelpMessage, true
	}

	tags, err := b.tagDB.GetTags(ctx, chatID)
	if err != nil {
		slog.Warn("unable to get tags", "err", err)
		return err.Error(), true
	}

	name := args[1]
	if !b.isValidTagName(name) {
		return ErrInvalidTag.Error(), true
	}

	i := slices.IndexFunc[[]db.Tag](tags, func(t db.Tag) bool {
		return t.Name == name
	})
	if i == -1 {
		return db.ErrTagDoesntExist.Error(), true
	}

	var info []string
	fields := strings.Fields(tags[i].Mentions)
	info = append(info, fmt.Sprintf("Tag name: %s", tags[i].Name))
	info = append(info, fmt.Sprintf("User count: %d", len(fields)))
	info = append(info, "Users:")
	for _, v := range fields {
		info = append(info, fmt.Sprintf("- %s", strings.TrimPrefix(v, "@")))
	}

	return strings.Join(info, "\n"), true
}
