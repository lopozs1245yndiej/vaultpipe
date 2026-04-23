// Package redact provides a Redactor type that scrubs known secret values from
// arbitrary strings, log lines, and error messages before they are displayed or
// persisted.
//
// Usage:
//
//	r := redact.New("[REDACTED]")
//	r.Load(secrets)          // register values fetched from Vault
//	safe := r.Redact(line)   // scrub a log line before writing
//
// RedactMap can be used to produce a display-safe copy of a secrets map where
// every value is replaced with the placeholder, useful for audit logging or
// dry-run output.
package redact
