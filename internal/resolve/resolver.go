// Package resolve provides variable interpolation for secret values,
// allowing secrets to reference other secrets using ${KEY} syntax.
package resolve

import (
	"fmt"
	"regexp"
	"strings"
)

var refPattern = regexp.MustCompile(`\$\{([A-Z0-9_]+)\}`)

// Resolver performs variable interpolation across a map of secrets.
type Resolver struct {
	maxDepth int
}

// New returns a Resolver with the given maximum interpolation depth.
// maxDepth prevents infinite loops caused by circular references.
func New(maxDepth int) (*Resolver, error) {
	if maxDepth < 1 {
		return nil, fmt.Errorf("resolve: maxDepth must be at least 1, got %d", maxDepth)
	}
	return &Resolver{maxDepth: maxDepth}, nil
}

// Resolve performs interpolation on all values in secrets, replacing
// ${KEY} references with the resolved value of KEY from the same map.
// Returns an error if a circular reference or undefined reference is found.
func (r *Resolver) Resolve(secrets map[string]string) (map[string]string, error) {
	result := make(map[string]string, len(secrets))
	for k, v := range secrets {
		resolved, err := r.resolveValue(v, secrets, 0)
		if err != nil {
			return nil, fmt.Errorf("resolve: key %q: %w", k, err)
		}
		result[k] = resolved
	}
	return result, nil
}

func (r *Resolver) resolveValue(value string, secrets map[string]string, depth int) (string, error) {
	if depth > r.maxDepth {
		return "", fmt.Errorf("max interpolation depth %d exceeded (possible circular reference)", r.maxDepth)
	}
	if !strings.Contains(value, "${") {
		return value, nil
	}
	var resolveErr error
	result := refPattern.ReplaceAllStringFunc(value, func(match string) string {
		if resolveErr != nil {
			return ""
		}
		key := refPattern.FindStringSubmatch(match)[1]
		ref, ok := secrets[key]
		if !ok {
			resolveErr = fmt.Errorf("undefined reference ${%s}", key)
			return ""
		}
		resolved, err := r.resolveValue(ref, secrets, depth+1)
		if err != nil {
			resolveErr = err
			return ""
		}
		return resolved
	})
	if resolveErr != nil {
		return "", resolveErr
	}
	return result, nil
}
