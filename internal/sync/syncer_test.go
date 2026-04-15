package sync_test

import (
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/your-org/vaultpipe/internal/audit"
	"github.com/your-org/vaultpipe/internal/config"
	"github.com/your-org/vaultpipe/internal/sync"
)

type mockVault struct {
	secrets map[string]string
	err     error
}

func (m *mockVault) ReadSecrets(_ string) (map[string]string, error) {
	return m.secrets, m.err
}

func newMockVault(secrets map[string]string, err error) *mockVault {
	return &mockVault{secrets: secrets, err: err}
}

func TestRun_WritesEnvFile(t *testing.T) {
	out := t.TempDir() + "/.env"
	cfg := &config.Config{
		SecretPath: "secret/data/app",
		OutputFile: out,
	}
	logger, _ := audit.NewLogger("")
	defer logger.Close()

	s := sync.New(cfg, newMockVault(map[string]string{"KEY": "val"}, nil), logger)
	if err := s.Run(); err != nil {
		t.Fatalf("Run: %v", err)
	}
	data, _ := os.ReadFile(out)
	if !strings.Contains(string(data), "KEY=val") {
		t.Errorf("expected KEY=val in output, got: %s", data)
	}
}

func TestRun_WithNamespaceFilter(t *testing.T) {
	out := t.TempDir() + "/.env"
	cfg := &config.Config{
		SecretPath: "secret/data/app",
		OutputFile: out,
		Namespace:  "APP",
	}
	logger, _ := audit.NewLogger("")
	defer logger.Close()

	secrets := map[string]string{"APP_KEY": "v1", "OTHER_KEY": "v2"}
	s := sync.New(cfg, newMockVault(secrets, nil), logger)
	if err := s.Run(); err != nil {
		t.Fatalf("Run: %v", err)
	}
	data, _ := os.ReadFile(out)
	if !strings.Contains(string(data), "APP_KEY=v1") {
		t.Errorf("expected APP_KEY in output")
	}
	if strings.Contains(string(data), "OTHER_KEY") {
		t.Errorf("unexpected OTHER_KEY in filtered output")
	}
}

func TestRun_VaultError_ReturnsError(t *testing.T) {
	out := t.TempDir() + "/.env"
	cfg := &config.Config{
		SecretPath: "secret/data/app",
		OutputFile: out,
	}
	logger, _ := audit.NewLogger("")
	defer logger.Close()

	s := sync.New(cfg, newMockVault(nil, errors.New("vault unavailable")), logger)
	if err := s.Run(); err == nil {
		t.Fatal("expected error from vault failure, got nil")
	}
}
