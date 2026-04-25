package dedupe

import (
	"testing"
)

func TestApply_KeepFirst_RetainsEarliestValue(t *testing.T) {
	d := New(StrategyKeepFirst)
	a := map[string]string{"KEY": "first", "ONLY_A": "a"}
	b := map[string]string{"KEY": "second", "ONLY_B": "b"}

	got, err := d.Apply(a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["KEY"] != "first" {
		t.Errorf("expected 'first', got %q", got["KEY"])
	}
	if got["ONLY_A"] != "a" || got["ONLY_B"] != "b" {
		t.Error("non-duplicate keys should be present")
	}
}

func TestApply_KeepLast_RetainsLatestValue(t *testing.T) {
	d := New(StrategyKeepLast)
	a := map[string]string{"KEY": "first"}
	b := map[string]string{"KEY": "second"}

	got, err := d.Apply(a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["KEY"] != "second" {
		t.Errorf("expected 'second', got %q", got["KEY"])
	}
}

func TestApply_ErrorStrategy_ReturnsDuplicateError(t *testing.T) {
	d := New(StrategyError)
	a := map[string]string{"KEY": "v1"}
	b := map[string]string{"KEY": "v2"}

	_, err := d.Apply(a, b)
	if err == nil {
		t.Fatal("expected error for duplicate key, got nil")
	}
}

func TestApply_NoConflict_ErrorStrategy_OK(t *testing.T) {
	d := New(StrategyError)
	a := map[string]string{"A": "1"}
	b := map[string]string{"B": "2"}

	got, err := d.Apply(a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Errorf("expected 2 keys, got %d", len(got))
	}
}

func TestApply_EmptySources_ReturnsEmptyMap(t *testing.T) {
	d := New(StrategyKeepFirst)
	got, err := d.Apply()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("expected empty map, got %d entries", len(got))
	}
}

func TestFingerprint_DeterministicAcrossOrders(t *testing.T) {
	a := map[string]string{"FOO": "bar", "BAZ": "qux"}
	b := map[string]string{"BAZ": "qux", "FOO": "bar"}

	if Fingerprint(a) != Fingerprint(b) {
		t.Error("fingerprints should be equal regardless of map iteration order")
	}
}

func TestFingerprint_DiffersOnValueChange(t *testing.T) {
	a := map[string]string{"KEY": "original"}
	b := map[string]string{"KEY": "changed"}

	if Fingerprint(a) == Fingerprint(b) {
		t.Error("fingerprints should differ when values change")
	}
}
