package truncate_test

import (
	"testing"

	"github.com/your-org/vaultpipe/internal/truncate"
)

func TestNew_InvalidMaxLen(t *testing.T) {
	_, err := truncate.New(0, "...")
	if err == nil {
		t.Fatal("expected error for maxLen=0, got nil")
	}
	_, err = truncate.New(-5, "...")
	if err == nil {
		t.Fatal("expected error for negative maxLen, got nil")
	}
}

func TestNew_Valid(t *testing.T) {
	tr, err := truncate.New(10, "...")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tr == nil {
		t.Fatal("expected non-nil Truncator")
	}
}

func TestApply_ShortValues_Unchanged(t *testing.T) {
	tr, _ := truncate.New(20, "...")
	input := map[string]string{
		"KEY": "short",
		"OTHER": "also fine",
	}
	out := tr.Apply(input)
	for k, v := range input {
		if out[k] != v {
			t.Errorf("key %q: expected %q, got %q", k, v, out[k])
		}
	}
}

func TestApply_LongValue_Truncated(t *testing.T) {
	tr, _ := truncate.New(10, "...")
	out := tr.Apply(map[string]string{
		"SECRET": "this-is-a-very-long-secret-value",
	})
	got := out["SECRET"]
	// expect 7 runes kept + "..."
	want := "this-is-..."
	if got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}

func TestApply_NoSuffix(t *testing.T) {
	tr, _ := truncate.New(5, "")
	out := tr.Apply(map[string]string{"K": "abcdefgh"})
	if got := out["K"]; got != "abcde" {
		t.Errorf("expected %q, got %q", "abcde", got)
	}
}

func TestApply_SuffixFillsBudget(t *testing.T) {
	// maxLen equals suffix length — only suffix should remain
	tr, _ := truncate.New(3, "...")
	out := tr.Apply(map[string]string{"K": "toolongvalue"})
	if got := out["K"]; got != "..." {
		t.Errorf("expected %q, got %q", "...", got)
	}
}

func TestApply_UnicodeValues(t *testing.T) {
	tr, _ := truncate.New(4, "…")
	// each of these is a multi-byte rune
	out := tr.Apply(map[string]string{"K": "日本語テスト"})
	got := out["K"]
	// keep 3 runes + suffix rune "…"
	want := "日本語…"
	if got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}

func TestApply_EmptySecrets(t *testing.T) {
	tr, _ := truncate.New(10, "...")
	out := tr.Apply(map[string]string{})
	if len(out) != 0 {
		t.Errorf("expected empty map, got %v", out)
	}
}
