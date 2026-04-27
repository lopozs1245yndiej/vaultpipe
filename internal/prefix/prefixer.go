// Package prefix provides utilities for adding, removing, and replacing
// key prefixes in secret maps. It is useful when secrets from Vault need
// to be namespaced or re-namespaced before being written to .env files.
package prefix

import (
	"errors"
	"strings"
)

// ErrEmptyPrefix is returned when an empty prefix string is provided.
var ErrEmptyPrefix = errors.New("prefix: prefix must not be empty")

// Prefixer adds or removes a prefix from secret map keys.
type Prefixer struct {
	prefix    string
	separator string
}

// Option configures a Prefixer.
type Option func(*Prefixer)

// WithSeparator sets the separator placed between the prefix and the key.
// Defaults to "_".
func WithSeparator(sep string) Option {
	return func(p *Prefixer) {
		p.separator = sep
	}
}

// New creates a new Prefixer with the given prefix.
// Returns ErrEmptyPrefix if prefix is empty.
func New(prefix string, opts ...Option) (*Prefixer, error) {
	if strings.TrimSpace(prefix) == "" {
		return nil, ErrEmptyPrefix
	}
	p := &Prefixer{
		prefix:    prefix,
		separator: "_",
	}
	for _, o := range opts {
		o(p)
	}
	return p, nil
}

// Add returns a new map with the prefix prepended to every key.
func (p *Prefixer) Add(secrets map[string]string) map[string]string {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[p.prefix+p.separator+k] = v
	}
	return out
}

// Strip returns a new map with the prefix (and separator) removed from
// keys that carry it. Keys that do not match are passed through unchanged.
func (p *Prefixer) Strip(secrets map[string]string) map[string]string {
	out := make(map[string]string, len(secrets))
	target := p.prefix + p.separator
	for k, v := range secrets {
		if strings.HasPrefix(k, target) {
			out[strings.TrimPrefix(k, target)] = v
		} else {
			out[k] = v
		}
	}
	return out
}

// Replace swaps oldPrefix for newPrefix on every matching key.
// Keys that do not carry oldPrefix are passed through unchanged.
func (p *Prefixer) Replace(secrets map[string]string, newPrefix string) map[string]string {
	out := make(map[string]string, len(secrets))
	target := p.prefix + p.separator
	newTarget := newPrefix + p.separator
	for k, v := range secrets {
		if strings.HasPrefix(k, target) {
			out[newTarget+strings.TrimPrefix(k, target)] = v
		} else {
			out[k] = v
		}
	}
	return out
}
