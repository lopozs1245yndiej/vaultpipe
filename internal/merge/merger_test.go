package merge_test

import (
	"testing"

	"github.com/yourusername/vaultpipe/internal/merge"
)

func TestMerge_LastWins(t *testing.T) {
	m := merge.New(merge.StrategyLastWins)
	a := map[string]string{"KEY": "first", "A": "1"}
	b := map[string]string{"KEY": "second", "B": "2"}

	result, err := m.Merge(a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["KEY"] != "second" {
		t.Errorf("expected 'second', got %q", result["KEY"])
	}
	if result["A"] != "1" || result["B"] != "2" {
		t.Errorf("missing non-conflicting keys")
	}
}

func TestMerge_FirstWins(t *testing.T) {
	m := merge.New(merge.StrategyFirstWins)
	a := map[string]string{"KEY": "first"}
	b := map[string]string{"KEY": "second"}

	result, err := m.Merge(a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["KEY"] != "first" {
		t.Errorf("expected 'first', got %q", result["KEY"])
	}
}

func TestMerge_ErrorOnConflict(t *testing.T) {
	m := merge.New(merge.StrategyError)
	a := map[string]string{"KEY": "first"}
	b := map[string]string{"KEY": "second"}

	_, err := m.Merge(a, b)
	if err == nil {
		t.Fatal("expected conflict error, got nil")
	}
}

func TestMerge_NoConflict_ErrorStrategy(t *testing.T) {
	m := merge.New(merge.StrategyError)
	a := map[string]string{"A": "1"}
	b := map[string]string{"B": "2"}

	result, err := m.Merge(a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 keys, got %d", len(result))
	}
}

func TestMerge_EmptySources(t *testing.T) {
	m := merge.New(merge.StrategyLastWins)
	result, err := m.Merge()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected empty map, got %d keys", len(result))
	}
}

func TestMergeAll_Convenience(t *testing.T) {
	a := map[string]string{"X": "old"}
	b := map[string]string{"X": "new", "Y": "1"}

	result := merge.MergeAll(a, b)
	if result["X"] != "new" {
		t.Errorf("expected 'new', got %q", result["X"])
	}
	if result["Y"] != "1" {
		t.Errorf("expected Y=1")
	}
}
