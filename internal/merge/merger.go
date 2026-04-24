// Package merge provides functionality for merging multiple secret maps
// into a single map, with configurable conflict resolution strategies.
package merge

import "fmt"

// Strategy defines how key conflicts are resolved during a merge.
type Strategy int

const (
	// StrategyLastWins causes the last source to overwrite earlier values.
	StrategyLastWins Strategy = iota
	// StrategyFirstWins preserves the first value seen for a key.
	StrategyFirstWins
	// StrategyError returns an error when a key conflict is detected.
	StrategyError
)

// Merger merges multiple secret maps using a defined strategy.
type Merger struct {
	strategy Strategy
}

// New creates a new Merger with the given conflict resolution strategy.
func New(strategy Strategy) *Merger {
	return &Merger{strategy: strategy}
}

// Merge combines the provided secret maps into a single map.
// Sources are processed in order; conflict behaviour depends on the strategy.
func (m *Merger) Merge(sources ...map[string]string) (map[string]string, error) {
	result := make(map[string]string)

	for _, src := range sources {
		for k, v := range src {
			existing, exists := result[k]
			if !exists {
				result[k] = v
				continue
			}

			switch m.strategy {
			case StrategyLastWins:
				result[k] = v
			case StrategyFirstWins:
				// keep existing — do nothing
				_ = existing
			case StrategyError:
				return nil, fmt.Errorf("merge conflict: key %q appears in multiple sources", k)
			}
		}
	}

	return result, nil
}

// MergeAll is a convenience function that merges sources using StrategyLastWins.
func MergeAll(sources ...map[string]string) map[string]string {
	m := New(StrategyLastWins)
	result, _ := m.Merge(sources...)
	return result
}
