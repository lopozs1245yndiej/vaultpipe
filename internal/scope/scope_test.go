package scope_test

import (
	"testing"

	"github.com/your-org/vaultpipe/internal/scope"
)

func TestNew_EmptyNamespace_ReturnsError(t *testing.T) {
	_, err := scope.New("")
	if err == nil {
		t.Fatal("expected error for empty namespace")
	}
}

func TestNew_ValidNamespace_OK(t *testing.T) {
	s, err := scope.New("APP")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s == nil {
		t.Fatal("expected non-nil Scoper")
	}
}

func TestApply_StripsNamespacePrefix(t *testing.T) {
	s, _ := scope.New("APP_")
	secrets := map[string]string{
		"APP_DB_HOST": "localhost",
		"APP_DB_PORT": "5432",
		"OTHER_KEY":   "ignored",
	}
	out := s.Apply(secrets)
	if len(out) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(out))
	}
	if out["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %q", out["DB_HOST"])
	}
	if out["DB_PORT"] != "5432" {
		t.Errorf("expected DB_PORT=5432, got %q", out["DB_PORT"])
	}
}

func TestApply_CaseInsensitiveByDefault(t *testing.T) {
	s, _ := scope.New("app_")
	secrets := map[string]string{
		"APP_SECRET": "value",
	}
	out := s.Apply(secrets)
	if _, ok := out["SECRET"]; !ok {
		t.Error("expected case-insensitive match for APP_SECRET")
	}
}

func TestApply_CaseSensitive_NoMatch(t *testing.T) {
	s, _ := scope.New("app_", scope.WithCaseSensitive(true))
	secrets := map[string]string{
		"APP_SECRET": "value",
	}
	out := s.Apply(secrets)
	if len(out) != 0 {
		t.Errorf("expected 0 keys with case-sensitive match, got %d", len(out))
	}
}

func TestApply_WithReplacePrefix(t *testing.T) {
	s, _ := scope.New("APP_", scope.WithReplacePrefix("SVC_"))
	secrets := map[string]string{
		"APP_TOKEN": "abc123",
	}
	out := s.Apply(secrets)
	if out["SVC_TOKEN"] != "abc123" {
		t.Errorf("expected SVC_TOKEN=abc123, got %v", out)
	}
}

func TestApply_EmptySecrets_ReturnsEmpty(t *testing.T) {
	s, _ := scope.New("APP_")
	out := s.Apply(map[string]string{})
	if len(out) != 0 {
		t.Errorf("expected empty map, got %d entries", len(out))
	}
}

func TestApply_KeyEqualsNamespace_Excluded(t *testing.T) {
	// A key that is exactly the namespace (no suffix) should be excluded.
	s, _ := scope.New("APP_")
	secrets := map[string]string{
		"APP_": "bare",
	}
	out := s.Apply(secrets)
	if len(out) != 0 {
		t.Errorf("expected bare namespace key to be excluded, got %v", out)
	}
}
