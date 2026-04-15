// Package audit provides functionality for logging sync operations
// and tracking secret changes over time.
package audit

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Entry represents a single audit log entry for a sync operation.
type Entry struct {
	Timestamp  time.Time         `json:"timestamp"`
	SecretPath string            `json:"secret_path"`
	OutputFile string            `json:"output_file"`
	Namespace  string            `json:"namespace,omitempty"`
	Changes    map[string]string `json:"changes"`
	Status     string            `json:"status"`
	Error      string            `json:"error,omitempty"`
}

// Logger writes audit entries to a file or stdout.
type Logger struct {
	filePath string
	file     *os.File
}

// NewLogger creates a new audit Logger. If filePath is empty, logs are
// written to stdout.
func NewLogger(filePath string) (*Logger, error) {
	if filePath == "" {
		return &Logger{file: os.Stdout}, nil
	}
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return nil, fmt.Errorf("audit: open log file: %w", err)
	}
	return &Logger{filePath: filePath, file: f}, nil
}

// Log writes an audit entry as a JSON line to the underlying writer.
func (l *Logger) Log(entry Entry) error {
	if entry.Timestamp.IsZero() {
		entry.Timestamp = time.Now().UTC()
	}
	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("audit: marshal entry: %w", err)
	}
	_, err = fmt.Fprintf(l.file, "%s\n", data)
	if err != nil {
		return fmt.Errorf("audit: write entry: %w", err)
	}
	return nil
}

// Close releases resources held by the Logger.
func (l *Logger) Close() error {
	if l.file != nil && l.file != os.Stdout {
		return l.file.Close()
	}
	return nil
}
