// Package drain provides graceful shutdown coordination for the vaultpipe
// CLI. It ensures that any in-flight Vault sync operations are allowed to
// finish (or time out) before the process exits.
//
// Typical usage:
//
//	 d := drain.New(10 * time.Second)
//
//	 // Before starting work:
//	 if !d.Acquire() {
//	     return // shutting down
//	 }
//	 defer d.Release()
//
//	 // On SIGTERM / SIGINT:
//	 if err := d.Drain(ctx); err != nil {
//	     log.Printf("drain error: %v", err)
//	 }
package drain
