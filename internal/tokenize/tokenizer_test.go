package tokenize

import (
	"testing"
)

func TestNew_EmptyDelimiter(t *testing.T) {
	_, err := New("")
	if err != ErrEmptyDelimiter {
		t.Fatalf("expected ErrEmptyDelimiter, got %v", err)
	}
}

func TestNew_ValidDelimiter(t *testing.T) {
	tok, err := New(":")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tok == nil {
		t.Fatal("expected non-nil tokenizer")
	}
}

func TestSplit_BasicColon(t *testing.T) {
	tok, _ := New(":")
	parts := tok.Split("localhost:5432")
	if len(parts) != 2 || parts[0] != "localhost" || parts[1] != "5432" {
		t.Fatalf("unexpected parts: %v", parts)
	}
}

func TestToken_ValidIndex(t *testing.T) {
	tok, _ := New(":")
	v, err := tok.Token("user:pass:extra", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v != "pass" {
		t.Fatalf("expected 'pass', got %q", v)
	}
}

func TestToken_OutOfRange(t *testing.T) {
	tok, _ := New(":")
	_, err := tok.Token("host:port", 5)
	if err != ErrIndexOutOfRange {
		t.Fatalf("expected ErrIndexOutOfRange, got %v", err)
	}
}

func TestToken_NegativeIndex(t *testing.T) {
	tok, _ := New(":")
	_, err := tok.Token("host:port", -1)
	if err != ErrIndexOutOfRange {
		t.Fatalf("expected ErrIndexOutOfRange, got %v", err)
	}
}

func TestJoin_Reassembles(t *testing.T) {
	tok, _ := New("-")
	result := tok.Join([]string{"a", "b", "c"})
	if result != "a-b-c" {
		t.Fatalf("expected 'a-b-c', got %q", result)
	}
}

func TestReplace_ValidIndex(t *testing.T) {
	tok, _ := New(":")
	result, err := tok.Replace("host:oldport", 1, "9999")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "host:9999" {
		t.Fatalf("expected 'host:9999', got %q", result)
	}
}

func TestReplace_OutOfRange(t *testing.T) {
	tok, _ := New(":")
	_, err := tok.Replace("host:port", 10, "x")
	if err != ErrIndexOutOfRange {
		t.Fatalf("expected ErrIndexOutOfRange, got %v", err)
	}
}

func TestSplitMap_MultipleKeys(t *testing.T) {
	tok, _ := New("@")
	secrets := map[string]string{
		"DB_URL": "user@host",
		"CACHE":  "admin@redis",
	}
	out := tok.SplitMap(secrets)
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
	if out["DB_URL"][0] != "user" || out["DB_URL"][1] != "host" {
		t.Fatalf("unexpected tokens for DB_URL: %v", out["DB_URL"])
	}
}
