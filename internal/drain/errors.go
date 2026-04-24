package drain

import "errors"

// ErrDrainTimeout is returned by Drain when the internal timeout elapses
// before all in-flight operations have completed.
var ErrDrainTimeout = errors.New("drain: timed out waiting for operations to complete")

// ErrAcquireAfterClose is returned when Acquire is called on a closed Drainer.
var ErrAcquireAfterClose = errors.New("drain: cannot acquire on a closed drainer")
