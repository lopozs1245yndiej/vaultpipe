package rotate

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Rotator handles backup and rotation of .env files before overwriting.
type Rotator struct {
	backupDir string
	maxBackups int
}

// New creates a new Rotator with the given backup directory and max backup count.
func New(backupDir string, maxBackups int) *Rotator {
	if maxBackups <= 0 {
		maxBackups = 5
	}
	return &Rotator{
		backupDir:  backupDir,
		maxBackups: maxBackups,
	}
}

// Rotate creates a timestamped backup of the given file, then prunes old backups.
func (r *Rotator) Rotate(envFilePath string) error {
	if _, err := os.Stat(envFilePath); os.IsNotExist(err) {
		return nil
	}

	if err := os.MkdirAll(r.backupDir, 0700); err != nil {
		return fmt.Errorf("rotate: create backup dir: %w", err)
	}

	data, err := os.ReadFile(envFilePath)
	if err != nil {
		return fmt.Errorf("rotate: read source file: %w", err)
	}

	baseName := filepath.Base(envFilePath)
	timestamp := time.Now().UTC().Format("20060102T150405Z")
	backupName := fmt.Sprintf("%s.%s.bak", baseName, timestamp)
	backupPath := filepath.Join(r.backupDir, backupName)

	if err := os.WriteFile(backupPath, data, 0600); err != nil {
		return fmt.Errorf("rotate: write backup: %w", err)
	}

	return r.pruneBackups(baseName)
}

// pruneBackups removes the oldest backups when the count exceeds maxBackups.
func (r *Rotator) pruneBackups(baseName string) error {
	pattern := filepath.Join(r.backupDir, baseName+".*.bak")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return fmt.Errorf("rotate: glob backups: %w", err)
	}

	if len(matches) <= r.maxBackups {
		return nil
	}

	// matches are lexicographically sorted; oldest timestamps come first
	toRemove := matches[:len(matches)-r.maxBackups]
	for _, f := range toRemove {
		if err := os.Remove(f); err != nil {
			return fmt.Errorf("rotate: remove old backup %s: %w", f, err)
		}
	}
	return nil
}
