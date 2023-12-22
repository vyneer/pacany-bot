package tg

import "strings"

var helpMessages = []string{
	"/help - Print this message",
	newHelpMessage,
	removeHelpMessage,
	addUserHelpMessage,
	removeUserHelpMessage,
	listHelpMessage,
}

func (b *Bot) SmallHelp() string {
	return "/tag <new|remove|add-user|remove-user|list> ...\nFor more information use /help"
}

func (b *Bot) Help() string {
	return strings.Join(helpMessages, "\n")
}
