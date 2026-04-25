package sanitize_test

import (
	"testing"

	"github.com/yourusername/vaultpipe/internal/sanitize"
)

func TestApply_NoOptions_ReturnsUnchanged(t *testing.T) {
	s := sanitize.New()
	input := map[string]string{"my-key": "  hello  ", "OTHER": "world"}
	out := s.Apply(input)
	if out["my-key"] != "  hello  " {
		t.Errorf("expected value unchanged, got %q", out["my-key"])
	}
}

func TestApply_NormalizeKeys(t *testing.T) {
	s := sanitize.New(sanitize.WithNormalizeKeys())
	out := s.Apply(map[string]string{
		"db-host":   "localhost",
		"api.token": "secret",
		"PORT":      "5432",
	})
	if _, ok := out["DB_HOST"]; !ok {
		t.Error("expected DB_HOST key")
	}
	if _, ok := out["API_TOKEN"]; !ok {
		t.Error("expected API_TOKEN key")
	}
	if _, ok := out["PORT"]; !ok {
		t.Error("expected PORT key")
	}
}

func TestApply_TrimValues(t *testing.T) {
	s := sanitize.New(sanitize.WithTrimValues())
	out := s.Apply(map[string]string{"KEY": "  value  "})
	if out["KEY"] != "value" {
		t.Errorf("expected trimmed value, got %q", out["KEY"])
	}
}

func TestApply_StripControlChars(t *testing.T) {
	s := sanitize.New(sanitize.WithStripControlChars())
	out := s.Apply(map[string]string{"KEY": "val\x00ue\x01"})
	if out["KEY"] != "value" {
		t.Errorf("expected control chars stripped, got %q", out["KEY"])
	}
}

func TestApply_StripControlChars_PreservesTab(t *testing.T) {
	s := sanitize.New(sanitize.WithStripControlChars())
	out := s.Apply(map[string]string{"KEY": "col1\tcol2"})
	if out["KEY"] != "col1\tcol2" {
		t.Errorf("expected tab preserved, got %q", out["KEY"])
	}
}

func TestApply_ChainedOptions(t *testing.T) {
	s := sanitize.New(
		sanitize.WithNormalizeKeys(),
		sanitize.WithTrimValues(),
		sanitize.WithStripControlChars(),
	)
	out := s.Apply(map[string]string{"db-host": "  local\x00host  "})
	if out["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %q", out["DB_HOST"])
	}
}

func TestApply_DoesNotMutateInput(t *testing.T) {
	s := sanitize.New(sanitize.WithNormalizeKeys(), sanitize.WithTrimValues())
	input := map[string]string{"my-key": "  val  "}
	_ = s.Apply(input)
	if _, ok := input["my-key"]; !ok {
		t.Error("original input was mutated")
	}
	if input["my-key"] != "  val  " {
		t.Error("original value was mutated")
	}
}
