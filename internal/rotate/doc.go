// Package rotate provides backup and rotation functionality for .env files
// managed by vaultpipe.
//
// Before vaultpipe overwrites a local .env file with fresh secrets from Vault,
// the Rotator creates a timestamped backup copy in a configurable backup
// directory. Once the number of backups exceeds the configured maximum, the
// oldest backups are pruned automatically.
//
// Example usage:
//
//	r := rotate.New(".vaultpipe/backups", 5)
//	if err := r.Rotate(".env"); err != nil {
//	    log.Fatal(err)
//	}
package rotate
