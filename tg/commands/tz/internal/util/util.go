package util

import "strings"

func IsValidUserName(username string) bool {
	return strings.HasPrefix(username, "@") && len(username) > 1
}
