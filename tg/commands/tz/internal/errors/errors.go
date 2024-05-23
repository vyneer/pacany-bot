package errors

import "errors"

var (
	ErrInvalidTimezone = errors.New("invalid timezone: please consult https://en.wikipedia.org/wiki/List_of_tz_database_time_zones")
	ErrTimezoneNotSet  = errors.New("please set your timezone before converting")
	ErrUnableToParse   = errors.New("unable to parse provided time, valid formats are: `2006-01-02 15:04:05`, `15:04:05` and `15:04`")
	ErrInvalidUsername = errors.New("invalid username provided")
)
