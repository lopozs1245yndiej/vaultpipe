// Package redact provides utilities for redacting sensitive secret values
// from log output, error messages, and terminal display.
package redact

import "strings"

// Redactor holds a set of known secret values and can scrub them from strings.
type Redactor struct {
	values []string
	placeholder string
}

// New creates a new Redactor with the given placeholder string used to replace
// sensitive values. If placeholder is empty, "[REDACTED]" is used.
func New(placeholder string) *Redactor {
	if placeholder == "" {
		placeholder = "[REDACTED]"
	}
	return &Redactor{placeholder: placeholder}
}

// Load registers secret values that should be redacted. Empty strings are
// ignored to avoid inadvertently replacing all empty substrings.
func (r *Redactor) Load(secrets map[string]string) {
	for _, v := range secrets {
		if v != "" {
			r.values = append(r.values, v)
		}
	}
}

// Redact replaces all registered secret values found in s with the placeholder.
func (r *Redactor) Redact(s string) string {
	for _, v := range r.values {
		s = strings.ReplaceAll(s, v, r.placeholder)
	}
	return s
}

// RedactMap returns a copy of the provided map where every value has been
// replaced with the placeholder. Keys are preserved as-is.
func (r *Redactor) RedactMap(secrets map[string]string) map[string]string {
	out := make(map[string]string, len(secrets))
	for k := range secrets {
		out[k] = r.placeholder
	}
	return out
}

// Placeholder returns the string used to replace redacted values.
func (r *Redactor) Placeholder() string {
	return r.placeholder
}
