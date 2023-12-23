package util

import (
	"regexp"
	"strings"
)

var prefixRegex = regexp.MustCompile(`(?i)^[@%#!&]{1}`)

func IsValidUserName(username string) bool {
	return strings.HasPrefix(username, "@") && len(username) > 1
}

func FilterInvalidUsernames(usernames []string) []string {
	var validUsernames []string

	for _, v := range usernames {
		if IsValidUserName(v) {
			validUsernames = append(validUsernames, v)
		}
	}

	return validUsernames
}

func FilterMentions(mentions []string, ignore string) (string, bool) {
	var filteredMentions []string

	for _, v := range mentions {
		if strings.TrimPrefix(v, "@") != ignore {
			filteredMentions = append(filteredMentions, v)
		}
	}

	return strings.Join(filteredMentions, " "), len(filteredMentions) > 0
}

func IsValidTagName(name string) bool {
	return prefixRegex.MatchString(name) && len(name) > 1
}
