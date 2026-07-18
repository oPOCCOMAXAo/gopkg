package errors

import "errors"

var (
	ErrDuplicate      = errors.New("duplicate")
	ErrFailed         = errors.New("failed")
	ErrInvalidAuth    = errors.New("invalid auth")
	ErrInvalidParam   = errors.New("invalid parameter")
	ErrLimitExceeded  = errors.New("limit exceeded")
	ErrNoAccess       = errors.New("no access")
	ErrNotFound       = errors.New("not found")
	ErrNotImplemented = errors.New("not implemented")
	ErrRequestFailed  = errors.New("request failed")
	ErrRetryLater     = errors.New("retry later")
)
