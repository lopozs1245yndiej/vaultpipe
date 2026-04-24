package rollback_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourusername/vaultpipe/internal/rollback"
)

func writeBackup(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatalf("writeBackup: %v", err)
	}
	return p
}

func TestNew_EmptyDir_ReturnsError(t *testing.T) {
	_, err := rollback.New("")
	if err == nil {
		t.Fatal("expected error for empty backup dir")
	}
}

func TestNew_CreatesDirectory(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "backups")
	r, err := rollback.New(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r == nil {
		t.Fatal("expected non-nil rollbacker")
	}
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Fatal("expected backup directory to be created")
	}
}

func TestList_ReturnsBackupEntries(t *testing.T) {
	dir := t.TempDir()
	writeBackup(t, dir, "a.env.bak", "KEY=a")
	writeBackup(t, dir, "b.env.bak", "KEY=b")

	r, _ := rollback.New(dir)
	entries, err := r.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
}

func TestRestore_WritesContentToDest(t *testing.T) {
	dir := t.TempDir()
	writeBackup(t, dir, "snap.env.bak", "SECRET=hello")

	r, _ := rollback.New(dir)
	dest := filepath.Join(t.TempDir(), ".env")
	if err := r.Restore("snap.env.bak", dest); err != nil {
		t.Fatalf("Restore: %v", err)
	}
	data, _ := os.ReadFile(dest)
	if string(data) != "SECRET=hello" {
		t.Fatalf("unexpected content: %q", string(data))
	}
}

func TestRestore_MissingBackup_ReturnsError(t *testing.T) {
	dir := t.TempDir()
	r, _ := rollback.New(dir)
	err := r.Restore("missing.env.bak", filepath.Join(t.TempDir(), ".env"))
	if err == nil {
		t.Fatal("expected error for missing backup")
	}
}

func TestLatest_ReturnsNewest(t *testing.T) {
	dir := t.TempDir()
	old := writeBackup(t, dir, "old.env.bak", "KEY=old")
	new_ := writeBackup(t, dir, "new.env.bak", "KEY=new")

	// ensure mtime difference
	_ = old
	os.Chtimes(new_, time.Now().Add(time.Hour), time.Now().Add(time.Hour))

	r, _ := rollback.New(dir)
	e, err := r.Latest()
	if err != nil {
		t.Fatalf("Latest: %v", err)
	}
	if e.ID != "new.env.bak" {
		t.Fatalf("expected new.env.bak, got %q", e.ID)
	}
}

func TestLatest_NoBackups_ReturnsError(t *testing.T) {
	dir := t.TempDir()
	r, _ := rollback.New(dir)
	_, err := r.Latest()
	if err == nil {
		t.Fatal("expected error when no backups exist")
	}
}
