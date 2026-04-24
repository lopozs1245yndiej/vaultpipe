package drain

import "time"

// Options holds configuration for creating a Drainer.
type Options struct {
	// Timeout is the maximum duration to wait for in-flight operations to
	// complete during Drain. Defaults to 30s if zero.
	Timeout time.Duration
}

// NewFromOptions creates a Drainer from the provided Options.
func NewFromOptions(opts Options) *Drainer {
	return New(opts.Timeout)
}
