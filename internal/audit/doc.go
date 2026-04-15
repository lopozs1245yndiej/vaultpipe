// Package audit provides structured audit logging for vaultpipe sync
// operations.
//
// Each sync run can produce one or more [Entry] values that capture which
// secrets were read, which output file was written, what namespace filter
// was applied, and a summary of changes (added / removed / changed keys).
//
// Entries are serialised as newline-delimited JSON (NDJSON) so that the
// audit log can be consumed by standard log-aggregation tooling such as
// Loki, Splunk, or a simple `jq` pipeline.
//
// Usage:
//
//	logger, err := audit.NewLogger("/var/log/vaultpipe/audit.jsonl")
//	if err != nil { ... }
//	defer logger.Close()
//
//	logger.Log(audit.Entry{
//		SecretPath: "secret/data/myapp",
//		OutputFile: ".env",
//		Status:     "success",
//		Changes:    map[string]string{"DB_PASS": "changed"},
//	})
package audit
