package validate

import (
	"strings"
	"testing"
)

func TestValidate_AllValid(t *testing.T) {
	v := New(nil)
	secrets := map[string]string{
		"DATABASE_URL": "postgres://localhost",
		"API_KEY":      "abc123",
	}
	r := v.Validate(secrets)
	if !r.Valid {
		t.Errorf("expected valid, got errors: %+v", r)
	}
}

func TestValidate_InvalidKeyName(t *testing.T) {
	v := New(nil)
	secrets := map[string]string{
		"bad-key":  "value",
		"GOOD_KEY": "value",
	}
	r := v.Validate(secrets)
	if r.Valid {
		t.Error("expected invalid result")
	}
	if len(r.InvalidKeys) != 1 || r.InvalidKeys[0] != "bad-key" {
		t.Errorf("unexpected invalid keys: %v", r.InvalidKeys)
	}
}

func TestValidate_MissingRequiredKey(t *testing.T) {
	v := New([]string{"REQUIRED_KEY", "ANOTHER_KEY"})
	secrets := map[string]string{
		"REQUIRED_KEY": "present",
	}
	r := v.Validate(secrets)
	if r.Valid {
		t.Error("expected invalid due to missing key")
	}
	if len(r.MissingKeys) != 1 || r.MissingKeys[0] != "ANOTHER_KEY" {
		t.Errorf("unexpected missing keys: %v", r.MissingKeys)
	}
}

func TestValidate_EmptySecrets(t *testing.T) {
	v := New([]string{"MUST_EXIST"})
	r := v.Validate(map[string]string{})
	if r.Valid {
		t.Error("expected invalid")
	}
	if len(r.MissingKeys) != 1 {
		t.Errorf("expected 1 missing key, got %v", r.MissingKeys)
	}
}

func TestFormatErrors_Valid(t *testing.T) {
	r := Result{Valid: true}
	if out := FormatErrors(r); out != "" {
		t.Errorf("expected empty string, got %q", out)
	}
}

func TestFormatErrors_ShowsBothErrors(t *testing.T) {
	r := Result{
		Valid:       false,
		InvalidKeys: []string{"bad-key"},
		MissingKeys: []string{"SECRET_TOKEN"},
	}
	out := FormatErrors(r)
	if !strings.Contains(out, "bad-key") {
		t.Errorf("expected bad-key in output: %s", out)
	}
	if !strings.Contains(out, "SECRET_TOKEN") {
		t.Errorf("expected SECRET_TOKEN in output: %s", out)
	}
}
