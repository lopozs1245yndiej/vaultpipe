// Package checkpoint provides persistent tracking of the last successful
// sync operation for each Vault secret path.
//
// A Tracker records an Entry per path containing:
//   - the UTC timestamp of the last successful sync
//   - a checksum (e.g. SHA-256 of the written .env content) for change detection
//   - an optional namespace label
//
// Entries are stored as a JSON file on disk and are safe for concurrent use.
//
// Example usage:
//
//	tr, err := checkpoint.New(".vaultpipe_checkpoint.json")
//	if err != nil { ... }
//
//	// After a successful sync:
//	_ = tr.Set(checkpoint.Entry{
//		Path:     "secret/myapp",
//		Checksum: hex.EncodeToString(hash[:]),
//	})
//
//	// Before syncing, check whether secrets have changed:
//	prev, err := tr.Get("secret/myapp")
package checkpoint
