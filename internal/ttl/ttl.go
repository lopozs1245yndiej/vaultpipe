// Package ttl provides secret expiry tracking for vaultpipe.
// It records when a secret was last fetched and determines
// whether it has exceeded its configured time-to-live.
package ttl

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
	"time"
)

// ErrExpired is returned when a secret has exceeded its TTL.
var ErrExpired = errors.New("secret has expired")

// Entry holds the fetch timestamp for a single secret key.
type Entry struct {
	Key       string    `json:"key"`
	FetchedAt time.Time `json:"fetched_at"`
}

// Tracker records fetch times and checks expiry against a TTL.
type Tracker struct {
	mu      sync.RWMutex
	entries map[string]time.Time
	ttl     time.Duration
	path    string
}

// New creates a Tracker with the given TTL. If path is non-empty the
// tracker will persist entries to disk as a JSON file.
func New(ttl time.Duration, path string) (*Tracker, error) {
	if ttl <= 0 {
		return nil, errors.New("ttl must be positive")
	}
	t := &Tracker{
		entries: make(map[string]time.Time),
		ttl:     ttl,
		path:    path,
	}
	if path != "" {
		if err := t.load(); err != nil && !os.IsNotExist(err) {
			return nil, err
		}
	}
	return t, nil
}

// Touch records the current time as the fetch time for key.
func (t *Tracker) Touch(key string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.entries[key] = time.Now()
}

// IsExpired returns true when the key has not been touched or its
// entry is older than the configured TTL.
func (t *Tracker) IsExpired(key string) bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	at, ok := t.entries[key]
	if !ok {
		return true
	}
	return time.Since(at) > t.ttl
}

// Save persists all entries to the configured path. It is a no-op
// when no path was provided.
func (t *Tracker) Save() error {
	if t.path == "" {
		return nil
	}
	t.mu.RLock()
	defer t.mu.RUnlock()
	var entries []Entry
	for k, v := range t.entries {
		entries = append(entries, Entry{Key: k, FetchedAt: v})
	}
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(t.path, data, 0o600)
}

func (t *Tracker) load() error {
	data, err := os.ReadFile(t.path)
	if err != nil {
		return err
	}
	var entries []Entry
	if err := json.Unmarshal(data, &entries); err != nil {
		return err
	}
	for _, e := range entries {
		t.entries[e.Key] = e.FetchedAt
	}
	return nil
}
