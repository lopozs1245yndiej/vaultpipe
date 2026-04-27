package schema_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/vaultpipe/internal/schema"
)

func writeTempSchema(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "schema.json")
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatalf("write temp schema: %v", err)
	}
	return p
}

func TestLoadFromFile_Valid(t *testing.T) {
	p := writeTempSchema(t, `[
		{"key":"API_KEY","type":"string","required":true},
		{"key":"PORT","type":"int"}
	]`)
	s, err := schema.LoadFromFile(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	errs := s.Validate(map[string]string{"API_KEY": "abc", "PORT": "8080"})
	if len(errs) != 0 {
		t.Fatalf("unexpected validation errors: %v", errs)
	}
}

func TestLoadFromFile_Missing(t *testing.T) {
	_, err := schema.LoadFromFile("/nonexistent/schema.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoadFromFile_InvalidJSON(t *testing.T) {
	p := writeTempSchema(t, `{not valid json}`)
	_, err := schema.LoadFromFile(p)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestLoadFromFile_DefaultsTypeToString(t *testing.T) {
	p := writeTempSchema(t, `[{"key":"FOO"}]`)
	s, err := schema.LoadFromFile(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	errs := s.Validate(map[string]string{"FOO": "bar"})
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
}

func TestMustLoad_Panics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for missing file")
		}
	}()
	schema.MustLoad("/no/such/file.json")
}
