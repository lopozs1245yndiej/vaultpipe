package rollback_test

import (
	"testing"

	"github.com/yourusername/vaultpipe/internal/rollback"
)

func TestNewFromOptions_ValidDir(t *testing.T) {
	dir := t.TempDir()
	r, err := rollback.NewFromOptions(rollback.Options{BackupDir: dir})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r == nil {
		t.Fatal("expected non-nil rollbacker")
	}
}

func TestNewFromOptions_MissingDir_ReturnsError(t *testing.T) {
	_, err := rollback.NewFromOptions(rollback.Options{})
	if err == nil {
		t.Fatal("expected error when BackupDir is empty")
	}
}
