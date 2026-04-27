package batch_test

import (
	"context"
	"testing"

	"github.com/vaultpipe/vaultpipe/internal/batch"
)

func noop(_ context.Context, _ map[string]string) error { return nil }

func TestNewFromOptions_DefaultsConcurrency(t *testing.T) {
	b, err := batch.NewFromOptions(batch.Options{Concurrency: 0}, noop)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b == nil {
		t.Fatal("expected non-nil Batcher")
	}
}

func TestNewFromOptions_ExplicitConcurrency(t *testing.T) {
	b, err := batch.NewFromOptions(batch.Options{Concurrency: 4}, noop)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b == nil {
		t.Fatal("expected non-nil Batcher")
	}
}

func TestNewFromOptions_NilFunc_ReturnsError(t *testing.T) {
	_, err := batch.NewFromOptions(batch.Options{Concurrency: 2}, nil)
	if err == nil {
		t.Fatal("expected error for nil ProcessFunc")
	}
}
