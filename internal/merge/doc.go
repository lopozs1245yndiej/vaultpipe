// Package merge provides utilities for combining multiple secret maps
// retrieved from HashiCorp Vault into a single unified map.
//
// When secrets are sourced from multiple Vault paths or namespaces they may
// share key names. The Merger type lets callers choose how those conflicts
// are resolved:
//
//   - StrategyLastWins  — the last source wins (default, mirrors shell behaviour)
//   - StrategyFirstWins — the first value seen is kept
//   - StrategyError     — an error is returned on the first conflict
//
// Example:
//
//	m := merge.New(merge.StrategyLastWins)
//	combined, err := m.Merge(baseSecrets, overrideSecrets)
package merge
