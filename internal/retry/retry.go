// Package retry provides configurable retry logic with exponential backoff
// for transient failures when communicating with HashiCorp Vault.
package retry

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"
)

// ErrMaxAttemptsReached is returned when all retry attempts are exhausted.
var ErrMaxAttemptsReached = errors.New("retry: max attempts reached")

// Retrier executes a function with retry logic and exponential backoff.
type Retrier struct {
	maxAttempts int
	baseDelay   time.Duration
	maxDelay    time.Duration
	multiplier  float64
	retryIf     func(error) bool
}

// Option configures a Retrier.
type Option func(*Retrier)

// WithMaxAttempts sets the maximum number of attempts (including the first call).
func WithMaxAttempts(n int) Option {
	return func(r *Retrier) {
		if n > 0 {
			r.maxAttempts = n
		}
	}
}

// WithBaseDelay sets the initial delay between retries.
func WithBaseDelay(d time.Duration) Option {
	return func(r *Retrier) {
		if d > 0 {
			r.baseDelay = d
		}
	}
}

// WithMaxDelay caps the delay between retries.
func WithMaxDelay(d time.Duration) Option {
	return func(r *Retrier) {
		if d > 0 {
			r.maxDelay = d
		}
	}
}

// WithMultiplier sets the exponential backoff multiplier.
func WithMultiplier(m float64) Option {
	return func(r *Retrier) {
		if m > 1 {
			r.multiplier = m
		}
	}
}

// WithRetryIf sets a predicate that determines whether an error should trigger
// a retry. By default, all non-nil errors are retried.
func WithRetryIf(fn func(error) bool) Option {
	return func(r *Retrier) {
		if fn != nil {
			r.retryIf = fn
		}
	}
}

// New creates a Retrier with sensible defaults:
// 3 max attempts, 200ms base delay, 10s max delay, multiplier 2.0.
func New(opts ...Option) *Retrier {
	r := &Retrier{
		maxAttempts: 3,
		baseDelay:   200 * time.Millisecond,
		maxDelay:    10 * time.Second,
		multiplier:  2.0,
		retryIf:     func(err error) bool { return err != nil },
	}
	for _, o := range opts {
		o(r)
	}
	return r
}

// Do executes fn, retrying on transient errors according to the configured
// policy. The context is checked before each sleep so cancellation is prompt.
func (r *Retrier) Do(ctx context.Context, fn func() error) error {
	var lastErr error
	for attempt := 0; attempt < r.maxAttempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("retry: context cancelled before attempt %d: %w", attempt+1, err)
		}

		lastErr = fn()
		if lastErr == nil {
			return nil
		}

		if !r.retryIf(lastErr) {
			return lastErr
		}

		if attempt == r.maxAttempts-1 {
			break
		}

		delay := r.delayFor(attempt)
		select {
		case <-ctx.Done():
			return fmt.Errorf("retry: context cancelled during backoff: %w", ctx.Err())
		case <-time.After(delay):
		}
	}
	return fmt.Errorf("%w after %d attempts: %v", ErrMaxAttemptsReached, r.maxAttempts, lastErr)
}

// delayFor calculates the backoff duration for a given attempt index (0-based).
func (r *Retrier) delayFor(attempt int) time.Duration {
	delay := float64(r.baseDelay) * math.Pow(r.multiplier, float64(attempt))
	if delay > float64(r.maxDelay) {
		delay = float64(r.maxDelay)
	}
	return time.Duration(delay)
}
