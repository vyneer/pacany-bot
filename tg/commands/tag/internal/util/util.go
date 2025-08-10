//nolint:revive
package util

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	prefixRegex  *regexp.Regexp
	matcherRegex *regexp.Regexp
)

func SetTagPrefix(prefixString string) {
	prefixRegex = regexp.MustCompile(fmt.Sprintf("(?i)^[%s]{1}", prefixString))
	matcherRegex = regexp.MustCompile(fmt.Sprintf("(?i)[%s]{1}", prefixString))
}

func IsValidUserName(username string) bool {
	return strings.HasPrefix(username, "@") && len(username) > 1
}

func FilterInvalidUsernames(usernames []string) []string {
	var validUsernames []string
	usernameMap := map[string]struct{}{}

	for _, v := range usernames {
		if _, ok := usernameMap[strings.ToLower(v)]; IsValidUserName(v) && !ok {
			validUsernames = append(validUsernames, v)
			usernameMap[strings.ToLower(v)] = struct{}{}
		}
	}

	return validUsernames
}

func FilterMentions(mentions []string, ignore string) ([]string, bool) {
	var filteredMentions []string
	mentionsMap := map[string]struct{}{}

	for _, v := range mentions {
		if _, ok := mentionsMap[strings.ToLower(v)]; strings.TrimPrefix(strings.ToLower(v), "@") != strings.ToLower(ignore) && !ok {
			filteredMentions = append(filteredMentions, v)
			mentionsMap[strings.ToLower(v)] = struct{}{}
		}
	}

	return filteredMentions, len(filteredMentions) > 0
}

func IsValidTagName(name string) bool {
	return prefixRegex.MatchString(name) && len(name) > 1
}

func GetMatcherRegex() *regexp.Regexp {
	return matcherRegex
}
