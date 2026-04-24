// Package lock provides a PID-based file locking mechanism for vaultpipe.
//
// It prevents multiple concurrent vaultpipe processes from writing to the
// same .env file simultaneously, which could result in partial or corrupted
// output.
//
// Usage:
//
//	l := lock.New("/tmp/vaultpipe.lock")
//	if err := l.Acquire(); err != nil {
//		if errors.Is(err, lock.ErrLocked) {
//			log.Fatal("another vaultpipe process is running")
//		}
//		log.Fatalf("lock error: %v", err)
//	}
//	defer l.Release()
//
// Stale locks (where the owning process is no longer alive) are automatically
// removed on the next Acquire call.
package lock
