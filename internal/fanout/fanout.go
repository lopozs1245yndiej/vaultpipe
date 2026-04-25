// Package fanout provides a mechanism for broadcasting secrets to multiple
// destinations concurrently, collecting all errors without short-circuiting.
package fanout

import (
	"context"
	"fmt"
	"strings"
	"sync"
)

// Destination is any target that can receive a map of secrets.
type Destination interface {
	Write(ctx context.Context, secrets map[string]string) error
	Name() string
}

// Result holds the outcome of writing to a single destination.
type Result struct {
	Destination string
	Err         error
}

// Fanout broadcasts secrets to multiple destinations concurrently.
type Fanout struct {
	destinations []Destination
}

// New creates a Fanout for the given destinations.
// Returns an error if no destinations are provided.
func New(destinations ...Destination) (*Fanout, error) {
	if len(destinations) == 0 {
		return nil, fmt.Errorf("fanout: at least one destination is required")
	}
	return &Fanout{destinations: destinations}, nil
}

// Broadcast writes secrets to all destinations concurrently.
// It waits for all writes to complete and returns a slice of Results.
// The returned slice always has one entry per destination.
func (f *Fanout) Broadcast(ctx context.Context, secrets map[string]string) []Result {
	results := make([]Result, len(f.destinations))
	var wg sync.WaitGroup

	for i, dest := range f.destinations {
		wg.Add(1)
		go func(idx int, d Destination) {
			defer wg.Done()
			results[idx] = Result{
				Destination: d.Name(),
				Err:         d.Write(ctx, secrets),
			}
		}(i, dest)
	}

	wg.Wait()
	return results
}

// Errors returns a combined error if any destination failed, or nil if all succeeded.
func Errors(results []Result) error {
	var msgs []string
	for _, r := range results {
		if r.Err != nil {
			msgs = append(msgs, fmt.Sprintf("%s: %v", r.Destination, r.Err))
		}
	}
	if len(msgs) == 0 {
		return nil
	}
	return fmt.Errorf("fanout: %d destination(s) failed: %s", len(msgs), strings.Join(msgs, "; "))
}
