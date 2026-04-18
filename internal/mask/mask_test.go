package mask

import (
	"strings"
	"testing"
)

func TestMask_FullyRedacts(t *testing.T) {
	m := New(0)
	got := m.Mask("supersecret")
	if got != "********" {
		t.Fatalf("expected ********, got %q", got)
	}
}

func TestMask_RevealsPrefix(t *testing.T) {
	m := New(3)
	got := m.Mask("supersecret")
	if !strings.HasPrefix(got, "sup") {
		t.Fatalf("expected prefix 'sup', got %q", got)
	}
	if !strings.Contains(got, "****") {
		t.Fatalf("expected masked suffix, got %q", got)
	}
}

func TestMask_EmptyValue(t *testing.T) {
	m := New(0)
	if got := m.Mask(""); got != "" {
		t.Fatalf("expected empty string, got %q", got)
	}
}

func TestMask_RevealCharsExceedsLength(t *testing.T) {
	m := New(100)
	got := m.Mask("short")
	if got != "********" {
		t.Fatalf("expected full mask when revealChars >= len, got %q", got)
	}
}

func TestMask_NegativeRevealChars(t *testing.T) {
	m := New(-5)
	got := m.Mask("value")
	if got != "********" {
		t.Fatalf("expected ********, got %q", got)
	}
}

func TestMaskMap(t *testing.T) {
	m := New(0)
	secrets := map[string]string{
		"DB_PASSWORD": "hunter2",
		"API_KEY":     "abc123",
	}
	masked := m.MaskMap(secrets)
	for k, v := range masked {
		if v != "********" {
			t.Errorf("key %s: expected ********, got %q", k, v)
		}
	}
	// original must be unchanged
	if secrets["DB_PASSWORD"] != "hunter2" {
		t.Fatal("original map was mutated")
	}
}
