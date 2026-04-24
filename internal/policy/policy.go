package policy

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// Rule defines a single access policy rule.
type Rule struct {
	Key     string // key pattern (supports * glob)
	Allowed bool   // true = allow, false = deny
}

// Policy enforces key-level access rules for secrets.
type Policy struct {
	rules []Rule
}

// New creates a Policy from a slice of rules.
// Rules are evaluated in order; the first match wins.
// If no rule matches, the key is allowed by default.
func New(rules []Rule) (*Policy, error) {
	for i, r := range rules {
		if strings.TrimSpace(r.Key) == "" {
			return nil, fmt.Errorf("rule %d: key pattern must not be empty", i)
		}
	}
	return &Policy{rules: rules}, nil
}

// IsAllowed reports whether the given key is permitted under the policy.
func (p *Policy) IsAllowed(key string) bool {
	for _, r := range p.rules {
		if matchGlob(r.Key, key) {
			return r.Allowed
		}
	}
	return true // default allow
}

// Apply filters a secrets map, returning only permitted keys.
func (p *Policy) Apply(secrets map[string]string) map[string]string {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		if p.IsAllowed(k) {
			out[k] = v
		}
	}
	return out
}

// Violations returns keys that are denied by the policy.
func (p *Policy) Violations(secrets map[string]string) []string {
	var denied []string
	for k := range secrets {
		if !p.IsAllowed(k) {
			denied = append(denied, k)
		}
	}
	return denied
}

// ErrDenied is returned when a required key is denied by policy.
var ErrDenied = errors.New("policy: key denied")

// matchGlob matches a key against a glob pattern (only * wildcard).
func matchGlob(pattern, key string) bool {
	regexPat := "^" + regexp.QuoteMeta(pattern) + "$"
	regexPat = strings.ReplaceAll(regexPat, `\*`, `.*`)
	re, err := regexp.Compile(regexPat)
	if err != nil {
		return false
	}
	return re.MatchString(key)
}
