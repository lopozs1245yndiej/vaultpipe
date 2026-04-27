package batch

import (
	"context"
	"errors"
	"sync"
)

// ErrEmptyBatch is returned when no items are provided.
var ErrEmptyBatch = errors.New("batch: no items provided")

// ProcessFunc is a function that processes a single map of secrets.
type ProcessFunc func(ctx context.Context, secrets map[string]string) error

// Batcher runs a ProcessFunc against multiple secret maps concurrently,
// up to a configurable concurrency limit.
type Batcher struct {
	concurrency int
	fn          ProcessFunc
}

// New creates a Batcher with the given concurrency limit and processing function.
// concurrency must be >= 1.
func New(concurrency int, fn ProcessFunc) (*Batcher, error) {
	if concurrency < 1 {
		return nil, errors.New("batch: concurrency must be at least 1")
	}
	if fn == nil {
		return nil, errors.New("batch: process function must not be nil")
	}
	return &Batcher{concurrency: concurrency, fn: fn}, nil
}

// Result holds the outcome of processing a single item.
type Result struct {
	Index int
	Err   error
}

// Run processes all items concurrently and returns a slice of Results.
// If ctx is cancelled, in-flight work is abandoned and remaining items are skipped.
func (b *Batcher) Run(ctx context.Context, items []map[string]string) ([]Result, error) {
	if len(items) == 0 {
		return nil, ErrEmptyBatch
	}

	results := make([]Result, len(items))
	sem := make(chan struct{}, b.concurrency)
	var wg sync.WaitGroup

	for i, item := range items {
		select {
		case <-ctx.Done():
			results[i] = Result{Index: i, Err: ctx.Err()}
			continue
		case sem <- struct{}{}:
		}

		wg.Add(1)
		go func(idx int, secrets map[string]string) {
			defer wg.Done()
			defer func() { <-sem }()
			err := b.fn(ctx, secrets)
			results[idx] = Result{Index: idx, Err: err}
		}(i, item)
	}

	wg.Wait()
	return results, nil
}

// Errors returns only the Results that contain a non-nil error.
func Errors(results []Result) []Result {
	var errs []Result
	for _, r := range results {
		if r.Err != nil {
			errs = append(errs, r)
		}
	}
	return errs
}
