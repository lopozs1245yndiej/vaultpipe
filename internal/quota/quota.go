// Package quota provides a simple per-key read quota tracker that limits
// how many times a secret path may be fetched within a rolling time window.
package quota

import (
	"fmt"
	"sync"
	"time"
)

// ErrQuotaExceeded is returned when a key has exceeded its allowed reads.
type ErrQuotaExceeded struct {
	Key   string
	Limit int
}

func (e *ErrQuotaExceeded) Error() string {
	return fmt.Sprintf("quota exceeded for key %q: limit is %d reads per window", e.Key, e.Limit)
}

type entry struct {
	count     int
	windowEnd time.Time
}

// Tracker enforces per-key read quotas within a rolling time window.
type Tracker struct {
	mu     sync.Mutex
	store  map[string]*entry
	limit  int
	window time.Duration
}

// New creates a Tracker that allows at most limit reads per key within window.
// Returns an error if limit < 1 or window <= 0.
func New(limit int, window time.Duration) (*Tracker, error) {
	if limit < 1 {
		return nil, fmt.Errorf("quota: limit must be at least 1, got %d", limit)
	}
	if window <= 0 {
		return nil, fmt.Errorf("quota: window must be positive, got %v", window)
	}
	return &Tracker{
		store:  make(map[string]*entry),
		limit:  limit,
		window: window,
	}, nil
}

// Check returns nil if the key is within quota, or *ErrQuotaExceeded otherwise.
// Each successful call increments the counter for the key.
func (t *Tracker) Check(key string) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := time.Now()
	e, ok := t.store[key]
	if !ok || now.After(e.windowEnd) {
		t.store[key] = &entry{count: 1, windowEnd: now.Add(t.window)}
		return nil
	}
	if e.count >= t.limit {
		return &ErrQuotaExceeded{Key: key, Limit: t.limit}
	}
	e.count++
	return nil
}

// Reset clears the quota counter for a specific key.
func (t *Tracker) Reset(key string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.store, key)
}

// Flush clears all quota counters.
func (t *Tracker) Flush() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.store = make(map[string]*entry)
}

// Remaining returns how many reads are left for key in the current window.
// If the window has expired or the key is unseen, the full limit is returned.
func (t *Tracker) Remaining(key string) int {
	t.mu.Lock()
	defer t.mu.Unlock()
	e, ok := t.store[key]
	if !ok || time.Now().After(e.windowEnd) {
		return t.limit
	}
	rem := t.limit - e.count
	if rem < 0 {
		return 0
	}
	return rem
}
