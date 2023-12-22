package tg

import (
	"context"
	"fmt"
	"log/slog"
)

const addUserHelpMessage string = "/tag add-user <tag_name> <username_1> <username_2> ... <username_n> - Add specified users to an existing tag"

func (b *Bot) AddUsers(ctx context.Context, chatID int64, args ...string) string {
	if len(args) < 2 {
		return addUserHelpMessage
	}

	name := args[1]
	if !b.isValidTagName(name) {
		return ErrInvalidTag.Error()
	}
	mentions := b.filterUsernames(args[2:])
	if len(mentions) == 0 {
		return ErrNoValidUsers.Error()
	}

	err := b.tagDB.AddMentionsToTag(ctx, chatID, name, mentions...)
	if err != nil {
		slog.Warn("unable to add mentions to tag", "err", err)
		return err.Error()
	}

	return fmt.Sprintf("Added users to tag \"%s\"", name)
}
