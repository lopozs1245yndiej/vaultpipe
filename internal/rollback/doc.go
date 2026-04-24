// Package rollback provides functionality for restoring .env files
// from previously created backup snapshots.
//
// A Rollbacker scans a backup directory for *.env.bak files and can
// restore any of them to a target path. It integrates naturally with
// the rotate package, which produces the backup files, and the diff
// package, which can show what changed between versions.
//
// Basic usage:
//
//	r, err := rollback.New("/var/vaultpipe/backups")
//	if err != nil { ... }
//
//	// Restore the most recent backup
//	entry, err := r.Latest()
//	if err != nil { ... }
//	if err := r.Restore(entry.ID, ".env"); err != nil { ... }
package rollback
