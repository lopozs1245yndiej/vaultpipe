package resolve

import (
	"strings"
	"testing"
)

func TestNew_InvalidMaxDepth(t *testing.T) {
	_, err := New(0)
	if err == nil {
		t.Fatal("expected error for maxDepth=0")
	}
	if !strings.Contains(err.Error(), "maxDepth") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestNew_Valid(t *testing.T) {
	r, err := New(5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r == nil {
		t.Fatal("expected non-nil resolver")
	}
}

func TestResolve_NoReferences(t *testing.T) {
	r, _ := New(5)
	secrets := map[string]string{"FOO": "bar", "BAZ": "qux"}
	out, err := r.Resolve(secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["FOO"] != "bar" || out["BAZ"] != "qux" {
		t.Errorf("unexpected output: %v", out)
	}
}

func TestResolve_SimpleReference(t *testing.T) {
	r, _ := New(5)
	secrets := map[string]string{
		"HOST": "localhost",
		"DSN":  "postgres://${HOST}/db",
	}
	out, err := r.Resolve(secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DSN"] != "postgres://localhost/db" {
		t.Errorf("got %q, want %q", out["DSN"], "postgres://localhost/db")
	}
}

func TestResolve_ChainedReference(t *testing.T) {
	r, _ := New(5)
	secrets := map[string]string{
		"A": "hello",
		"B": "${A}_world",
		"C": "${B}!",
	}
	out, err := r.Resolve(secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["C"] != "hello_world!" {
		t.Errorf("got %q, want %q", out["C"], "hello_world!")
	}
}

func TestResolve_UndefinedReference(t *testing.T) {
	r, _ := New(5)
	secrets := map[string]string{
		"URL": "https://${MISSING_HOST}/path",
	}
	_, err := r.Resolve(secrets)
	if err == nil {
		t.Fatal("expected error for undefined reference")
	}
	if !strings.Contains(err.Error(), "MISSING_HOST") {
		t.Errorf("error should mention missing key, got: %v", err)
	}
}

func TestResolve_CircularReference_ExceedsDepth(t *testing.T) {
	r, _ := New(3)
	secrets := map[string]string{
		"A": "${B}",
		"B": "${A}",
	}
	_, err := r.Resolve(secrets)
	if err == nil {
		t.Fatal("expected error for circular reference")
	}
	if !strings.Contains(err.Error(), "depth") {
		t.Errorf("error should mention depth, got: %v", err)
	}
}

func TestResolve_MultipleReferencesInValue(t *testing.T) {
	r, _ := New(5)
	secrets := map[string]string{
		"USER": "admin",
		"PASS": "secret",
		"DSN":  "${USER}:${PASS}@host",
	}
	out, err := r.Resolve(secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DSN"] != "admin:secret@host" {
		t.Errorf("got %q, want %q", out["DSN"], "admin:secret@host")
	}
}
