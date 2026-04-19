package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad_MissingToken(t *testing.T) {
	t.Setenv("VAULT_TOKEN", "")
	t.Setenv("VAULTPIPE_VAULT_TOKEN", "")

	_, err := Load("")
	if err == nil {
		t.Fatal("expected error when vault_token is missing, got nil")
	}
}

func TestLoad_MissingSecretPath(t *testing.T) {
	t.Setenv("VAULT_TOKEN", "test-token")
	t.Setenv("VAULTPIPE_SECRET_PATH", "")

	_, err := Load("")
	if err == nil {
		t.Fatal("expected error when secret_path is missing, got nil")
	}
}

func TestLoad_FromEnv(t *testing.T) {
	t.Setenv("VAULT_TOKEN", "s.testtoken")
	t.Setenv("VAULT_ADDR", "http://localhost:8200")
	t.Setenv("VAULTPIPE_SECRET_PATH", "secret/data/myapp")
	t.Setenv("VAULTPIPE_NAMESPACE", "dev")

	cfg, err := Load("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.VaultToken != "s.testtoken" {
		t.Errorf("expected token 's.testtoken', got %q", cfg.VaultToken)
	}
	if cfg.SecretPath != "secret/data/myapp" {
		t.Errorf("expected secret_path 'secret/data/myapp', got %q", cfg.SecretPath)
	}
	if cfg.Namespace != "dev" {
		t.Errorf("expected namespace 'dev', got %q", cfg.Namespace)
	}
	if cfg.OutputFile != ".env" {
		t.Errorf("expected default output_file '.env', got %q", cfg.OutputFile)
	}
}

func TestLoad_FromFile(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "vaultpipe.yaml")

	content := []byte(`
vault_addr: "http://vault.example.com:8200"
vault_token: "file-token"
secret_path: "secret/data/service"
output_file: "secrets.env"
filter:
  - DB_
  - API_
`)
	if err := os.WriteFile(cfgPath, content, 0o600); err != nil {
		t.Fatalf("failed to write temp config: %v", err)
	}

	cfg, err := Load(cfgPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.VaultAddr != "http://vault.example.com:8200" {
		t.Errorf("unexpected vault_addr: %q", cfg.VaultAddr)
	}
	if cfg.OutputFile != "secrets.env" {
		t.Errorf("unexpected output_file: %q", cfg.OutputFile)
	}
	if len(cfg.Filter) != 2 {
		t.Errorf("expected 2 filters, got %d", len(cfg.Filter))
	}
}

func TestLoad_NonExistentFile(t *testing.T) {
	_, err := Load("/nonexistent/path/vaultpipe.yaml")
	if err == nil {
		t.Fatal("expected error when config file does not exist, got nil")
	}
}
