package policy

import (
	"encoding/json"
	"fmt"
	"os"
)

// FileRule is the JSON-serialisable form of a Rule.
type FileRule struct {
	Key     string `json:"key"`
	Allowed bool   `json:"allowed"`
}

// policyFile is the top-level JSON structure.
type policyFile struct {
	Rules []FileRule `json:"rules"`
}

// LoadFromFile reads a JSON policy file and returns a Policy.
//
// Example file:
//
//	{
//	  "rules": [
//	    {"key": "INTERNAL_*", "allowed": false},
//	    {"key": "*",          "allowed": true}
//	  ]
//	}
func LoadFromFile(path string) (*Policy, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("policy: read file %q: %w", path, err)
	}

	var pf policyFile
	if err := json.Unmarshal(data, &pf); err != nil {
		return nil, fmt.Errorf("policy: parse %q: %w", path, err)
	}

	rules := make([]Rule, len(pf.Rules))
	for i, r := range pf.Rules {
		rules[i] = Rule{Key: r.Key, Allowed: r.Allowed}
	}

	return New(rules)
}

// MustLoad loads a policy from file and panics on error.
// Intended for use in tests or CLI initialisation.
func MustLoad(path string) *Policy {
	p, err := LoadFromFile(path)
	if err != nil {
		panic(err)
	}
	return p
}
