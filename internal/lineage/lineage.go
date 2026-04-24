// Package lineage tracks the origin and history of secrets synced from Vault,
// recording which path, version, and timestamp each secret was sourced from.
package lineage

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

// Record holds provenance information for a single secret key.
type Record struct {
	Key       string    `json:"key"`
	VaultPath string    `json:"vault_path"`
	Version   int       `json:"version"`
	SyncedAt  time.Time `json:"synced_at"`
}

// Tracker maintains an in-memory map of secret lineage records and can
// persist them to a JSON file for auditing purposes.
type Tracker struct {
	mu      sync.RWMutex
	records map[string]Record
	path    string
}

// New creates a new Tracker that persists lineage data to the given file path.
// If the file already exists its records are loaded into memory.
func New(path string) (*Tracker, error) {
	t := &Tracker{
		records: make(map[string]Record),
		path:    path,
	}
	if err := t.load(); err != nil {
		return nil, fmt.Errorf("lineage: load: %w", err)
	}
	return t, nil
}

// Track records the provenance of a secret key.
func (t *Tracker) Track(key, vaultPath string, version int, syncedAt time.Time) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.records[key] = Record{
		Key:       key,
		VaultPath: vaultPath,
		Version:   version,
		SyncedAt:  syncedAt,
	}
}

// Get returns the lineage record for the given key, if present.
func (t *Tracker) Get(key string) (Record, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	r, ok := t.records[key]
	return r, ok
}

// All returns a copy of all tracked records.
func (t *Tracker) All() []Record {
	t.mu.RLock()
	defer t.mu.RUnlock()
	out := make([]Record, 0, len(t.records))
	for _, r := range t.records {
		out = append(out, r)
	}
	return out
}

// Save persists the current records to the configured file path as JSON.
func (t *Tracker) Save() error {
	t.mu.RLock()
	defer t.mu.RUnlock()
	f, err := os.Create(t.path)
	if err != nil {
		return fmt.Errorf("lineage: save: %w", err)
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(t.records)
}

func (t *Tracker) load() error {
	f, err := os.Open(t.path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewDecoder(f).Decode(&t.records)
}
