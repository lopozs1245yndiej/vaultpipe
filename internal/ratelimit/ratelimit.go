package ratelimit

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Limiter enforces a maximum number of operations per time window.
type Limiter struct {
	mu       sync.Mutex
	rate     int
	window   time.Duration
	tokens   int
	resetAt  time.Time
}

// New creates a Limiter that allows up to rate operations per window.
// Returns an error if rate <= 0 or window <= 0.
func New(rate int, window time.Duration) (*Limiter, error) {
	if rate <= 0 {
		return nil, fmt.Errorf("ratelimit: rate must be greater than zero")
	}
	if window <= 0 {
		return nil, fmt.Errorf("ratelimit: window must be greater than zero")
	}
	return &Limiter{
		rate:    rate,
		window:  window,
		tokens:  rate,
		resetAt: time.Now().Add(window),
	}, nil
}

// Wait blocks until a token is available or the context is cancelled.
func (l *Limiter) Wait(ctx context.Context) error {
	for {
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("ratelimit: context cancelled: %w", err)
		}
		l.mu.Lock()
		now := time.Now()
		if now.After(l.resetAt) {
			l.tokens = l.rate
			l.resetAt = now.Add(l.window)
		}
		if l.tokens > 0 {
			l.tokens--
			l.mu.Unlock()
			return nil
		}
		waitUntil := l.resetAt
		l.mu.Unlock()
		select {
		case <-ctx.Done():
			return fmt.Errorf("ratelimit: context cancelled: %w", ctx.Err())
		case <-time.After(time.Until(waitUntil)):
		}
	}
}

// Remaining returns the number of tokens left in the current window.
func (l *Limiter) Remaining() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	if time.Now().After(l.resetAt) {
		return l.rate
	}
	return l.tokens
}

// Reset manually resets the token bucket to its full capacity.
func (l *Limiter) Reset() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.tokens = l.rate
	l.resetAt = time.Now().Add(l.window)
}
