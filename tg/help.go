package tg

import "strings"

var helpMessages = []string{
	"/help - Print this message",
	newHelpMessage,
	removeHelpMessage,
	addUserHelpMessage,
	removeUserHelpMessage,
	infoHelpMessage,
	listHelpMessage,
}

func (b *Bot) SmallHelp() (string, bool) {
	return "/tag <new|remove|add-user|remove-user|info|list> ...\n\nFor more information use /help", true
}

func (b *Bot) Help() (string, bool) {
	return strings.Join(helpMessages, "\n\n"), true
}
