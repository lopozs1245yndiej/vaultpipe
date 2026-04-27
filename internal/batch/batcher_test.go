package batch_test

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"

	"github.com/vaultpipe/vaultpipe/internal/batch"
)

func TestNew_InvalidConcurrency(t *testing.T) {
	_, err := batch.New(0, func(_ context.Context, _ map[string]string) error { return nil })
	if err == nil {
		t.Fatal("expected error for concurrency=0")
	}
}

func TestNew_NilFunc(t *testing.T) {
	_, err := batch.New(2, nil)
	if err == nil {
		t.Fatal("expected error for nil ProcessFunc")
	}
}

func TestNew_Valid(t *testing.T) {
	b, err := batch.New(2, func(_ context.Context, _ map[string]string) error { return nil })
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b == nil {
		t.Fatal("expected non-nil Batcher")
	}
}

func TestRun_EmptyItems(t *testing.T) {
	b, _ := batch.New(2, func(_ context.Context, _ map[string]string) error { return nil })
	_, err := b.Run(context.Background(), nil)
	if !errors.Is(err, batch.ErrEmptyBatch) {
		t.Fatalf("expected ErrEmptyBatch, got %v", err)
	}
}

func TestRun_ProcessesAllItems(t *testing.T) {
	var count atomic.Int32
	b, _ := batch.New(3, func(_ context.Context, _ map[string]string) error {
		count.Add(1)
		return nil
	})

	items := []map[string]string{
		{"A": "1"}, {"B": "2"}, {"C": "3"}, {"D": "4"},
	}

	results, err := b.Run(context.Background(), items)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if int(count.Load()) != len(items) {
		t.Errorf("expected %d calls, got %d", len(items), count.Load())
	}
	if len(results) != len(items) {
		t.Errorf("expected %d results, got %d", len(items), len(results))
	}
}

func TestRun_CollectsErrors(t *testing.T) {
	sentinel := errors.New("proc error")
	b, _ := batch.New(2, func(_ context.Context, m map[string]string) error {
		if _, ok := m["fail"]; ok {
			return sentinel
		}
		return nil
	})

	items := []map[string]string{
		{"ok": "1"},
		{"fail": "yes"},
		{"ok": "2"},
	}

	results, err := b.Run(context.Background(), items)
	if err != nil {
		t.Fatalf("unexpected top-level error: %v", err)
	}
	errs := batch.Errors(results)
	if len(errs) != 1 {
		t.Fatalf("expected 1 error result, got %d", len(errs))
	}
	if !errors.Is(errs[0].Err, sentinel) {
		t.Errorf("expected sentinel error, got %v", errs[0].Err)
	}
}

func TestRun_CancelledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	b, _ := batch.New(1, func(_ context.Context, _ map[string]string) error {
		return nil
	})

	items := []map[string]string{{"x": "1"}, {"y": "2"}, {"z": "3"}}
	results, err := b.Run(ctx, items)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	errs := batch.Errors(results)
	if len(errs) == 0 {
		t.Error("expected at least one context-cancelled error")
	}
}
