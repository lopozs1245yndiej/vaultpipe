package batch

import "errors"

// Options configures a Batcher.
type Options struct {
	// Concurrency is the maximum number of items processed simultaneously.
	// Defaults to 1 if not set.
	Concurrency int
}

// NewFromOptions creates a Batcher from Options and the given ProcessFunc.
func NewFromOptions(opts Options, fn ProcessFunc) (*Batcher, error) {
	if opts.Concurrency <= 0 {
		opts.Concurrency = 1
	}
	if fn == nil {
		return nil, errors.New("batch: process function must not be nil")
	}
	return New(opts.Concurrency, fn)
}
