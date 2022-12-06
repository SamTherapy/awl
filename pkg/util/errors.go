package util

import (
	"errors"
	"fmt"
)

// ErrHTTPStatus is returned when DoH returns a bad status code.
type ErrHTTPStatus struct {
	// Status code
	Code int
}

func (e *ErrHTTPStatus) Error() string {
	return fmt.Sprintf("doh server responded with HTTP %d", e.Code)
}

// ErrNotError is an error that is not actually an error.
var ErrNotError = errors.New("not an error")
