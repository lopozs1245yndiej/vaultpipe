// Package snapshot provides point-in-time captures of Vault secrets.
//
// Snapshots are saved as JSON files in a configurable directory and can be
// loaded, listed, and deleted by their unique ID. The package also exposes
// a Diff function to compare two snapshots and identify added, removed, and
// changed keys — useful for auditing secret drift over time.
//
// Basic usage:
//
//	m, err := snapshot.New("/var/lib/vaultpipe/snapshots")
//	snap, err := m.Save("vault/myapp", secrets)
//	loaded, err := m.Load(snap.ID)
//	result := snapshot.Diff(snap, loaded)
package snapshot
