// Package scope provides namespace-scoped secret key filtering and
// remapping for vaultpipe. It allows callers to strip a namespace prefix
// from keys, optionally renaming them for the target environment.
package scope

import (
	"fmt"
	"strings"
)

// Scoper strips a namespace prefix from secret keys and optionally applies
// a replacement prefix to the resulting key names.
type Scoper struct {
	namespace      string
	replacePrefix  string
	caseSensitive  bool
}

// Option configures a Scoper.
type Option func(*Scoper)

// WithReplacePrefix sets a prefix to prepend to keys after the namespace is
// stripped. If empty, keys are left without any prefix.
func WithReplacePrefix(p string) Option {
	return func(s *Scoper) { s.replacePrefix = p }
}

// WithCaseSensitive controls whether namespace matching is case-sensitive.
// Default is case-insensitive.
func WithCaseSensitive(v bool) Option {
	return func(s *Scoper) { s.caseSensitive = v }
}

// New creates a Scoper that strips keys beginning with namespace.
// namespace must not be empty.
func New(namespace string, opts ...Option) (*Scoper, error) {
	if strings.TrimSpace(namespace) == "" {
		return nil, fmt.Errorf("scope: namespace must not be empty")
	}
	s := &Scoper{namespace: namespace}
	for _, o := range opts {
		o(s)
	}
	return s, nil
}

// Apply filters secrets to those whose keys begin with the configured
// namespace, strips the namespace prefix, and optionally prepends
// replacePrefix. Keys that do not match the namespace are excluded.
func (s *Scoper) Apply(secrets map[string]string) map[string]string {
	out := make(map[string]string, len(secrets))
	ns := s.namespace
	if !s.caseSensitive {
		ns = strings.ToLower(ns)
	}
	for k, v := range secrets {
		compare := k
		if !s.caseSensitive {
			compare = strings.ToLower(k)
		}
		if !strings.HasPrefix(compare, ns) {
			continue
		}
		trimmed := k[len(ns):]
		trimmed = strings.TrimLeft(trimmed, "_")
		if trimmed == "" {
			continue
		}
		newKey := trimmed
		if s.replacePrefix != "" {
			newKey = s.replacePrefix + trimmed
		}
		out[newKey] = v
	}
	return out
}
