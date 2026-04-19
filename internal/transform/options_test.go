package transform_test

import (
	"testing"

	"github.com/yourusername/vaultpipe/internal/transform"
)

func TestNewFromOptions_AllEnabled(t *testing.T) {
	tr := transform.NewFromOptions(transform.Options{
		Uppercase:   true,
		Prefix:      "SVC_",
		TrimSpace:   true,
		ReplaceFrom: "-",
		ReplaceTo:   "_",
	})
	out := tr.Apply(map[string]string{"db-url": "  postgres  "})
	if out["SVC_DB_URL"] != "postgres" {
		t.Fatalf("unexpected result: %v", out)
	}
}

func TestNewFromOptions_OnlyUppercase(t *testing.T) {
	tr := transform.NewFromOptions(transform.Options{Uppercase: true})
	out := tr.Apply(map[string]string{"host": "localhost"})
	if out["HOST"] != "localhost" {
		t.Fatalf("unexpected result: %v", out)
	}
}

func TestNewFromOptions_NoOptions(t *testing.T) {
	tr := transform.NewFromOptions(transform.Options{})
	out := tr.Apply(map[string]string{"key": "value"})
	if out["key"] != "value" {
		t.Fatalf("expected passthrough, got %v", out)
	}
}
