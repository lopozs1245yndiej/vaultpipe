package lock_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/vaultpipe/internal/lock"
)

func tmpLockFile(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, "vaultpipe.lock")
}

func TestAcquire_CreatesLockFile(t *testing.T) {
	path := tmpLockFile(t)
	l := lock.New(path)

	if err := l.Acquire(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	defer l.Release() //nolint:errcheck

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatal("expected lock file to exist")
	}
}

func TestAcquire_ReturnErrLockedWhenAlive(t *testing.T) {
	path := tmpLockFile(t)
	l := lock.New(path)

	if err := l.Acquire(); err != nil {
		t.Fatalf("first acquire failed: %v", err)
	}
	defer l.Release() //nolint:errcheck

	err := l.Acquire()
	if err == nil {
		t.Fatal("expected ErrLocked, got nil")
	}
	if err != lock.ErrLocked {
		t.Fatalf("expected ErrLocked, got %v", err)
	}
}

func TestRelease_RemovesLockFile(t *testing.T) {
	path := tmpLockFile(t)
	l := lock.New(path)

	_ = l.Acquire()
	if err := l.Release(); err != nil {
		t.Fatalf("release failed: %v", err)
	}

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatal("expected lock file to be removed")
	}
}

func TestRelease_NoErrorIfFileMissing(t *testing.T) {
	path := tmpLockFile(t)
	l := lock.New(path)

	if err := l.Release(); err != nil {
		t.Fatalf("expected no error releasing non-existent lock, got %v", err)
	}
}

func TestIsLocked_FalseAfterRelease(t *testing.T) {
	path := tmpLockFile(t)
	l := lock.New(path)

	_ = l.Acquire()
	_ = l.Release()

	if l.IsLocked() {
		t.Fatal("expected IsLocked to return false after release")
	}
}

func TestAcquire_RemovesStaleLock(t *testing.T) {
	path := tmpLockFile(t)
	// Write a stale lock with PID 0 (never a real process).
	_ = os.WriteFile(path, []byte("0\n0\n"), 0o644)

	l := lock.New(path)
	if err := l.Acquire(); err != nil {
		t.Fatalf("expected stale lock to be removed, got %v", err)
	}
	defer l.Release() //nolint:errcheck
}
