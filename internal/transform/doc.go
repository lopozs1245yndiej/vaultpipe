// Package transform provides a composable pipeline for transforming secret
// key-value pairs before they are written to .env files or exported.
//
// Transformations are applied in order and can be chained:
//
//	tr := transform.New(
//		transform.ReplaceKeyChars("-", "_"),
//		transform.UppercaseKeys(),
//		transform.PrefixKeys("APP_"),
//		transform.TrimValueSpace(),
//	)
//	out := tr.Apply(secrets)
package transform
