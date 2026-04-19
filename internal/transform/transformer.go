// Package transform provides key/value transformation pipelines for secrets.
package transform

import (
	"strings"
)

// TransformFunc is a function that transforms a key-value pair.
type TransformFunc func(key, value string) (string, string)

// Transformer applies a chain of transformations to secrets.
type Transformer struct {
	fns []TransformFunc
}

// New creates a new Transformer with the given transform functions.
func New(fns ...TransformFunc) *Transformer {
	return &Transformer{fns: fns}
}

// Apply runs all transformations over the provided secrets map and returns a new map.
func (t *Transformer) Apply(secrets map[string]string) map[string]string {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		for _, fn := range t.fns {
			k, v = fn(k, v)
		}
		out[k] = v
	}
	return out
}

// UppercaseKeys returns a TransformFunc that uppercases all keys.
func UppercaseKeys() TransformFunc {
	return func(key, value string) (string, string) {
		return strings.ToUpper(key), value
	}
}

// PrefixKeys returns a TransformFunc that prepends a prefix to all keys.
func PrefixKeys(prefix string) TransformFunc {
	return func(key, value string) (string, string) {
		return prefix + key, value
	}
}

// TrimValueSpace returns a TransformFunc that trims whitespace from values.
func TrimValueSpace() TransformFunc {
	return func(key, value string) (string, string) {
		return key, strings.TrimSpace(value)
	}
}

// ReplaceKeyChars returns a TransformFunc that replaces occurrences of old with new in keys.
func ReplaceKeyChars(old, new string) TransformFunc {
	return func(key, value string) (string, string) {
		return strings.ReplaceAll(key, old, new), value
	}
}
