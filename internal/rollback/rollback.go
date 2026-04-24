package rollback

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Rollbacker restores a .env file from a previous snapshot backup.
type Rollbacker struct {
	backupDir string
}

// Entry represents a restorable backup entry.
type Entry struct {
	ID        string
	Timestamp time.Time
	Path      string
}

// New creates a new Rollbacker that reads backups from backupDir.
func New(backupDir string) (*Rollbacker, error) {
	if backupDir == "" {
		return nil, fmt.Errorf("rollback: backup directory must not be empty")
	}
	if err := os.MkdirAll(backupDir, 0o700); err != nil {
		return nil, fmt.Errorf("rollback: create backup dir: %w", err)
	}
	return &Rollbacker{backupDir: backupDir}, nil
}

// List returns all available backup entries sorted by filename (chronological).
func (r *Rollbacker) List() ([]Entry, error) {
	matches, err := filepath.Glob(filepath.Join(r.backupDir, "*.env.bak"))
	if err != nil {
		return nil, fmt.Errorf("rollback: list backups: %w", err)
	}
	entries := make([]Entry, 0, len(matches))
	for _, p := range matches {
		info, err := os.Stat(p)
		if err != nil {
			continue
		}
		entries = append(entries, Entry{
			ID:        info.Name(),
			Timestamp: info.ModTime(),
			Path:      p,
		})
	}
	return entries, nil
}

// Restore copies the backup identified by id to dest, overwriting it.
func (r *Rollbacker) Restore(id, dest string) error {
	if id == "" {
		return fmt.Errorf("rollback: id must not be empty")
	}
	if dest == "" {
		return fmt.Errorf("rollback: destination path must not be empty")
	}
	src := filepath.Join(r.backupDir, id)
	data, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("rollback: read backup %q: %w", id, err)
	}
	if err := os.WriteFile(dest, data, 0o600); err != nil {
		return fmt.Errorf("rollback: write destination %q: %w", dest, err)
	}
	return nil
}

// Latest returns the most recently modified backup entry, or an error if none exist.
func (r *Rollbacker) Latest() (Entry, error) {
	entries, err := r.List()
	if err != nil {
		return Entry{}, err
	}
	if len(entries) == 0 {
		return Entry{}, fmt.Errorf("rollback: no backups found in %q", r.backupDir)
	}
	latest := entries[0]
	for _, e := range entries[1:] {
		if e.Timestamp.After(latest.Timestamp) {
			latest = e
		}
	}
	return latest, nil
}
