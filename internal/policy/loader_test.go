package policy

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func writePolicyFile(t *testing.T, rules []FileRule) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "policy.json")
	data, err := json.Marshal(policyFile{Rules: rules})
	if err != nil {
		t.Fatalf("marshal policy: %v", err)
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatalf("write policy file: %v", err)
	}
	return path
}

func TestLoadFromFile_Valid(t *testing.T) {
	path := writePolicyFile(t, []FileRule{
		{Key: "SECRET_*", Allowed: false},
		{Key: "*", Allowed: true},
	})
	p, err := LoadFromFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.IsAllowed("SECRET_TOKEN") {
		t.Error("SECRET_TOKEN should be denied")
	}
	if !p.IsAllowed("PUBLIC_KEY") {
		t.Error("PUBLIC_KEY should be allowed")
	}
}

func TestLoadFromFile_Missing(t *testing.T) {
	_, err := LoadFromFile("/nonexistent/policy.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoadFromFile_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	_ = os.WriteFile(path, []byte("not-json"), 0o600)
	_, err := LoadFromFile(path)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestLoadFromFile_EmptyRules(t *testing.T) {
	path := writePolicyFile(t, []FileRule{})
	p, err := LoadFromFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !p.IsAllowed("ANY_KEY") {
		t.Error("expected default-allow with no rules")
	}
}

func TestMustLoad_Panics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for missing file")
		}
	}()
	MustLoad("/no/such/file.json")
}
