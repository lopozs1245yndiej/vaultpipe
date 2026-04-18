// Package validate provides secret key validation utilities for vaultpipe.
// It checks that secret keys conform to env variable naming conventions
// and optionally enforces a required-keys list.
package validate

import (
	"fmt"
	"regexp"
	"strings"
)

var validKeyRe = regexp.MustCompile(`^[A-Z][A-Z0-9_]*$`)

// Result holds the outcome of a validation run.
type Result struct {
	InvalidKeys  []string
	MissingKeys  []string
	Valid        bool
}

// Validator validates secret keys against naming rules and required keys.
type Validator struct {
	requiredKeys []string
}

// New returns a Validator. requiredKeys may be nil or empty.
func New(requiredKeys []string) *Validator {
	return &Validator{requiredKeys: requiredKeys}
}

// Validate checks secrets for invalid key names and missing required keys.
func (v *Validator) Validate(secrets map[string]string) Result {
	var invalid, missing []string

	for k := range secrets {
		upper := strings.ToUpper(k)
		if !validKeyRe.MatchString(upper) {
			invalid = append(invalid, k)
		}
	}

	for _, req := range v.requiredKeys {
		if _, ok := secrets[req]; !ok {
			missing = append(missing, req)
		}
	}

	return Result{
		InvalidKeys: invalid,
		MissingKeys: missing,
		Valid:        len(invalid) == 0 && len(missing) == 0,
	}
}

// FormatErrors returns a human-readable summary of validation failures.
func FormatErrors(r Result) string {
	if r.Valid {
		return ""
	}
	var sb strings.Builder
	if len(r.InvalidKeys) > 0 {
		sb.WriteString(fmt.Sprintf("invalid key names: %s\n", strings.Join(r.InvalidKeys, ", ")))
	}
	if len(r.MissingKeys) > 0 {
		sb.WriteString(fmt.Sprintf("missing required keys: %s\n", strings.Join(r.MissingKeys, ", ")))
	}
	return strings.TrimRight(sb.String(), "\n")
}
