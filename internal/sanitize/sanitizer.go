// Package sanitize provides utilities for cleaning and normalizing secret
// key-value pairs before they are written to .env files or passed downstream.
package sanitize

import (
	"strings"
	"unicode"
)

// Sanitizer cleans secret maps according to configured rules.
type Sanitizer struct {
	stripControlChars bool
	normalizeKeys     bool
	trimValues        bool
}

// Option configures a Sanitizer.
type Option func(*Sanitizer)

// WithStripControlChars removes non-printable control characters from values.
func WithStripControlChars() Option {
	return func(s *Sanitizer) { s.stripControlChars = true }
}

// WithNormalizeKeys uppercases keys and replaces hyphens and dots with underscores.
func WithNormalizeKeys() Option {
	return func(s *Sanitizer) { s.normalizeKeys = true }
}

// WithTrimValues trims leading and trailing whitespace from values.
func WithTrimValues() Option {
	return func(s *Sanitizer) { s.trimValues = true }
}

// New returns a Sanitizer configured with the given options.
func New(opts ...Option) *Sanitizer {
	s := &Sanitizer{}
	for _, o := range opts {
		o(s)
	}
	return s
}

// Apply sanitizes a copy of the provided secrets map and returns it.
func (s *Sanitizer) Apply(secrets map[string]string) map[string]string {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		if s.normalizeKeys {
			k = normalizeKey(k)
		}
		if s.trimValues {
			v = strings.TrimSpace(v)
		}
		if s.stripControlChars {
			v = stripControl(v)
		}
		out[k] = v
	}
	return out
}

func normalizeKey(k string) string {
	k = strings.ToUpper(k)
	k = strings.ReplaceAll(k, "-", "_")
	k = strings.ReplaceAll(k, ".", "_")
	return k
}

func stripControl(s string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsControl(r) && r != '\t' {
			return -1
		}
		return r
	}, s)
}
