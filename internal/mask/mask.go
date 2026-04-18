// Package mask provides utilities for redacting sensitive secret values
// before displaying them in logs, diffs, or terminal output.
package mask

import "strings"

const defaultMask = "********"

// Masker redacts secret values while optionally revealing a prefix.
type Masker struct {
	revealChars int
}

// New returns a Masker. revealChars controls how many leading characters of a
// value remain visible (0 means fully masked).
func New(revealChars int) *Masker {
	if revealChars < 0 {
		revealChars = 0
	}
	return &Masker{revealChars: revealChars}
}

// Mask redacts a single secret value.
func (m *Masker) Mask(value string) string {
	if value == "" {
		return ""
	}
	if m.revealChars == 0 || m.revealChars >= len(value) {
		return defaultMask
	}
	return value[:m.revealChars] + strings.Repeat("*", len(defaultMask))
}

// MaskMap returns a copy of secrets with all values redacted.
func (m *Masker) MaskMap(secrets map[string]string) map[string]string {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[k] = m.Mask(v)
	}
	return out
}
