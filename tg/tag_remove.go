package tg

import (
	"context"
	"fmt"
	"log/slog"
)

const removeHelpMessage string = "/tag remove <tag_name> - Remove the specified tag"

func (b *Bot) RemoveTag(ctx context.Context, chatID int64, args ...string) string {
	if len(args) != 2 {
		return removeHelpMessage
	}

	name := args[1]
	if !b.isValidTagName(name) {
		return ErrInvalidTag.Error()
	}

	err := b.tagDB.RemoveTag(ctx, chatID, name)
	if err != nil {
		slog.Warn("unable to remove tag", "err", err)
		return err.Error()
	}

	return fmt.Sprintf("Removed tag \"%s\"", name)
}
