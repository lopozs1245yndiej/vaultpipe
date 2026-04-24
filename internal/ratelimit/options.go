package ratelimit

import (
	"fmt"
	"time"
)

// Options holds configuration for creating a Limiter via NewFromOptions.
type Options struct {
	// Rate is the maximum number of allowed operations per Window.
	Rate int
	// Window is the duration of each rate-limit window.
	Window time.Duration
}

// NewFromOptions creates a Limiter from an Options struct.
// Returns an error if Rate or Window are invalid.
func NewFromOptions(opts Options) (*Limiter, error) {
	if opts.Rate <= 0 {
		return nil, fmt.Errorf("ratelimit: options.Rate must be greater than zero")
	}
	if opts.Window <= 0 {
		return nil, fmt.Errorf("ratelimit: options.Window must be greater than zero")
	}
	return New(opts.Rate, opts.Window)
}
