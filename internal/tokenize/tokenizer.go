// Package tokenize provides utilities for splitting and reassembling
// secret values using configurable delimiters, useful when secrets
// contain structured data (e.g. "host:port" or "user@host").
package tokenize

import (
	"errors"
	"strings"
)

// ErrEmptyDelimiter is returned when an empty delimiter is provided.
var ErrEmptyDelimiter = errors.New("tokenize: delimiter must not be empty")

// ErrIndexOutOfRange is returned when the requested token index does not exist.
var ErrIndexOutOfRange = errors.New("tokenize: index out of range")

// Tokenizer splits and joins secret values by a configured delimiter.
type Tokenizer struct {
	delimiter string
}

// New creates a Tokenizer using the given delimiter.
func New(delimiter string) (*Tokenizer, error) {
	if delimiter == "" {
		return nil, ErrEmptyDelimiter
	}
	return &Tokenizer{delimiter: delimiter}, nil
}

// Split breaks value into a slice of tokens using the configured delimiter.
func (t *Tokenizer) Split(value string) []string {
	return strings.Split(value, t.delimiter)
}

// Token returns the token at position index from value.
// Returns ErrIndexOutOfRange if the index exceeds the number of tokens.
func (t *Tokenizer) Token(value string, index int) (string, error) {
	parts := t.Split(value)
	if index < 0 || index >= len(parts) {
		return "", ErrIndexOutOfRange
	}
	return parts[index], nil
}

// Join reassembles tokens into a single string using the configured delimiter.
func (t *Tokenizer) Join(tokens []string) string {
	return strings.Join(tokens, t.delimiter)
}

// Replace replaces the token at position index with replacement and returns
// the reassembled string. Returns ErrIndexOutOfRange if index is invalid.
func (t *Tokenizer) Replace(value string, index int, replacement string) (string, error) {
	parts := t.Split(value)
	if index < 0 || index >= len(parts) {
		return "", ErrIndexOutOfRange
	}
	parts[index] = replacement
	return t.Join(parts), nil
}

// SplitMap applies Split to every value in secrets and returns a map of
// key -> []string tokens.
func (t *Tokenizer) SplitMap(secrets map[string]string) map[string][]string {
	out := make(map[string][]string, len(secrets))
	for k, v := range secrets {
		out[k] = t.Split(v)
	}
	return out
}
