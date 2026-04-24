// Package policy provides key-level access control for secrets managed by
// vaultpipe.
//
// A Policy is a prioritised list of Rules. Each Rule associates a key pattern
// (supporting a single * glob wildcard) with an allow/deny decision. Rules are
// evaluated in declaration order and the first matching rule wins. If no rule
// matches, the key is allowed by default.
//
// Policies can be constructed programmatically via New, or loaded from a JSON
// file with LoadFromFile:
//
//	{
//	  "rules": [
//	    {"key": "INTERNAL_*", "allowed": false},
//	    {"key": "*",          "allowed": true}
//	  ]
//	}
//
// Use Apply to filter a secrets map to only permitted keys, and Violations to
// audit which keys in a map would be denied.
package policy
