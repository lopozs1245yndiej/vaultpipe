// Package coalesce provides a mechanism for selecting the first non-empty
// value from a prioritised list of secret maps. This is useful when secrets
// may come from multiple sources (e.g. Vault namespaces, environment tiers)
// and a fallback chain is required.
package coalesce

import (
	"errors"
	"fmt"
)

// ErrAllEmpty is returned when every source map is empty or nil.
var ErrAllEmpty = errors.New("coalesce: all sources are empty")

// Strategy controls how values are selected when a key appears in multiple
// sources.
type Strategy int

const (
	// FirstNonEmpty picks the value from the first source that contains the key
	// with a non-empty string value.
	FirstNonEmpty Strategy = iota

	// FirstPresent picks the value from the first source that contains the key,
	// even if the value is an empty string.
	FirstPresent
)

// Coalescer merges a prioritised slice of secret maps into a single map.
type Coalescer struct {
	strategy Strategy
}

// New returns a Coalescer configured with the given Strategy.
func New(s Strategy) (*Coalescer, error) {
	if s != FirstNonEmpty && s != FirstPresent {
		return nil, fmt.Errorf("coalesce: unknown strategy %d", s)
	}
	return &Coalescer{strategy: s}, nil
}

// Apply merges sources in priority order (index 0 = highest priority).
// For each key found across all sources the strategy determines which value
// wins. Keys present in no source are omitted from the result.
//
// An error is returned only when every source is nil or empty and
// requireNonEmpty is true.
func (c *Coalescer) Apply(sources []map[string]string) (map[string]string, error) {
	if len(sources) == 0 {
		return nil, ErrAllEmpty
	}

	// Collect all keys across all sources.
	keySet := make(map[string]struct{})
	for _, src := range sources {
		for k := range src {
			keySet[k] = struct{}{}
		}
	}

	if len(keySet) == 0 {
		return nil, ErrAllEmpty
	}

	result := make(map[string]string, len(keySet))

	for key := range keySet {
		for _, src := range sources {
			val, ok := src[key]
			if !ok {
				continue
			}
			if c.strategy == FirstPresent || val != "" {
				result[key] = val
				break
			}
		}
	}

	return result, nil
}

// MustApply is like Apply but panics on error. Intended for tests or
// initialisation paths where the caller has already validated the sources.
func (c *Coalescer) MustApply(sources []map[string]string) map[string]string {
	out, err := c.Apply(sources)
	if err != nil {
		panic(err)
	}
	return out
}
