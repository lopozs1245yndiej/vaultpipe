package ttl_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/vaultpipe/vaultpipe/internal/ttl"
)

func tmpFile(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "ttl.json")
}

func TestNew_InvalidTTL(t *testing.T) {
	_, err := ttl.New(0, "")
	if err == nil {
		t.Fatal("expected error for zero TTL")
	}
}

func TestNew_Valid(t *testing.T) {
	tr, err := ttl.New(time.Minute, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tr == nil {
		t.Fatal("expected non-nil tracker")
	}
}

func TestIsExpired_NewKey(t *testing.T) {
	tr, _ := ttl.New(time.Minute, "")
	if !tr.IsExpired("MISSING_KEY") {
		t.Error("untouched key should be expired")
	}
}

func TestIsExpired_AfterTouch(t *testing.T) {
	tr, _ := ttl.New(time.Minute, "")
	tr.Touch("DB_PASSWORD")
	if tr.IsExpired("DB_PASSWORD") {
		t.Error("recently touched key should not be expired")
	}
}

func TestIsExpired_AfterTTLElapsed(t *testing.T) {
	tr, _ := ttl.New(50*time.Millisecond, "")
	tr.Touch("API_KEY")
	time.Sleep(80 * time.Millisecond)
	if !tr.IsExpired("API_KEY") {
		t.Error("key should be expired after TTL elapsed")
	}
}

func TestSave_AndReload(t *testing.T) {
	path := tmpFile(t)

	tr, err := ttl.New(time.Hour, path)
	if err != nil {
		t.Fatalf("create tracker: %v", err)
	}
	tr.Touch("SECRET_ONE")
	tr.Touch("SECRET_TWO")

	if err := tr.Save(); err != nil {
		t.Fatalf("save: %v", err)
	}

	tr2, err := ttl.New(time.Hour, path)
	if err != nil {
		t.Fatalf("reload tracker: %v", err)
	}
	for _, key := range []string{"SECRET_ONE", "SECRET_TWO"} {
		if tr2.IsExpired(key) {
			t.Errorf("key %s should not be expired after reload", key)
		}
	}
}

func TestSave_NoPath_IsNoop(t *testing.T) {
	tr, _ := ttl.New(time.Minute, "")
	tr.Touch("X")
	if err := tr.Save(); err != nil {
		t.Fatalf("save with no path should not error: %v", err)
	}
}

func TestNew_NonExistentFile_ReturnsEmptyTracker(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nonexistent.json")
	tr, err := ttl.New(time.Minute, path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !tr.IsExpired("ANY") {
		t.Error("expected key to be expired in empty tracker")
	}
}

func TestNew_CorruptFile_ReturnsError(t *testing.T) {
	path := tmpFile(t)
	if err := os.WriteFile(path, []byte("not json{"), 0o600); err != nil {
		t.Fatal(err)
	}
	_, err := ttl.New(time.Minute, path)
	if err == nil {
		t.Fatal("expected error for corrupt file")
	}
}
