package errors

import "errors"

var (
	ErrBadType       = errors.New("bad type error")
	ErrNotFound      = errors.New("not found")
	ErrUnknownAction = errors.New("unknown action")
	ErrTooManyArgs   = errors.New("too many arguments")
)
