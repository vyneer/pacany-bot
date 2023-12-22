package tg

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/vyneer/tg-tagbot/db"
)

const removeUserHelpMessage string = "/tag remove-user <tag_name> <username_1> <username_2> ... <username_n> - Remove specified users from an existing tag"

func (b *Bot) RemoveUsers(ctx context.Context, chatID int64, args ...string) (string, bool) {
	if len(args) < 2 {
		return removeUserHelpMessage, true
	}

	name := args[1]
	if !b.isValidTagName(name) {
		return ErrInvalidTag.Error(), true
	}
	mentions := b.filterInvalidUsernames(args[2:])
	if len(mentions) == 0 {
		return ErrNoValidUsers.Error(), true
	}

	err := b.tagDB.RemoveMentionsFromTag(ctx, chatID, name, mentions...)
	if err != nil {
		if errors.Is(err, db.ErrEmptyTag) {
			err := b.tagDB.RemoveTag(ctx, chatID, name)
			if err != nil {
				slog.Warn("unable to remove tag", "err", err)
				return err.Error(), true
			}
			return fmt.Sprintf("Removed tag \"%s\"", name), true
		}
		slog.Warn("unable to remove users from tag", "err", err)
		return err.Error(), true
	}

	return fmt.Sprintf("Removed users from tag \"%s\"", name), true
}
