package ratelimit_test

import (
	"context"
	"testing"
	"time"

	"github.com/yourusername/vaultpipe/internal/ratelimit"
)

func TestNew_InvalidRate(t *testing.T) {
	_, err := ratelimit.New(0, time.Second)
	if err == nil {
		t.Fatal("expected error for zero rate")
	}
}

func TestNew_InvalidWindow(t *testing.T) {
	_, err := ratelimit.New(5, 0)
	if err == nil {
		t.Fatal("expected error for zero window")
	}
}

func TestNew_Valid(t *testing.T) {
	l, err := ratelimit.New(10, time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if l == nil {
		t.Fatal("expected non-nil limiter")
	}
}

func TestWait_ConsumesTokens(t *testing.T) {
	l, _ := ratelimit.New(3, time.Second)
	ctx := context.Background()
	for i := 0; i < 3; i++ {
		if err := l.Wait(ctx); err != nil {
			t.Fatalf("unexpected error on wait %d: %v", i, err)
		}
	}
	if got := l.Remaining(); got != 0 {
		t.Fatalf("expected 0 remaining, got %d", got)
	}
}

func TestWait_CancelledContext(t *testing.T) {
	l, _ := ratelimit.New(1, 10*time.Second)
	ctx := context.Background()
	// exhaust tokens
	_ = l.Wait(ctx)

	ctx2, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	err := l.Wait(ctx2)
	if err == nil {
		t.Fatal("expected error when context cancelled")
	}
}

func TestRemaining_AfterReset(t *testing.T) {
	l, _ := ratelimit.New(5, time.Second)
	ctx := context.Background()
	_ = l.Wait(ctx)
	_ = l.Wait(ctx)
	l.Reset()
	if got := l.Remaining(); got != 5 {
		t.Fatalf("expected 5 after reset, got %d", got)
	}
}

func TestRemaining_ReturnsFullOnExpiredWindow(t *testing.T) {
	l, _ := ratelimit.New(4, 10*time.Millisecond)
	ctx := context.Background()
	_ = l.Wait(ctx)
	time.Sleep(20 * time.Millisecond)
	if got := l.Remaining(); got != 4 {
		t.Fatalf("expected 4 after window expiry, got %d", got)
	}
}
