// Package pipeline provides a composable, ordered processing pipeline for
// secret maps read from HashiCorp Vault.
//
// A Pipeline is constructed from one or more Stage functions.  Each Stage
// receives the current secret map and returns a (possibly modified) copy.
// Stages are executed sequentially; the output of stage N becomes the input
// of stage N+1.
//
// Built-in stages:
//
//   - UppercaseValues – converts all values to upper-case.
//   - PrefixKeys      – prepends a fixed string to every key.
//   - FilterKeys      – retains only keys with a given prefix.
//   - TrimValueSpace  – strips leading/trailing whitespace from values.
//
// Custom stages can be provided as any function with the signature:
//
//	func(ctx context.Context, secrets map[string]string) (map[string]string, error)
//
// Example:
//
//	p := pipeline.New(
//		pipeline.TrimValueSpace(),
//		pipeline.PrefixKeys("APP_"),
//	)
//	result, err := p.Run(ctx, rawSecrets)
package pipeline
