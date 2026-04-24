package checkpoint_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourusername/vaultpipe/internal/checkpoint"
)

func tmpCheckpointFile(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "checkpoint.json")
}

func TestNew_EmptyFile_ReturnsTracker(t *testing.T) {
	p := tmpCheckpointFile(t)
	tr, err := checkpoint.New(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tr == nil {
		t.Fatal("expected non-nil tracker")
	}
}

func TestNew_NonExistentFile_ReturnsEmptyTracker(t *testing.T) {
	p := filepath.Join(t.TempDir(), "missing.json")
	tr, err := checkpoint.New(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tr.All()) != 0 {
		t.Fatalf("expected empty tracker, got %d entries", len(tr.All()))
	}
}

func TestSet_AndGet_RoundTrip(t *testing.T) {
	p := tmpCheckpointFile(t)
	tr, _ := checkpoint.New(p)

	now := time.Now().UTC().Truncate(time.Second)
	entry := checkpoint.Entry{
		Path:     "secret/myapp",
		SyncedAt: now,
		Checksum: "abc123",
	}
	if err := tr.Set(entry); err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	got, err := tr.Get("secret/myapp")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if got.Checksum != "abc123" {
		t.Errorf("expected checksum abc123, got %s", got.Checksum)
	}
	if !got.SyncedAt.Equal(now) {
		t.Errorf("expected time %v, got %v", now, got.SyncedAt)
	}
}

func TestSet_SetsTimestampIfZero(t *testing.T) {
	p := tmpCheckpointFile(t)
	tr, _ := checkpoint.New(p)

	_ = tr.Set(checkpoint.Entry{Path: "secret/x", Checksum: "z"})
	got, _ := tr.Get("secret/x")
	if got.SyncedAt.IsZero() {
		t.Error("expected SyncedAt to be set automatically")
	}
}

func TestGet_MissingKey_ReturnsErrNotFound(t *testing.T) {
	p := tmpCheckpointFile(t)
	tr, _ := checkpoint.New(p)

	_, err := tr.Get("secret/nonexistent")
	if err != checkpoint.ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestDelete_RemovesEntry(t *testing.T) {
	p := tmpCheckpointFile(t)
	tr, _ := checkpoint.New(p)

	_ = tr.Set(checkpoint.Entry{Path: "secret/del", Checksum: "x"})
	_ = tr.Delete("secret/del")

	_, err := tr.Get("secret/del")
	if err != checkpoint.ErrNotFound {
		t.Errorf("expected ErrNotFound after delete, got %v", err)
	}
}

func TestSet_PersistsToDisk(t *testing.T) {
	p := tmpCheckpointFile(t)
	tr, _ := checkpoint.New(p)
	_ = tr.Set(checkpoint.Entry{Path: "secret/persist", Checksum: "cksum"})

	tr2, err := checkpoint.New(p)
	if err != nil {
		t.Fatalf("reload failed: %v", err)
	}
	got, err := tr2.Get("secret/persist")
	if err != nil {
		t.Fatalf("Get after reload failed: %v", err)
	}
	if got.Checksum != "cksum" {
		t.Errorf("expected cksum, got %s", got.Checksum)
	}
}

func TestSet_EmptyPath_ReturnsError(t *testing.T) {
	p := tmpCheckpointFile(t)
	tr, _ := checkpoint.New(p)
	err := tr.Set(checkpoint.Entry{Path: "", Checksum: "x"})
	if err == nil {
		t.Error("expected error for empty path")
	}
}

func TestNew_InvalidJSON_ReturnsError(t *testing.T) {
	p := tmpCheckpointFile(t)
	_ = os.WriteFile(p, []byte("not-json{"), 0600)
	_, err := checkpoint.New(p)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}
