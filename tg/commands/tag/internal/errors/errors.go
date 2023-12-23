package errors

import "errors"

var (
	ErrInvalidTag   = errors.New("invalid tag: only @%#!& can be used")
	ErrNoValidUsers = errors.New("no valid users provided")
)
