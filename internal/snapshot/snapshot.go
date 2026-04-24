package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Snapshot represents a point-in-time capture of secrets.
type Snapshot struct {
	ID        string            `json:"id"`
	CreatedAt time.Time         `json:"created_at"`
	Secrets   map[string]string `json:"secrets"`
	Source    string            `json:"source"`
}

// Manager handles saving and loading snapshots to disk.
type Manager struct {
	dir string
}

// New creates a new snapshot Manager that stores snapshots in dir.
func New(dir string) (*Manager, error) {
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return nil, fmt.Errorf("snapshot: create dir: %w", err)
	}
	return &Manager{dir: dir}, nil
}

// Save writes secrets to a new snapshot file and returns the snapshot.
func (m *Manager) Save(source string, secrets map[string]string) (*Snapshot, error) {
	snap := &Snapshot{
		ID:        fmt.Sprintf("%d", time.Now().UnixNano()),
		CreatedAt: time.Now().UTC(),
		Secrets:   secrets,
		Source:    source,
	}
	path := filepath.Join(m.dir, snap.ID+".json")
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	if err != nil {
		return nil, fmt.Errorf("snapshot: open file: %w", err)
	}
	defer f.Close()
	if err := json.NewEncoder(f).Encode(snap); err != nil {
		return nil, fmt.Errorf("snapshot: encode: %w", err)
	}
	return snap, nil
}

// Load reads a snapshot by ID from disk.
func (m *Manager) Load(id string) (*Snapshot, error) {
	path := filepath.Join(m.dir, id+".json")
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("snapshot: open: %w", err)
	}
	defer f.Close()
	var snap Snapshot
	if err := json.NewDecoder(f).Decode(&snap); err != nil {
		return nil, fmt.Errorf("snapshot: decode: %w", err)
	}
	return &snap, nil
}

// List returns all snapshot IDs stored in the directory, sorted oldest-first.
func (m *Manager) List() ([]string, error) {
	entries, err := filepath.Glob(filepath.Join(m.dir, "*.json"))
	if err != nil {
		return nil, fmt.Errorf("snapshot: list: %w", err)
	}
	ids := make([]string, 0, len(entries))
	for _, e := range entries {
		base := filepath.Base(e)
		ids = append(ids, base[:len(base)-5])
	}
	return ids, nil
}

// Delete removes a snapshot by ID.
func (m *Manager) Delete(id string) error {
	path := filepath.Join(m.dir, id+".json")
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("snapshot: delete: %w", err)
	}
	return nil
}
