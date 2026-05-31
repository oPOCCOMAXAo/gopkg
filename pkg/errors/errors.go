package errors

import "errors"

// Common error values for the library.
var (
	ErrDuplicate = errors.New("duplicate")
	ErrNotFound  = errors.New("not found")
)
