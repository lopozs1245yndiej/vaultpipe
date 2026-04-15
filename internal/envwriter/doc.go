// Package envwriter provides utilities for serialising Vault secrets into
// .env files suitable for use with docker-compose, direnv, or any
// twelve-factor application.
//
// # Namespace filtering
//
// When a namespace is supplied to NewWriter, only secret keys whose
// upper-cased names begin with "<NAMESPACE>_" are written to the output
// file.  This lets a single Vault path serve multiple services without
// each service seeing the other's credentials.
//
// # File safety
//
// Output files are always written with mode 0600 so that secrets are not
// readable by other users on the same system.
package envwriter
