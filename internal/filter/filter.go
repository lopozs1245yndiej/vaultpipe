// Package filter provides key filtering utilities for vaultpipe,
// allowing inclusion and exclusion of secrets by glob-style patterns.
package filter

import (
	"path"
	"strings"
)

// Filter holds compiled include/exclude patterns.
type Filter struct {
	include []string
	exclude []string
}

// New creates a Filter with the given include and exclude glob patterns.
// An empty include list means "include all".
func New(include, exclude []string) *Filter {
	return &Filter{
		include: include,
		exclude: exclude,
	}
}

// Apply returns a filtered copy of secrets, keeping only keys that match
// at least one include pattern (if any) and no exclude patterns.
func (f *Filter) Apply(secrets map[string]string) map[string]string {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		if f.isExcluded(k) {
			continue
		}
		if !f.isIncluded(k) {
			continue
		}
		out[k] = v
	}
	return out
}

func (f *Filter) isIncluded(key string) bool {
	if len(f.include) == 0 {
		return true
	}
	for _, pat := range f.include {
		if matchPattern(pat, key) {
			return true
		}
	}
	return false
}

func (f *Filter) isExcluded(key string) bool {
	for _, pat := range f.exclude {
		if matchPattern(pat, key) {
			return true
		}
	}
	return false
}

func matchPattern(pattern, key string) bool {
	// Case-insensitive prefix shorthand: patterns without wildcards match as prefix
	if !strings.ContainsAny(pattern, "*?[") {
		return strings.HasPrefix(key, pattern)
	}
	matched, err := path.Match(pattern, key)
	return err == nil && matched
}
