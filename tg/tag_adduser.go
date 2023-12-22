package tg

import (
	"context"
	"fmt"
	"log/slog"
)

const addUserHelpMessage string = "/tag add-user <tag_name> <username_1> <username_2> ... <username_n> - Add specified users to an existing tag"

func (b *Bot) AddUsers(ctx context.Context, chatID int64, args ...string) (string, bool) {
	if len(args) < 2 {
		return addUserHelpMessage, true
	}

	name := args[1]
	if !b.isValidTagName(name) {
		return ErrInvalidTag.Error(), true
	}
	mentions := b.filterInvalidUsernames(args[2:])
	if len(mentions) == 0 {
		return ErrNoValidUsers.Error(), true
	}

	err := b.tagDB.AddMentionsToTag(ctx, chatID, name, mentions...)
	if err != nil {
		slog.Warn("unable to add mentions to tag", "err", err)
		return err.Error(), true
	}

	return fmt.Sprintf("Added users to tag \"%s\"", name), true
}
