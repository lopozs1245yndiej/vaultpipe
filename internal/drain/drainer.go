// Package drain provides graceful shutdown coordination for vaultpipe,
// ensuring in-flight sync operations complete before the process exits.
package drain

import (
	"context"
	"sync"
	"time"
)

// Drainer tracks active operations and blocks shutdown until all complete
// or the drain timeout is exceeded.
type Drainer struct {
	mu      sync.Mutex
	wg      sync.WaitGroup
	closed  bool
	timeout time.Duration
}

// New creates a new Drainer with the given drain timeout.
// If timeout is zero, it defaults to 30 seconds.
func New(timeout time.Duration) *Drainer {
	if timeout <= 0 {
		timeout = 30 * time.Second
	}
	return &Drainer{timeout: timeout}
}

// Acquire registers a new in-flight operation. It returns false if the
// Drainer has already been closed and no new work should be accepted.
func (d *Drainer) Acquire() bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.closed {
		return false
	}
	d.wg.Add(1)
	return true
}

// Release marks one in-flight operation as complete.
func (d *Drainer) Release() {
	d.wg.Done()
}

// Drain closes the Drainer to new acquisitions and waits for all active
// operations to finish, or until the drain timeout elapses.
// It returns ctx.Err() if the parent context is cancelled before draining
// completes, or ErrDrainTimeout if the internal timeout fires first.
func (d *Drainer) Drain(ctx context.Context) error {
	d.mu.Lock()
	d.closed = true
	d.mu.Unlock()

	done := make(chan struct{})
	go func() {
		d.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(d.timeout):
		return ErrDrainTimeout
	}
}

// IsClosed reports whether the Drainer has been closed to new acquisitions.
func (d *Drainer) IsClosed() bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.closed
}
