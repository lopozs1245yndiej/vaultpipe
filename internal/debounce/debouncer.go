// Package debounce provides a debouncer that delays execution of a function
// until a specified quiet period has elapsed since the last invocation.
package debounce

import (
	"context"
	"sync"
	"time"
)

// Debouncer delays repeated calls to a function until a quiet period elapses.
type Debouncer struct {
	delay  time.Duration
	mu     sync.Mutex
	timer  *time.Timer
	cancel context.CancelFunc
}

// New creates a new Debouncer with the given delay.
// Returns an error if delay is zero or negative.
func New(delay time.Duration) (*Debouncer, error) {
	if delay <= 0 {
		return nil, ErrInvalidDelay
	}
	return &Debouncer{delay: delay}, nil
}

// Call schedules fn to run after the debounce delay.
// If Call is invoked again before the delay elapses, the previous
// scheduled call is cancelled and the timer resets.
// The provided context governs the scheduled execution; if it is
// cancelled before the timer fires, fn will not be called.
func (d *Debouncer) Call(ctx context.Context, fn func()) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.cancel != nil {
		d.cancel()
	}

	if d.timer != nil {
		d.timer.Stop()
	}

	runCtx, cancel := context.WithCancel(ctx)
	d.cancel = cancel

	d.timer = time.AfterFunc(d.delay, func() {
		select {
		case <-runCtx.Done():
			return
		default:
			fn()
		}
	})
}

// Flush cancels any pending scheduled call.
func (d *Debouncer) Flush() {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.cancel != nil {
		d.cancel()
		d.cancel = nil
	}

	if d.timer != nil {
		d.timer.Stop()
		d.timer = nil
	}
}
