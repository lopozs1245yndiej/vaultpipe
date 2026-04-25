// Package dedupe provides secret deduplication by detecting and removing
// duplicate key-value pairs across multiple secret maps.
package dedupe

import (
	"crypto/sha256"
	"fmt"
)

// Strategy controls how duplicates are handled.
type Strategy int

const (
	// StrategyKeepFirst retains the first occurrence of a duplicate key.
	StrategyKeepFirst Strategy = iota
	// StrategyKeepLast retains the last occurrence of a duplicate key.
	StrategyKeepLast
	// StrategyError returns an error when a duplicate key is detected.
	StrategyError
)

// Deduper removes duplicate secrets according to a chosen strategy.
type Deduper struct {
	strategy Strategy
}

// New creates a new Deduper with the given strategy.
func New(strategy Strategy) *Deduper {
	return &Deduper{strategy: strategy}
}

// Apply merges multiple secret maps, resolving duplicates per the strategy.
// Sources are processed in order; earlier indices are considered "first".
func (d *Deduper) Apply(sources ...map[string]string) (map[string]string, error) {
	result := make(map[string]string)
	seen := make(map[string]bool)

	for _, src := range sources {
		for k, v := range src {
			if seen[k] {
				switch d.strategy {
				case StrategyKeepFirst:
					continue
				case StrategyKeepLast:
					result[k] = v
				case StrategyError:
					return nil, fmt.Errorf("dedupe: duplicate key %q", k)
				}
			} else {
				result[k] = v
				seen[k] = true
			}
		}
	}

	return result, nil
}

// Fingerprint returns a stable SHA-256 hex digest of a secret map's
// key-value content, useful for change detection.
func Fingerprint(secrets map[string]string) string {
	h := sha256.New()
	// Iterate sorted for determinism.
	keys := make([]string, 0, len(secrets))
	for k := range secrets {
		keys = append(keys, k)
	}
	sortStrings(keys)
	for _, k := range keys {
		fmt.Fprintf(h, "%s=%s\n", k, secrets[k])
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

// sortStrings sorts a string slice in place (insertion sort — small N).
func sortStrings(s []string) {
	for i := 1; i < len(s); i++ {
		for j := i; j > 0 && s[j] < s[j-1]; j-- {
			s[j], s[j-1] = s[j-1], s[j]
		}
	}
}
