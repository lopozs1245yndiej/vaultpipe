package lock

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// ErrLocked is returned when a lock file already exists and is still valid.
var ErrLocked = errors.New("lock: resource is already locked")

// Locker manages a PID-based lock file to prevent concurrent vaultpipe runs.
type Locker struct {
	path string
}

// New creates a Locker for the given lock file path.
func New(path string) *Locker {
	return &Locker{path: path}
}

// Acquire creates the lock file. Returns ErrLocked if a live lock exists.
func (l *Locker) Acquire() error {
	if err := os.MkdirAll(filepath.Dir(l.path), 0o755); err != nil {
		return fmt.Errorf("lock: create dir: %w", err)
	}

	if existing, err := l.readPID(); err == nil {
		if isProcessAlive(existing) {
			return ErrLocked
		}
		// Stale lock — remove it.
		_ = os.Remove(l.path)
	}

	content := fmt.Sprintf("%d\n%d\n", os.Getpid(), time.Now().Unix())
	return os.WriteFile(l.path, []byte(content), 0o644)
}

// Release removes the lock file.
func (l *Locker) Release() error {
	if err := os.Remove(l.path); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("lock: release: %w", err)
	}
	return nil
}

// IsLocked reports whether a live lock file currently exists.
func (l *Locker) IsLocked() bool {
	pid, err := l.readPID()
	if err != nil {
		return false
	}
	return isProcessAlive(pid)
}

func (l *Locker) readPID() (int, error) {
	data, err := os.ReadFile(l.path)
	if err != nil {
		return 0, err
	}
	lines := strings.SplitN(strings.TrimSpace(string(data)), "\n", 2)
	if len(lines) == 0 {
		return 0, errors.New("lock: empty file")
	}
	return strconv.Atoi(strings.TrimSpace(lines[0]))
}

func isProcessAlive(pid int) bool {
	proc, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	// On Unix, Signal(0) checks existence without sending a real signal.
	return proc.Signal(os.Signal(nil)) == nil
}
