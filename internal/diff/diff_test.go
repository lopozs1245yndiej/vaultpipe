package diff

import (
	"os"
	"testing"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "vaultpipe-diff-*.env")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestCompare_Added(t *testing.T) {
	path := writeTempEnv(t, "EXISTING_KEY=value1\n")
	incoming := map[string]string{
		"EXISTING_KEY": "value1",
		"NEW_KEY":      "value2",
	}
	changes, err := Compare(path, incoming)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(changes) != 1 || changes[0].Key != "NEW_KEY" || changes[0].Type != Added {
		t.Errorf("expected one Added change for NEW_KEY, got %+v", changes)
	}
}

func TestCompare_Removed(t *testing.T) {
	path := writeTempEnv(t, "OLD_KEY=value\nKEEP_KEY=keep\n")
	incoming := map[string]string{
		"KEEP_KEY": "keep",
	}
	changes, err := Compare(path, incoming)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(changes) != 1 || changes[0].Key != "OLD_KEY" || changes[0].Type != Removed {
		t.Errorf("expected one Removed change for OLD_KEY, got %+v", changes)
	}
}

func TestCompare_Changed(t *testing.T) {
	path := writeTempEnv(t, "MY_KEY=old_value\n")
	incoming := map[string]string{
		"MY_KEY": "new_value",
	}
	changes, err := Compare(path, incoming)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(changes) != 1 || changes[0].Key != "MY_KEY" || changes[0].Type != Changed {
		t.Errorf("expected one Changed change for MY_KEY, got %+v", changes)
	}
}

func TestCompare_NoChanges(t *testing.T) {
	path := writeTempEnv(t, "KEY_A=alpha\nKEY_B=beta\n")
	incoming := map[string]string{
		"KEY_A": "alpha",
		"KEY_B": "beta",
	}
	changes, err := Compare(path, incoming)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(changes) != 0 {
		t.Errorf("expected no changes, got %+v", changes)
	}
}

func TestCompare_MissingFile(t *testing.T) {
	incoming := map[string]string{"NEW_KEY": "value"}
	changes, err := Compare("/nonexistent/path/.env", incoming)
	if err != nil {
		t.Fatalf("unexpected error for missing file: %v", err)
	}
	if len(changes) != 1 || changes[0].Type != Added {
		t.Errorf("expected Added change for missing file, got %+v", changes)
	}
}
