package policy

import (
	"testing"
)

func TestNew_ValidRules(t *testing.T) {
	_, err := New([]Rule{
		{Key: "SECRET_*", Allowed: false},
		{Key: "*", Allowed: true},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNew_EmptyKeyPattern(t *testing.T) {
	_, err := New([]Rule{{Key: "", Allowed: true}})
	if err == nil {
		t.Fatal("expected error for empty key pattern")
	}
}

func TestIsAllowed_ExplicitDeny(t *testing.T) {
	p, _ := New([]Rule{
		{Key: "SECRET_KEY", Allowed: false},
	})
	if p.IsAllowed("SECRET_KEY") {
		t.Error("expected SECRET_KEY to be denied")
	}
}

func TestIsAllowed_DefaultAllow(t *testing.T) {
	p, _ := New([]Rule{})
	if !p.IsAllowed("ANY_KEY") {
		t.Error("expected ANY_KEY to be allowed by default")
	}
}

func TestIsAllowed_GlobDeny(t *testing.T) {
	p, _ := New([]Rule{
		{Key: "INTERNAL_*", Allowed: false},
	})
	if p.IsAllowed("INTERNAL_TOKEN") {
		t.Error("expected INTERNAL_TOKEN to be denied by glob")
	}
	if !p.IsAllowed("PUBLIC_KEY") {
		t.Error("expected PUBLIC_KEY to be allowed")
	}
}

func TestApply_FiltersSecrets(t *testing.T) {
	p, _ := New([]Rule{
		{Key: "DENY_*", Allowed: false},
	})
	secrets := map[string]string{
		"DENY_ME":  "secret",
		"KEEP_ME":  "value",
		"DENY_TOO": "other",
	}
	out := p.Apply(secrets)
	if _, ok := out["DENY_ME"]; ok {
		t.Error("DENY_ME should have been filtered")
	}
	if _, ok := out["KEEP_ME"]; !ok {
		t.Error("KEEP_ME should be present")
	}
	if len(out) != 1 {
		t.Errorf("expected 1 key, got %d", len(out))
	}
}

func TestViolations_ReturnsDeniedKeys(t *testing.T) {
	p, _ := New([]Rule{
		{Key: "PRIVATE_*", Allowed: false},
	})
	secrets := map[string]string{
		"PRIVATE_KEY":    "x",
		"PRIVATE_SECRET": "y",
		"PUBLIC_KEY":     "z",
	}
	v := p.Violations(secrets)
	if len(v) != 2 {
		t.Errorf("expected 2 violations, got %d", len(v))
	}
}

func TestMatchGlob_ExactMatch(t *testing.T) {
	if !matchGlob("EXACT", "EXACT") {
		t.Error("exact match should succeed")
	}
	if matchGlob("EXACT", "EXACT_EXTRA") {
		t.Error("should not match with extra suffix")
	}
}
