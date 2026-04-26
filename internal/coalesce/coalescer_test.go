package coalesce_test

import (
	"testing"

	"github.com/yourusername/vaultpipe/internal/coalesce"
)

func TestNew_NilSources_ReturnsError(t *testing.T) {
	_, err := coalesce.New(nil)
	if err == nil {
		t.Fatal("expected error for nil sources, got nil")
	}
}

func TestNew_EmptySources_ReturnsError(t *testing.T) {
	_, err := coalesce.New([]map[string]string{})
	if err == nil {
		t.Fatal("expected error for empty sources, got nil")
	}
}

func TestNew_ValidSources_OK(t *testing.T) {
	sources := []map[string]string{
		{"KEY": "value"},
	}
	_, err := coalesce.New(sources)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestApply_FirstNonEmptyWins(t *testing.T) {
	sources := []map[string]string{
		{"DB_HOST": "", "APP_ENV": "production"},
		{"DB_HOST": "localhost", "APP_ENV": "staging"},
		{"DB_HOST": "remote", "APP_ENV": "development"},
	}
	c, err := coalesce.New(sources)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	result, err := c.Apply()
	if err != nil {
		t.Fatalf("unexpected error applying coalesce: %v", err)
	}
	if got := result["DB_HOST"]; got != "localhost" {
		t.Errorf("DB_HOST: expected %q, got %q", "localhost", got)
	}
	if got := result["APP_ENV"]; got != "production" {
		t.Errorf("APP_ENV: expected %q, got %q", "production", got)
	}
}

func TestApply_AllEmptyValues_ReturnsEmpty(t *testing.T) {
	sources := []map[string]string{
		{"KEY": ""},
		{"KEY": ""},
	}
	c, err := coalesce.New(sources)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	result, err := c.Apply()
	if err != nil {
		t.Fatalf("unexpected error applying coalesce: %v", err)
	}
	if got, ok := result["KEY"]; ok && got != "" {
		t.Errorf("expected empty or absent KEY, got %q", got)
	}
}

func TestApply_KeyOnlyInLaterSource(t *testing.T) {
	sources := []map[string]string{
		{"A": "first"},
		{"A": "second", "B": "only-in-second"},
	}
	c, err := coalesce.New(sources)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	result, err := c.Apply()
	if err != nil {
		t.Fatalf("unexpected error applying coalesce: %v", err)
	}
	if got := result["A"]; got != "first" {
		t.Errorf("A: expected %q, got %q", "first", got)
	}
	if got := result["B"]; got != "only-in-second" {
		t.Errorf("B: expected %q, got %q", "only-in-second", got)
	}
}

func TestApply_SingleSource_ReturnsItself(t *testing.T) {
	source := map[string]string{
		"FOO": "bar",
		"BAZ": "qux",
	}
	c, err := coalesce.New([]map[string]string{source})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	result, err := c.Apply()
	if err != nil {
		t.Fatalf("unexpected error applying coalesce: %v", err)
	}
	for k, v := range source {
		if got := result[k]; got != v {
			t.Errorf("%s: expected %q, got %q", k, v, got)
		}
	}
}

func TestApply_DoesNotMutateSource(t *testing.T) {
	original := map[string]string{"X": "original"}
	sources := []map[string]string{
		original,
		{"X": "override"},
	}
	c, err := coalesce.New(sources)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, err = c.Apply()
	if err != nil {
		t.Fatalf("unexpected error applying coalesce: %v", err)
	}
	if original["X"] != "original" {
		t.Errorf("source map was mutated: expected %q, got %q", "original", original["X"])
	}
}
