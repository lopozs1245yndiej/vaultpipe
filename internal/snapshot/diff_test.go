package snapshot_test

import (
	"testing"
	"time"

	"github.com/your-org/vaultpipe/internal/snapshot"
)

func makeSnap(secrets map[string]string) *snapshot.Snapshot {
	return &snapshot.Snapshot{
		ID:        "test",
		CreatedAt: time.Now(),
		Secrets:   secrets,
		Source:    "vault/test",
	}
}

func TestDiff_Added(t *testing.T) {
	old := makeSnap(map[string]string{"A": "1"})
	new := makeSnap(map[string]string{"A": "1", "B": "2"})
	result := snapshot.Diff(old, new)
	if result.Added["B"] != "2" {
		t.Errorf("expected B=2 in Added, got %v", result.Added)
	}
	if result.HasChanges() == false {
		t.Error("expected HasChanges to be true")
	}
}

func TestDiff_Removed(t *testing.T) {
	old := makeSnap(map[string]string{"A": "1", "B": "2"})
	new := makeSnap(map[string]string{"A": "1"})
	result := snapshot.Diff(old, new)
	if result.Removed["B"] != "2" {
		t.Errorf("expected B in Removed, got %v", result.Removed)
	}
}

func TestDiff_Changed(t *testing.T) {
	old := makeSnap(map[string]string{"A": "old"})
	new := makeSnap(map[string]string{"A": "new"})
	result := snapshot.Diff(old, new)
	pair, ok := result.Changed["A"]
	if !ok {
		t.Fatal("expected A in Changed")
	}
	if pair[0] != "old" || pair[1] != "new" {
		t.Errorf("changed pair: got %v", pair)
	}
}

func TestDiff_NoChanges(t *testing.T) {
	old := makeSnap(map[string]string{"A": "1", "B": "2"})
	new := makeSnap(map[string]string{"A": "1", "B": "2"})
	result := snapshot.Diff(old, new)
	if result.HasChanges() {
		t.Error("expected no changes")
	}
}
