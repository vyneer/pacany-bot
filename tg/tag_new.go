package tg

import (
	"context"
	"fmt"
	"log/slog"
)

const newHelpMessage string = "/tag new <tag_name> <username_1> <username_2> ... <username_n> - Add a new tag"

func (b *Bot) NewTag(ctx context.Context, chatID int64, args ...string) (string, bool) {
	if len(args) < 2 {
		return newHelpMessage, true
	}

	name := args[1]
	if !b.isValidTagName(name) {
		return ErrInvalidTag.Error(), true
	}
	mentions := b.filterInvalidUsernames(args[2:])
	if len(mentions) == 0 {
		return ErrNoValidUsers.Error(), true
	}

	err := b.tagDB.NewTag(ctx, chatID, name, mentions...)
	if err != nil {
		slog.Warn("unable to create new tag", "err", err)
		return err.Error(), true
	}

	return fmt.Sprintf("Added tag \"%s\"", name), true
}
