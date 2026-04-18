// Package export provides multi-format secret export functionality for vaultpipe.
//
// It supports writing secrets to files in the following formats:
//
//   - env:  KEY=VALUE pairs, one per line
//   - json: Pretty-printed JSON object
//   - yaml: Simple YAML key-value pairs
//
// Files are written with 0600 permissions to prevent unauthorised access.
//
// Example usage:
//
//	ex, err := export.New(export.FormatJSON)
//	if err != nil { ... }
//	err = ex.Export(secrets, "/path/to/output.json")
package export
