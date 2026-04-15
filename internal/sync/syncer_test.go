package sync_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yourusername/vaultpipe/internal/config"
	"github.com/yourusername/vaultpipe/internal/sync"
)

func newMockVault(t *testing.T, data map[string]interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		payload := map[string]interface{}{
			"data": map[string]interface{}{"data": data},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(payload)
	}))
}

func TestRun_WritesEnvFile(t *testing.T) {
	srv := newMockVault(t, map[string]interface{}{
		"APP_KEY": "secret123",
		"APP_DEBUG": "false",
	})
	defer srv.Close()

	out := filepath.Join(t.TempDir(), ".env")

	cfg := &config.Config{
		VaultAddr:  srv.URL,
		VaultToken: "test-token",
		SecretPath: "secret/data/myapp",
		OutputFile: out,
		Namespace:  "",
	}

	s, err := sync.New(cfg)
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}

	if err := s.Run(); err != nil {
		t.Fatalf("Run() error: %v", err)
	}

	bytes, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("reading output file: %v", err)
	}

	content := string(bytes)
	if !strings.Contains(content, "APP_KEY") {
		t.Errorf("expected APP_KEY in output, got:\n%s", content)
	}
}

func TestRun_WithNamespaceFilter(t *testing.T) {
	srv := newMockVault(t, map[string]interface{}{
		"APP_KEY":  "secret123",
		"DB_PASS":  "dbpass",
	})
	defer srv.Close()

	out := filepath.Join(t.TempDir(), ".env")

	cfg := &config.Config{
		VaultAddr:  srv.URL,
		VaultToken: "test-token",
		SecretPath: "secret/data/myapp",
		OutputFile: out,
		Namespace:  "APP",
	}

	s, err := sync.New(cfg)
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}

	if err := s.Run(); err != nil {
		t.Fatalf("Run() error: %v", err)
	}

	bytes, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("reading output file: %v", err)
	}

	content := string(bytes)
	if strings.Contains(content, "DB_PASS") {
		t.Errorf("DB_PASS should be filtered out by namespace APP, got:\n%s", content)
	}
	if !strings.Contains(content, "APP_KEY") {
		t.Errorf("expected APP_KEY in output, got:\n%s", content)
	}
}
