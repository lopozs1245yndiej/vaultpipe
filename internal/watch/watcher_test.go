package watch_test

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/yourusername/vaultpipe/internal/watch"
)

type mockSyncer struct {
	calls int
	err   error
}

func (m *mockSyncer) Run(_ context.Context) error {
	m.calls++
	return m.err
}

func newWatcher(s watch.Syncable, interval time.Duration) *watch.Watcher {
	return watch.New(s, interval, log.Default())
}

func TestRun_CallsSyncOnTick(t *testing.T) {
	ms := &mockSyncer{}
	w := newWatcher(ms, 50*time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 180*time.Millisecond)
	defer cancel()

	_ = w.Run(ctx)

	if ms.calls < 2 {
		t.Errorf("expected at least 2 sync calls, got %d", ms.calls)
	}
}

func TestRun_StopsOnContextCancel(t *testing.T) {
	ms := &mockSyncer{}
	w := newWatcher(ms, 10*time.Second)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	err := w.Run(ctx)
	if err != context.Canceled {
		t.Errorf("expected context.Canceled, got %v", err)
	}
}

func TestRun_ContinuesOnSyncError(t *testing.T) {
	ms := &mockSyncer{err: fmt.Errorf("vault down")}
	w := newWatcher(ms, 40*time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 130*time.Millisecond)
	defer cancel()

	_ = w.Run(ctx)
	if ms.calls < 2 {
		t.Errorf("expected retries despite error, got %d calls", ms.calls)
	}
}

func TestRun_ZeroCallsBeforeFirstTick(t *testing.T) {
	ms := &mockSyncer{}
	// Use a very long interval so the ticker never fires within the timeout.
	w := newWatcher(ms, 10*time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	_ = w.Run(ctx)
	if ms.calls != 0 {
		t.Errorf("expected 0 sync calls before first tick, got %d", ms.calls)
	}
}
