// Package schema provides secret schema validation for vaultpipe.
//
// A Schema is a collection of Field definitions that describe the expected
// keys, types, and constraints for a set of secrets fetched from Vault.
//
// Supported field types:
//
//	"string" – any non-empty string (default)
//	"int"    – a decimal integer, optionally negative
//	"bool"   – one of: true, false, 1, 0
//	"url"    – a string starting with http:// or https://
//
// Fields may also carry an optional regex Pattern that the value must match.
//
// Example:
//
//	s, err := schema.New([]schema.Field{
//		{Key: "DATABASE_URL", Type: schema.TypeURL, Required: true},
//		{Key: "PORT",         Type: schema.TypeInt},
//		{Key: "REGION",       Pattern: `^us-[a-z]+-\d+$`},
//	})
//	errs := s.Validate(secrets)
//
// Schemas can also be loaded from a JSON file via LoadFromFile.
package schema
