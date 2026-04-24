// Package tokenize provides a Tokenizer for splitting and reassembling
// structured secret values using a configurable delimiter.
//
// Some secrets stored in Vault encode multiple fields in a single string,
// for example a database DSN of the form "user:password@host:port".
// The Tokenizer makes it easy to extract individual components, replace
// specific tokens, and reassemble the value without losing the original
// structure.
//
// Example:
//
//	tok, err := tokenize.New(":")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	host, _ := tok.Token("localhost:5432", 0) // "localhost"
//	updated, _ := tok.Replace("localhost:5432", 1, "9999") // "localhost:9999"
package tokenize
