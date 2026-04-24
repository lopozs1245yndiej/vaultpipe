package debounce

import "errors"

// ErrInvalidDelay is returned when a non-positive delay is provided to New.
var ErrInvalidDelay = errors.New("debounce: delay must be greater than zero")
