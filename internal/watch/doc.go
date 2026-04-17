// Package watch provides a polling watcher that periodically re-syncs secrets
// from HashiCorp Vault into local .env files.
//
// Usage:
//
//	s, _ := sync.New(cfg)
//	w := watch.New(s, 30*time.Second, nil)
//	if err := w.Run(ctx); err != nil && err != context.Canceled {
//		log.Fatal(err)
//	}
//
// The watcher respects context cancellation, making it safe to use with
// OS signal handlers for graceful shutdown.
package watch
