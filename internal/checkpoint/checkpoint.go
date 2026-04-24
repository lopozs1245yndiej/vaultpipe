// Package checkpoint tracks the last successful sync time and metadata
// for each configured secret path, enabling incremental and resumable syncs.
package checkpoint

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
	"time"
)

// ErrNotFound is returned when no checkpoint exists for a given key.
var ErrNotFound = errors.New("checkpoint: no entry found")

// Entry holds the recorded state for a single secret path.
type Entry struct {
	Path      string    `json:"path"`
	SyncedAt  time.Time `json:"synced_at"`
	Checksum  string    `json:"checksum"`
	Namespace string    `json:"namespace,omitempty"`
}

// Tracker persists and retrieves checkpoint entries.
type Tracker struct {
	mu       sync.RWMutex
	filePath string
	entries  map[string]Entry
}

// New loads an existing checkpoint file or creates an empty tracker.
func New(filePath string) (*Tracker, error) {
	t := &Tracker{
		filePath: filePath,
		entries:  make(map[string]Entry),
	}
	data, err := os.ReadFile(filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return t, nil
		}
		return nil, err
	}
	if err := json.Unmarshal(data, &t.entries); err != nil {
		return nil, err
	}
	return t, nil
}

// Set records or updates the checkpoint entry for the given path.
func (t *Tracker) Set(e Entry) error {
	if e.Path == "" {
		return errors.New("checkpoint: path must not be empty")
	}
	if e.SyncedAt.IsZero() {
		e.SyncedAt = time.Now().UTC()
	}
	t.mu.Lock()
	t.entries[e.Path] = e
	t.mu.Unlock()
	return t.save()
}

// Get retrieves the checkpoint entry for the given path.
func (t *Tracker) Get(path string) (Entry, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	e, ok := t.entries[path]
	if !ok {
		return Entry{}, ErrNotFound
	}
	return e, nil
}

// Delete removes the checkpoint for the given path and persists the change.
func (t *Tracker) Delete(path string) error {
	t.mu.Lock()
	delete(t.entries, path)
	t.mu.Unlock()
	return t.save()
}

// All returns a copy of all current checkpoint entries.
func (t *Tracker) All() []Entry {
	t.mu.RLock()
	defer t.mu.RUnlock()
	out := make([]Entry, 0, len(t.entries))
	for _, e := range t.entries {
		out = append(out, e)
	}
	return out
}

func (t *Tracker) save() error {
	t.mu.RLock()
	data, err := json.MarshalIndent(t.entries, "", "  ")
	t.mu.RUnlock()
	if err != nil {
		return err
	}
	return os.WriteFile(t.filePath, data, 0600)
}
