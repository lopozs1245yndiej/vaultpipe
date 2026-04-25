// Package semaphore provides a counting semaphore for bounding concurrent
// operations such as parallel Vault secret fetches.
package semaphore

import (
	"context"
	"errors"
	"time"
)

// ErrAcquireTimeout is returned when Acquire times out waiting for a slot.
var ErrAcquireTimeout = errors.New("semaphore: acquire timed out")

// Semaphore is a counting semaphore backed by a buffered channel.
type Semaphore struct {
	slots   chan struct{}
	timeout time.Duration
}

// New creates a Semaphore with the given concurrency limit and acquire timeout.
// limit must be >= 1. A zero timeout means no timeout (block until context is
// cancelled or a slot becomes available).
func New(limit int, timeout time.Duration) (*Semaphore, error) {
	if limit < 1 {
		return nil, errors.New("semaphore: limit must be >= 1")
	}
	return &Semaphore{
		slots:   make(chan struct{}, limit),
		timeout: timeout,
	}, nil
}

// Acquire claims one slot. It blocks until a slot is available, the context is
// cancelled, or the configured timeout elapses.
func (s *Semaphore) Acquire(ctx context.Context) error {
	if s.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, s.timeout)
		defer cancel()
	}
	select {
	case s.slots <- struct{}{}:
		return nil
	case <-ctx.Done():
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return ErrAcquireTimeout
		}
		return ctx.Err()
	}
}

// Release frees one slot. It panics if Release is called more times than
// Acquire, which indicates a programming error.
func (s *Semaphore) Release() {
	select {
	case <-s.slots:
	default:
		panic("semaphore: Release called without matching Acquire")
	}
}

// Available returns the number of free slots at the moment of the call.
func (s *Semaphore) Available() int {
	return cap(s.slots) - len(s.slots)
}
