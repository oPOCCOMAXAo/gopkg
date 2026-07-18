package errors

import (
	"fmt"
	"strconv"
	"unsafe"
)

// Join returns an error that wraps the given errors.
// Any nil error values are discarded.
// Join returns nil if every value in errs is nil.
// The error formats as the concatenation of the strings obtained
// by calling the Error method of each element of errs, with a newline
// between each string.
//
// A non-nil error returned by Join implements the Unwrap() []error method.
// The errors may be inspected with [errors.Is] and [errors.As].
func Join(errs ...error) error {
	n := 0

	for _, err := range errs {
		if err != nil {
			n++
		}
	}

	if n == 0 {
		return nil
	}

	e := &joinError{
		errs: make([]error, 0, n),
	}
	for _, err := range errs {
		if err != nil {
			e.errs = append(e.errs, err)
		}
	}

	return e
}

type joinError struct {
	errs []error
}

func (e *joinError) Error() string {
	// Since Join returns nil if every value in errs is nil,
	// e.errs cannot be empty.
	if len(e.errs) == 1 {
		return e.errs[0].Error()
	}

	b := []byte(e.errs[0].Error())
	for _, err := range e.errs[1:] {
		b = append(b, '\n')
		b = append(b, err.Error()...)
	}

	// At this point, b has at least one byte '\n'.
	return unsafe.String(&b[0], len(b))
}

func (e *joinError) Unwrap() []error {
	return e.errs
}

func (e *joinError) Format(s fmt.State, verb rune) {
	format := fmt.FormatString(s, verb)

	_, _ = s.Write([]byte("errors.Join:\n"))

	for i, err := range e.errs {
		_, _ = s.Write([]byte("error[" + strconv.Itoa(i) + "]:\n"))

		formatter, ok := err.(fmt.Formatter)
		if ok {
			formatter.Format(s, verb)
		} else {
			fmt.Fprintf(s, format, err)
		}

		_, _ = s.Write([]byte("\n\n"))
	}
}
