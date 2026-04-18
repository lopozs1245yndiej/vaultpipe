package export

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func tmpPath(t *testing.T, name string) string {
	t.Helper()
	return filepath.Join(t.TempDir(), name)
}

func TestNew_ValidFormats(t *testing.T) {
	for _, f := range []Format{FormatEnv, FormatJSON, FormatYAML} {
		_, err := New(f)
		if err != nil {
			t.Errorf("expected no error for format %q, got %v", f, err)
		}
	}
}

func TestNew_InvalidFormat(t *testing.T) {
	_, err := New("toml")
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
}

func TestExport_EnvFormat(t *testing.T) {
	ex, _ := New(FormatEnv)
	p := tmpPath(t, "out.env")
	secrets := map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432"}
	if err := ex.Export(secrets, p); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	b, _ := os.ReadFile(p)
	content := string(b)
	if !strings.Contains(content, "DB_HOST=localhost") || !strings.Contains(content, "DB_PORT=5432") {
		t.Errorf("unexpected env content: %s", content)
	}
}

func TestExport_JSONFormat(t *testing.T) {
	ex, _ := New(FormatJSON)
	p := tmpPath(t, "out.json")
	secrets := map[string]string{"KEY": "value"}
	if err := ex.Export(secrets, p); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	b, _ := os.ReadFile(p)
	var out map[string]string
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if out["KEY"] != "value" {
		t.Errorf("expected KEY=value, got %v", out)
	}
}

func TestExport_YAMLFormat(t *testing.T) {
	ex, _ := New(FormatYAML)
	p := tmpPath(t, "out.yaml")
	secrets := map[string]string{"APP_ENV": "production"}
	if err := ex.Export(secrets, p); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	b, _ := os.ReadFile(p)
	if !strings.Contains(string(b), "APP_ENV:") {
		t.Errorf("unexpected yaml content: %s", string(b))
	}
}

func TestExport_FilePermissions(t *testing.T) {
	ex, _ := New(FormatEnv)
	p := tmpPath(t, "secure.env")
	if err := ex.Export(map[string]string{"X": "1"}, p); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	info, _ := os.Stat(p)
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected 0600 permissions, got %v", info.Mode().Perm())
	}
}
