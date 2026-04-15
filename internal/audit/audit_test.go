package audit_test

import (
	"bufio"
	"encoding/json"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/your-org/vaultpipe/internal/audit"
)

func tmpAuditFile(t *testing.T) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "audit-*.jsonl")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestNewLogger_InvalidPath(t *testing.T) {
	_, err := audit.NewLogger("/nonexistent/dir/audit.log")
	if err == nil {
		t.Fatal("expected error for invalid path, got nil")
	}
}

func TestLog_WritesJSONLine(t *testing.T) {
	path := tmpAuditFile(t)
	logger, err := audit.NewLogger(path)
	if err != nil {
		t.Fatalf("NewLogger: %v", err)
	}
	defer logger.Close()

	entry := audit.Entry{
		Timestamp:  time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		SecretPath: "secret/data/myapp",
		OutputFile: ".env",
		Namespace:  "APP",
		Changes:    map[string]string{"APP_KEY": "added"},
		Status:     "success",
	}
	if err := logger.Log(entry); err != nil {
		t.Fatalf("Log: %v", err)
	}
	logger.Close()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read audit file: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if len(lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(lines))
	}
	var got audit.Entry
	if err := json.Unmarshal([]byte(lines[0]), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.SecretPath != "secret/data/myapp" {
		t.Errorf("SecretPath = %q, want %q", got.SecretPath, "secret/data/myapp")
	}
	if got.Status != "success" {
		t.Errorf("Status = %q, want %q", got.Status, "success")
	}
}

func TestLog_MultipleEntries(t *testing.T) {
	path := tmpAuditFile(t)
	logger, err := audit.NewLogger(path)
	if err != nil {
		t.Fatalf("NewLogger: %v", err)
	}
	defer logger.Close()

	for i := 0; i < 3; i++ {
		if err := logger.Log(audit.Entry{Status: "success"}); err != nil {
			t.Fatalf("Log[%d]: %v", i, err)
		}
	}
	logger.Close()

	f, _ := os.Open(path)
	defer f.Close()
	scanner := bufio.NewScanner(f)
	count := 0
	for scanner.Scan() {
		count++
	}
	if count != 3 {
		t.Errorf("expected 3 lines, got %d", count)
	}
}

func TestLog_SetsTimestampIfZero(t *testing.T) {
	path := tmpAuditFile(t)
	logger, _ := audit.NewLogger(path)
	defer logger.Close()

	entry := audit.Entry{Status: "success"}
	logger.Log(entry)
	logger.Close()

	data, _ := os.ReadFile(path)
	var got audit.Entry
	json.Unmarshal(data, &got)
	if got.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp to be set automatically")
	}
}
