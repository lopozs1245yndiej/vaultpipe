package drain_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/yourusername/vaultpipe/internal/drain"
)

func TestNew_DefaultTimeout(t *testing.T) {
	d := drain.New(0)
	if d == nil {
		t.Fatal("expected non-nil Drainer")
	}
}

func TestAcquire_ReturnsFalseAfterClose(t *testing.T) {
	d := drain.New(time.Second)
	if !d.Acquire() {
		t.Fatal("expected first Acquire to succeed")
	}
	d.Release()

	_ = d.Drain(context.Background())

	if d.Acquire() {
		t.Fatal("expected Acquire to fail after Drain")
	}
}

func TestDrain_WaitsForRelease(t *testing.T) {
	d := drain.New(5 * time.Second)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if !d.Acquire() {
			return
		}
		time.Sleep(50 * time.Millisecond)
		d.Release()
	}()

	// Give goroutine time to acquire.
	time.Sleep(10 * time.Millisecond)

	err := d.Drain(context.Background())
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	wg.Wait()
}

func TestDrain_ReturnsErrDrainTimeout(t *testing.T) {
	d := drain.New(50 * time.Millisecond)

	if !d.Acquire() {
		t.Fatal("expected Acquire to succeed")
	}
	// Intentionally never Release to trigger timeout.
	defer d.Release()

	err := d.Drain(context.Background())
	if err != drain.ErrDrainTimeout {
		t.Fatalf("expected ErrDrainTimeout, got %v", err)
	}
}

func TestDrain_RespectsContextCancellation(t *testing.T) {
	d := drain.New(10 * time.Second)

	if !d.Acquire() {
		t.Fatal("expected Acquire to succeed")
	}
	defer d.Release()

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	err := d.Drain(ctx)
	if err == nil {
		t.Fatal("expected context error, got nil")
	}
}

func TestIsClosed_BeforeAndAfterDrain(t *testing.T) {
	d := drain.New(time.Second)

	if d.IsClosed() {
		t.Fatal("expected IsClosed to be false before Drain")
	}

	_ = d.Drain(context.Background())

	if !d.IsClosed() {
		t.Fatal("expected IsClosed to be true after Drain")
	}
}

func TestDrain_MultipleWorkers(t *testing.T) {
	d := drain.New(5 * time.Second)
	const workers = 10

	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if !d.Acquire() {
				return
			}
			defer d.Release()
			time.Sleep(20 * time.Millisecond)
		}()
	}

	time.Sleep(5 * time.Millisecond)
	err := d.Drain(context.Background())
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	wg.Wait()
}
