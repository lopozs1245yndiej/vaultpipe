package prefix_test

import (
	"testing"

	"github.com/yourusername/vaultpipe/internal/prefix"
)

func TestNew_EmptyPrefix_ReturnsError(t *testing.T) {
	_, err := prefix.New("")
	if err == nil {
		t.Fatal("expected error for empty prefix, got nil")
	}
}

func TestNew_WhitespacePrefix_ReturnsError(t *testing.T) {
	_, err := prefix.New("   ")
	if err == nil {
		t.Fatal("expected error for whitespace prefix, got nil")
	}
}

func TestNew_ValidPrefix_OK(t *testing.T) {
	p, err := prefix.New("APP")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p == nil {
		t.Fatal("expected non-nil Prefixer")
	}
}

func TestAdd_PrependsPrefix(t *testing.T) {
	p, _ := prefix.New("APP")
	input := map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432"}
	out := p.Add(input)

	if out["APP_DB_HOST"] != "localhost" {
		t.Errorf("expected APP_DB_HOST=localhost, got %q", out["APP_DB_HOST"])
	}
	if out["APP_DB_PORT"] != "5432" {
		t.Errorf("expected APP_DB_PORT=5432, got %q", out["APP_DB_PORT"])
	}
	if len(out) != 2 {
		t.Errorf("expected 2 keys, got %d", len(out))
	}
}

func TestAdd_CustomSeparator(t *testing.T) {
	p, _ := prefix.New("APP", prefix.WithSeparator("."))
	out := p.Add(map[string]string{"KEY": "val"})
	if out["APP.KEY"] != "val" {
		t.Errorf("expected APP.KEY=val, got %v", out)
	}
}

func TestStrip_RemovesPrefix(t *testing.T) {
	p, _ := prefix.New("APP")
	input := map[string]string{"APP_DB_HOST": "localhost", "OTHER": "value"}
	out := p.Strip(input)

	if out["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %q", out["DB_HOST"])
	}
	if out["OTHER"] != "value" {
		t.Errorf("expected OTHER=value to pass through, got %q", out["OTHER"])
	}
}

func TestStrip_NonMatchingKeys_PassThrough(t *testing.T) {
	p, _ := prefix.New("APP")
	input := map[string]string{"UNRELATED": "x"}
	out := p.Strip(input)
	if out["UNRELATED"] != "x" {
		t.Errorf("expected UNRELATED to pass through")
	}
}

func TestReplace_SwapsPrefix(t *testing.T) {
	p, _ := prefix.New("OLD")
	input := map[string]string{"OLD_KEY": "val", "KEEP": "yes"}
	out := p.Replace(input, "NEW")

	if out["NEW_KEY"] != "val" {
		t.Errorf("expected NEW_KEY=val, got %v", out)
	}
	if out["KEEP"] != "yes" {
		t.Errorf("expected KEEP=yes to pass through, got %v", out)
	}
	if _, exists := out["OLD_KEY"]; exists {
		t.Error("OLD_KEY should have been replaced")
	}
}
