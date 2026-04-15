package rotate

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTemp(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatalf("writeTemp: %v", err)
	}
	return p
}

func TestRotate_CreatesBackup(t *testing.T) {
	tmpDir := t.TempDir()
	envFile := writeTemp(t, tmpDir, ".env", "KEY=value\n")
	backupDir := filepath.Join(tmpDir, "backups")

	r := New(backupDir, 5)
	if err := r.Rotate(envFile); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	matches, _ := filepath.Glob(filepath.Join(backupDir, ".env.*.bak"))
	if len(matches) != 1 {
		t.Fatalf("expected 1 backup, got %d", len(matches))
	}

	data, _ := os.ReadFile(matches[0])
	if string(data) != "KEY=value\n" {
		t.Errorf("backup content mismatch: %q", string(data))
	}
}

func TestRotate_NoErrorIfFileMissing(t *testing.T) {
	tmpDir := t.TempDir()
	r := New(filepath.Join(tmpDir, "backups"), 5)
	if err := r.Rotate(filepath.Join(tmpDir, "nonexistent.env")); err != nil {
		t.Fatalf("expected no error for missing file, got: %v", err)
	}
}

func TestRotate_PrunesOldBackups(t *testing.T) {
	tmpDir := t.TempDir()
	envFile := writeTemp(t, tmpDir, ".env", "KEY=value\n")
	backupDir := filepath.Join(tmpDir, "backups")
	if err := os.MkdirAll(backupDir, 0700); err != nil {
		t.Fatal(err)
	}

	// Pre-create 5 old backups
	for i := 0; i < 5; i++ {
		name := filepath.Join(backupDir, ".env.2024010"+string(rune('1'+i))+"T000000Z.bak")
		_ = os.WriteFile(name, []byte("old"), 0600)
	}

	r := New(backupDir, 5)
	if err := r.Rotate(envFile); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	matches, _ := filepath.Glob(filepath.Join(backupDir, ".env.*.bak"))
	if len(matches) != 5 {
		t.Errorf("expected 5 backups after prune, got %d", len(matches))
	}
}

func TestNew_DefaultsMaxBackups(t *testing.T) {
	r := New("/tmp", 0)
	if r.maxBackups != 5 {
		t.Errorf("expected default maxBackups=5, got %d", r.maxBackups)
	}
}
