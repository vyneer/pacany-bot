package errors

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidTag   error
	ErrNoValidUsers = errors.New("no valid users provided")
)

func SetErrInvalidTag(prefixString string) {
	ErrInvalidTag = fmt.Errorf("invalid tag: only %s can be used", prefixString)
}
