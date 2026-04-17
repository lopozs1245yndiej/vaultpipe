package watch

import (
	"context"
	"log"
	"time"

	"github.com/yourusername/vaultpipe/internal/sync"
)

// Watcher polls Vault at a fixed interval and re-syncs secrets.
type Watcher struct {
	syncer   *sync.Syncer
	interval time.Duration
	logger   *log.Logger
}

// New creates a new Watcher with the given syncer and poll interval.
func New(s *sync.Syncer, interval time.Duration, logger *log.Logger) *Watcher {
	if logger == nil {
		logger = log.Default()
	}
	return &Watcher{
		syncer:   s,
		interval: interval,
		logger:   logger,
	}
}

// Run starts the watch loop, blocking until ctx is cancelled.
func (w *Watcher) Run(ctx context.Context) error {
	w.logger.Printf("watch: starting poll every %s", w.interval)
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			w.logger.Println("watch: stopping")
			return ctx.Err()
		case <-ticker.C:
			if err := w.syncer.Run(ctx); err != nil {
				w.logger.Printf("watch: sync error: %v", err)
			} else {
				w.logger.Println("watch: sync completed")
			}
		}
	}
}
