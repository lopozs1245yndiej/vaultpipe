package semaphore_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/yourusername/vaultpipe/internal/semaphore"
)

func TestNew_InvalidLimit(t *testing.T) {
	_, err := semaphore.New(0, 0)
	if err == nil {
		t.Fatal("expected error for limit=0")
	}
}

func TestNew_Valid(t *testing.T) {
	s, err := semaphore.New(3, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Available() != 3 {
		t.Fatalf("expected 3 available slots, got %d", s.Available())
	}
}

func TestAcquire_Release_RoundTrip(t *testing.T) {
	s, _ := semaphore.New(2, time.Second)

	if err := s.Acquire(context.Background()); err != nil {
		t.Fatalf("acquire 1: %v", err)
	}
	if err := s.Acquire(context.Background()); err != nil {
		t.Fatalf("acquire 2: %v", err)
	}
	if s.Available() != 0 {
		t.Fatalf("expected 0 available, got %d", s.Available())
	}
	s.Release()
	if s.Available() != 1 {
		t.Fatalf("expected 1 available after release, got %d", s.Available())
	}
}

func TestAcquire_Timeout(t *testing.T) {
	s, _ := semaphore.New(1, 50*time.Millisecond)

	// Fill the only slot.
	_ = s.Acquire(context.Background())

	err := s.Acquire(context.Background())
	if err != semaphore.ErrAcquireTimeout {
		t.Fatalf("expected ErrAcquireTimeout, got %v", err)
	}
}

func TestAcquire_ContextCancelled(t *testing.T) {
	s, _ := semaphore.New(1, 0)
	_ = s.Acquire(context.Background())

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := s.Acquire(ctx)
	if err == nil {
		t.Fatal("expected error from cancelled context")
	}
	if err == semaphore.ErrAcquireTimeout {
		t.Fatal("expected context error, not timeout")
	}
}

func TestAcquire_ConcurrentLimit(t *testing.T) {
	const limit = 3
	const goroutines = 10

	s, _ := semaphore.New(limit, 2*time.Second)

	var mu sync.Mutex
	var peak int
	var current int
	var wg sync.WaitGroup

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := s.Acquire(context.Background()); err != nil {
				return
			}
			defer s.Release()

			mu.Lock()
			current++
			if current > peak {
				peak = current
			}
			mu.Unlock()

			time.Sleep(10 * time.Millisecond)

			mu.Lock()
			current--
			mu.Unlock()
		}()
	}
	wg.Wait()

	if peak > limit {
		t.Fatalf("peak concurrency %d exceeded limit %d", peak, limit)
	}
}
