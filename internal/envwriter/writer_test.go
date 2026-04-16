package envwriter

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func tmpFile(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), ".env")
}

func TestWrite_AllSecrets(t *testing.T) {
	p := tmpFile(t)
	w := NewWriter(p, "")

	secrets := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
	}

	if err := w.Write(secrets); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(p)
	content := string(data)

	if !strings.Contains(content, "DB_HOST=localhost") {
		t.Errorf("expected DB_HOST in output, got:\n%s", content)
	}
	if !strings.Contains(content, "DB_PORT=5432") {
		t.Errorf("expected DB_PORT in output, got:\n%s", content)
	}
}

func TestWrite_NamespaceFilter(t *testing.T) {
	p := tmpFile(t)
	w := NewWriter(p, "APP")

	secrets := map[string]string{
		"APP_SECRET":  "abc123",
		"DB_PASSWORD": "hunter2",
	}

	if err := w.Write(secrets); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(p)
	content := string(data)

	if !strings.Contains(content, "APP_SECRET=abc123") {
		t.Errorf("expected APP_SECRET in output, got:\n%s", content)
	}
	if strings.Contains(content, "DB_PASSWORD") {
		t.Errorf("expected DB_PASSWORD to be filtered out, got:\n%s", content)
	}
}

func TestWrite_EscapesValueWithSpaces(t *testing.T) {
	p := tmpFile(t)
	w := NewWriter(p, "")

	if err := w.Write(map[string]string{"GREETING": "hello world"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(p)
	if !strings.Contains(string(data), `GREETING="hello world"`) {
		t.Errorf("expected quoted value, got: %s", string(data))
	}
}

func TestWrite_EmptySecrets(t *testing.T) {
	p := tmpFile(t)
	w := NewWriter(p, "")

	if err := w.Write(map[string]string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(p)
	if len(data) != 0 {
		t.Errorf("expected empty file, got: %s", string(data))
	}
}

func TestWrite_FilePermissions(t *testing.T) {
	p := tmpFile(t)
	w := NewWriter(p, "")

	_ = w.Write(map[string]string{"KEY": "val"})

	info, err := os.Stat(p)
	if err != nil {
		t.Fatalf("stat failed: %v", err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Errorf("expected 0600 permissions, got %v", info.Mode().Perm())
	}
}

func TestWrite_OverwritesExistingFile(t *testing.T) {
	p := tmpFile(t)
	w := NewWriter(p, "")

	// Write initial content
	if err := w.Write(map[string]string{"OLD_KEY": "old_val"}); err != nil {
		t.Fatalf("first write failed: %v", err)
	}

	// Overwrite with new content
	if err := w.Write(map[string]string{"NEW_KEY": "new_val"}); err != nil {
		t.Fatalf("second write failed: %v", err)
	}

	data, _ := os.ReadFile(p)
	content := string(data)

	if strings.Contains(content, "OLD_KEY") {
		t.Errorf("expected OLD_KEY to be overwritten, got:\n%s", content)
	}
	if !strings.Contains(content, "NEW_KEY=new_val") {
		t.Errorf("expected NEW_KEY in output, got:\n%s", content)
	}
}
