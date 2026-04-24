// Package circuit implements a circuit breaker for Vault API calls.
// It tracks consecutive failures and opens the circuit after a configurable
// threshold, preventing cascading failures during Vault outages.
package circuit

import (
	"errors"
	"sync"
	"time"
)

// ErrOpen is returned when the circuit breaker is in the open state.
var ErrOpen = errors.New("circuit breaker is open")

// State represents the current state of the circuit breaker.
type State int

const (
	StateClosed State = iota
	StateOpen
	StateHalfOpen
)

// Breaker is a circuit breaker that tracks failures and opens after a threshold.
type Breaker struct {
	mu           sync.Mutex
	failures     int
	maxFailures  int
	resetTimeout time.Duration
	state        State
	openedAt     time.Time
}

// New creates a new Breaker with the given failure threshold and reset timeout.
func New(maxFailures int, resetTimeout time.Duration) (*Breaker, error) {
	if maxFailures <= 0 {
		return nil, errors.New("maxFailures must be greater than zero")
	}
	if resetTimeout <= 0 {
		return nil, errors.New("resetTimeout must be greater than zero")
	}
	return &Breaker{
		maxFailures:  maxFailures,
		resetTimeout: resetTimeout,
		state:        StateClosed,
	}, nil
}

// Allow reports whether a call should be allowed through.
// It returns ErrOpen if the circuit is open and the reset timeout has not elapsed.
func (b *Breaker) Allow() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	switch b.state {
	case StateOpen:
		if time.Since(b.openedAt) >= b.resetTimeout {
			b.state = StateHalfOpen
			return nil
		}
		return ErrOpen
	default:
		return nil
	}
}

// RecordSuccess resets the failure count and closes the circuit.
func (b *Breaker) RecordSuccess() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failures = 0
	b.state = StateClosed
}

// RecordFailure increments the failure count and opens the circuit if the threshold is reached.
func (b *Breaker) RecordFailure() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failures++
	if b.failures >= b.maxFailures {
		b.state = StateOpen
		b.openedAt = time.Now()
	}
}

// State returns the current state of the circuit breaker.
func (b *Breaker) State() State {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.state
}

// Failures returns the current failure count.
func (b *Breaker) Failures() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.failures
}
