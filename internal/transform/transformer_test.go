package transform_test

import (
	"testing"

	"github.com/yourusername/vaultpipe/internal/transform"
)

func TestApply_UppercaseKeys(t *testing.T) {
	tr := transform.New(transform.UppercaseKeys())
	out := tr.Apply(map[string]string{"foo": "bar", "baz": "qux"})
	if out["FOO"] != "bar" || out["BAZ"] != "qux" {
		t.Fatalf("expected uppercased keys, got %v", out)
	}
}

func TestApply_PrefixKeys(t *testing.T) {
	tr := transform.New(transform.PrefixKeys("APP_"))
	out := tr.Apply(map[string]string{"DB": "postgres"})
	if out["APP_DB"] != "postgres" {
		t.Fatalf("expected prefixed key, got %v", out)
	}
}

func TestApply_TrimValueSpace(t *testing.T) {
	tr := transform.New(transform.TrimValueSpace())
	out := tr.Apply(map[string]string{"KEY": "  hello  "})
	if out["KEY"] != "hello" {
		t.Fatalf("expected trimmed value, got %q", out["KEY"])
	}
}

func TestApply_ReplaceKeyChars(t *testing.T) {
	tr := transform.New(transform.ReplaceKeyChars("-", "_"))
	out := tr.Apply(map[string]string{"my-key": "val"})
	if out["my_key"] != "val" {
		t.Fatalf("expected replaced key, got %v", out)
	}
}

func TestApply_ChainedTransforms(t *testing.T) {
	tr := transform.New(
		transform.ReplaceKeyChars("-", "_"),
		transform.UppercaseKeys(),
		transform.PrefixKeys("APP_"),
		transform.TrimValueSpace(),
	)
	out := tr.Apply(map[string]string{"db-host": "  localhost  "})
	if out["APP_DB_HOST"] != "localhost" {
		t.Fatalf("expected chained transform result, got %v", out)
	}
}

func TestApply_EmptySecrets(t *testing.T) {
	tr := transform.New(transform.UppercaseKeys())
	out := tr.Apply(map[string]string{})
	if len(out) != 0 {
		t.Fatalf("expected empty map, got %v", out)
	}
}

func TestApply_NoTransforms(t *testing.T) {
	tr := transform.New()
	in := map[string]string{"key": "value"}
	out := tr.Apply(in)
	if out["key"] != "value" {
		t.Fatalf("expected passthrough, got %v", out)
	}
}
