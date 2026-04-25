// Package truncate provides utilities for truncating secret values
// to a maximum byte or character length before writing or displaying them.
package truncate

import (
	"errors"
	"unicode/utf8"
)

// ErrInvalidMaxLen is returned when MaxLen is not a positive integer.
var ErrInvalidMaxLen = errors.New("truncate: maxLen must be greater than zero")

// Truncator truncates string values in a secrets map to a configurable
// maximum rune length, appending an optional suffix (e.g. "...") when
// truncation occurs.
type Truncator struct {
	maxLen int
	suffix string
}

// New creates a Truncator that caps each value at maxLen runes.
// suffix is appended to truncated values; pass an empty string for none.
func New(maxLen int, suffix string) (*Truncator, error) {
	if maxLen <= 0 {
		return nil, ErrInvalidMaxLen
	}
	return &Truncator{maxLen: maxLen, suffix: suffix}, nil
}

// Apply returns a new map where every value longer than maxLen runes is
// truncated. Values within the limit are copied unchanged.
func (t *Truncator) Apply(secrets map[string]string) map[string]string {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[k] = t.truncate(v)
	}
	return out
}

func (t *Truncator) truncate(s string) string {
	if utf8.RuneCountInString(s) <= t.maxLen {
		return s
	}
	// Build truncated string rune-by-rune up to maxLen.
	suffixLen := utf8.RuneCountInString(t.suffix)
	keep := t.maxLen - suffixLen
	if keep <= 0 {
		// suffix alone already fills the budget; just return the suffix.
		return t.suffix
	}
	runes := []rune(s)
	return string(runes[:keep]) + t.suffix
}
