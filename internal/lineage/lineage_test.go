package lineage_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourusername/vaultpipe/internal/lineage"
)

func tmpLineageFile(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "lineage.json")
}

func TestNew_CreatesEmptyTracker(t *testing.T) {
	path := tmpLineageFile(t)
	tracker, err := lineage.New(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := tracker.All(); len(got) != 0 {
		t.Errorf("expected 0 records, got %d", len(got))
	}
}

func TestTrack_AndGet(t *testing.T) {
	tracker, _ := lineage.New(tmpLineageFile(t))
	now := time.Now().UTC().Truncate(time.Second)
	tracker.Track("DB_PASSWORD", "secret/data/app", 3, now)

	r, ok := tracker.Get("DB_PASSWORD")
	if !ok {
		t.Fatal("expected record to exist")
	}
	if r.VaultPath != "secret/data/app" {
		t.Errorf("vault path: want %q, got %q", "secret/data/app", r.VaultPath)
	}
	if r.Version != 3 {
		t.Errorf("version: want 3, got %d", r.Version)
	}
	if !r.SyncedAt.Equal(now) {
		t.Errorf("synced_at: want %v, got %v", now, r.SyncedAt)
	}
}

func TestGet_MissingKey_ReturnsFalse(t *testing.T) {
	tracker, _ := lineage.New(tmpLineageFile(t))
	_, ok := tracker.Get("NONEXISTENT")
	if ok {
		t.Error("expected false for missing key")
	}
}

func TestSave_AndReload(t *testing.T) {
	path := tmpLineageFile(t)
	tracker, _ := lineage.New(path)
	now := time.Now().UTC().Truncate(time.Second)
	tracker.Track("API_KEY", "secret/data/svc", 1, now)

	if err := tracker.Save(); err != nil {
		t.Fatalf("save error: %v", err)
	}

	reloaded, err := lineage.New(path)
	if err != nil {
		t.Fatalf("reload error: %v", err)
	}
	r, ok := reloaded.Get("API_KEY")
	if !ok {
		t.Fatal("expected API_KEY after reload")
	}
	if r.VaultPath != "secret/data/svc" {
		t.Errorf("vault path after reload: got %q", r.VaultPath)
	}
}

func TestSave_WritesValidJSON(t *testing.T) {
	path := tmpLineageFile(t)
	tracker, _ := lineage.New(path)
	tracker.Track("TOKEN", "secret/data/tokens", 2, time.Now())
	_ = tracker.Save()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read file: %v", err)
	}
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Errorf("invalid JSON: %v", err)
	}
}

func TestAll_ReturnsAllRecords(t *testing.T) {
	tracker, _ := lineage.New(tmpLineageFile(t))
	tracker.Track("A", "path/a", 1, time.Now())
	tracker.Track("B", "path/b", 1, time.Now())
	tracker.Track("C", "path/c", 1, time.Now())
	if got := len(tracker.All()); got != 3 {
		t.Errorf("expected 3 records, got %d", got)
	}
}
