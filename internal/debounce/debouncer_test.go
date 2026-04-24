package debounce_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/yourusername/vaultpipe/internal/debounce"
)

func TestNew_InvalidDelay(t *testing.T) {
	_, err := debounce.New(0)
	if err == nil {
		t.Fatal("expected error for zero delay, got nil")
	}
	_, err = debounce.New(-1 * time.Millisecond)
	if err == nil {
		t.Fatal("expected error for negative delay, got nil")
	}
}

func TestNew_Valid(t *testing.T) {
	d, err := debounce.New(10 * time.Millisecond)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d == nil {
		t.Fatal("expected non-nil debouncer")
	}
}

func TestCall_ExecutesFnAfterDelay(t *testing.T) {
	d, _ := debounce.New(20 * time.Millisecond)
	var called int32
	d.Call(context.Background(), func() {
		atomic.StoreInt32(&called, 1)
	})
	time.Sleep(50 * time.Millisecond)
	if atomic.LoadInt32(&called) != 1 {
		t.Fatal("expected fn to be called after delay")
	}
}

func TestCall_DebouncesRapidCalls(t *testing.T) {
	d, _ := debounce.New(40 * time.Millisecond)
	var count int32
	for i := 0; i < 5; i++ {
		d.Call(context.Background(), func() {
			atomic.AddInt32(&count, 1)
		})
		time.Sleep(5 * time.Millisecond)
	}
	time.Sleep(100 * time.Millisecond)
	if c := atomic.LoadInt32(&count); c != 1 {
		t.Fatalf("expected fn called once, got %d", c)
	}
}

func TestCall_CancelledContext_DoesNotExecute(t *testing.T) {
	d, _ := debounce.New(30 * time.Millisecond)
	var called int32
	ctx, cancel := context.WithCancel(context.Background())
	d.Call(ctx, func() {
		atomic.StoreInt32(&called, 1)
	})
	cancel()
	time.Sleep(60 * time.Millisecond)
	if atomic.LoadInt32(&called) != 0 {
		t.Fatal("expected fn NOT to be called after context cancel")
	}
}

func TestFlush_CancelsPendingCall(t *testing.T) {
	d, _ := debounce.New(50 * time.Millisecond)
	var called int32
	d.Call(context.Background(), func() {
		atomic.StoreInt32(&called, 1)
	})
	d.Flush()
	time.Sleep(80 * time.Millisecond)
	if atomic.LoadInt32(&called) != 0 {
		t.Fatal("expected fn NOT to be called after Flush")
	}
}
