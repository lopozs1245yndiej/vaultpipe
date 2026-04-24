package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/your-org/vaultpipe/internal/snapshot"
)

func tmpDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "snapshot-*")
	if err != nil {
		t.Fatalf("tmpDir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func TestSave_AndLoad_RoundTrip(t *testing.T) {
	m, err := snapshot.New(tmpDir(t))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	secrets := map[string]string{"DB_PASS": "secret", "API_KEY": "abc123"}
	snap, err := m.Save("vault/app", secrets)
	if err != nil {
		t.Fatalf("Save: %v", err)
	}
	loaded, err := m.Load(snap.ID)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if loaded.Source != "vault/app" {
		t.Errorf("source: got %q, want %q", loaded.Source, "vault/app")
	}
	if loaded.Secrets["DB_PASS"] != "secret" {
		t.Errorf("DB_PASS: got %q", loaded.Secrets["DB_PASS"])
	}
}

func TestList_ReturnsAllIDs(t *testing.T) {
	m, _ := snapshot.New(tmpDir(t))
	for i := 0; i < 3; i++ {
		if _, err := m.Save("vault/app", map[string]string{"K": "v"}); err != nil {
			t.Fatalf("Save: %v", err)
		}
	}
	ids, err := m.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(ids) != 3 {
		t.Errorf("expected 3 snapshots, got %d", len(ids))
	}
}

func TestDelete_RemovesSnapshot(t *testing.T) {
	m, _ := snapshot.New(tmpDir(t))
	snap, _ := m.Save("vault/app", map[string]string{"X": "y"})
	if err := m.Delete(snap.ID); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if _, err := m.Load(snap.ID); err == nil {
		t.Error("expected error loading deleted snapshot")
	}
}

func TestDelete_NoErrorIfMissing(t *testing.T) {
	m, _ := snapshot.New(tmpDir(t))
	if err := m.Delete("nonexistent"); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestNew_CreatesDirectory(t *testing.T) {
	base := tmpDir(t)
	dir := filepath.Join(base, "nested", "snapshots")
	if _, err := snapshot.New(dir); err != nil {
		t.Fatalf("New: %v", err)
	}
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Error("expected directory to be created")
	}
}
