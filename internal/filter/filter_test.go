package filter_test

import (
	"testing"

	"github.com/yourusername/vaultpipe/internal/filter"
)

func TestApply_NoPatterns_ReturnsAll(t *testing.T) {
	f := filter.New(nil, nil)
	secrets := map[string]string{"DB_HOST": "localhost", "API_KEY": "secret"}
	out := f.Apply(secrets)
	if len(out) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(out))
	}
}

func TestApply_IncludePrefix(t *testing.T) {
	f := filter.New([]string{"DB_"}, nil)
	secrets := map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432", "API_KEY": "secret"}
	out := f.Apply(secrets)
	if len(out) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(out))
	}
	if _, ok := out["API_KEY"]; ok {
		t.Error("API_KEY should have been excluded")
	}
}

func TestApply_ExcludePrefix(t *testing.T) {
	f := filter.New(nil, []string{"INTERNAL_"})
	secrets := map[string]string{"DB_HOST": "localhost", "INTERNAL_TOKEN": "xyz"}
	out := f.Apply(secrets)
	if len(out) != 1 {
		t.Fatalf("expected 1 key, got %d", len(out))
	}
	if _, ok := out["INTERNAL_TOKEN"]; ok {
		t.Error("INTERNAL_TOKEN should have been excluded")
	}
}

func TestApply_GlobPattern(t *testing.T) {
	f := filter.New([]string{"DB_*"}, nil)
	secrets := map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432", "API_KEY": "secret"}
	out := f.Apply(secrets)
	if len(out) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(out))
	}
}

func TestApply_ExcludeTakesPrecedence(t *testing.T) {
	f := filter.New([]string{"DB_"}, []string{"DB_PASSWORD"})
	secrets := map[string]string{"DB_HOST": "localhost", "DB_PASSWORD": "s3cr3t"}
	out := f.Apply(secrets)
	if len(out) != 1 {
		t.Fatalf("expected 1 key, got %d", len(out))
	}
	if _, ok := out["DB_PASSWORD"]; ok {
		t.Error("DB_PASSWORD should have been excluded")
	}
}

func TestApply_EmptySecrets(t *testing.T) {
	f := filter.New([]string{"DB_"}, nil)
	out := f.Apply(map[string]string{})
	if len(out) != 0 {
		t.Fatalf("expected 0 keys, got %d", len(out))
	}
}
