package watch_test

import (
	"context"
	"fmt"
	"log"
	"sync/atomic"
	"testing"
	"time"

	"github.com/yourusername/vaultpipe/internal/watch"
)

type countSyncer struct {
	count int64
}

func (c *countSyncer) Run(_ context.Context) error {
	atomic.AddInt64(&c.count, 1)
	return nil
}

func TestWatcher_TickCount(t *testing.T) {
	cs := &countSyncer{}
	w := watch.New(cs, 30*time.Millisecond, log.New(log.Writer(), "[test] ", 0))

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	_ = w.Run(ctx)

	got := atomic.LoadInt64(&cs.count)
	if got < 4 {
		t.Errorf("expected >=4 ticks in 200ms at 30ms interval, got %d", got)
	}
	fmt.Printf("integration: syncer called %d times\n", got)
}
