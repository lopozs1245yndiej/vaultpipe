package sync

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/user/vaultpipe/internal/config"
	"github.com/user/vaultpipe/internal/rotate"
	"github.com/user/vaultpipe/internal/envwriter"
)

type mockVault struct {
	secrets map[string]string
	err     error
}

func (m *mockVault) ReadSecrets(_ context.Context, _ string) (map[string]string, error) {
	return m.secrets, m.err
}

func newTestSyncer(t *testing.T, secrets map[string]string, vaultErr error) (*Syncer, string) {
	t.Helper()
	tmpDir := t.TempDir()
	envPath := filepath.Join(tmpDir, ".env")
	cfg := &config.Config{
		VaultAddr:   "http://127.0.0.1:8200",
		VaultToken:  "tok",
		SecretPath:  "secret/data/app",
		EnvFilePath: envPath,
		BackupDir:   filepath.Join(tmpDir, "backups"),
		MaxBackups:  3,
	}
	s := &Syncer{
		cfg:     cfg,
		vault:   &mockVault{secrets: secrets, err: vaultErr},
		writer:  envwriter.NewWriter(envPath, ""),
		rotator: rotate.New(cfg.BackupDir, cfg.MaxBackups),
		auditor: nil,
	}
	return s, envPath
}

func TestRun_WritesEnvFile(t *testing.T) {
	s, envPath := newTestSyncer(t, map[string]string{"KEY": "value"}, nil)
	if err := s.Run(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(envPath)
	if string(data) == "" {
		t.Error("expected non-empty .env file")
	}
}

func TestRun_WithNamespaceFilter(t *testing.T) {
	tmpDir := t.TempDir()
	envPath := filepath.Join(tmpDir, ".env")
	cfg := &config.Config{
		SecretPath:  "secret/data/app",
		EnvFilePath: envPath,
		Namespace:   "APP",
		BackupDir:   filepath.Join(tmpDir, "backups"),
		MaxBackups:  3,
	}
	secrets := map[string]string{"APP_KEY": "v1", "OTHER_KEY": "v2"}
	s := &Syncer{
		cfg:     cfg,
		vault:   &mockVault{secrets: secrets},
		writer:  envwriter.NewWriter(envPath, "APP"),
		rotator: rotate.New(cfg.BackupDir, cfg.MaxBackups),
	}
	if err := s.Run(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(envPath)
	if !contains(string(data), "APP_KEY") {
		t.Error("expected APP_KEY in output")
	}
	if contains(string(data), "OTHER_KEY") {
		t.Error("unexpected OTHER_KEY in namespace-filtered output")
	}
}

func TestRun_VaultError_ReturnsError(t *testing.T) {
	s, _ := newTestSyncer(t, nil, errors.New("vault unavailable"))
	if err := s.Run(context.Background()); err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestRun_RotatesExistingFile(t *testing.T) {
	s, envPath := newTestSyncer(t, map[string]string{"K": "v"}, nil)
	_ = os.WriteFile(envPath, []byte("OLD=1\n"), 0600)
	if err := s.Run(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	matches, _ := filepath.Glob(filepath.Join(s.cfg.BackupDir, ".env.*.bak"))
	if len(matches) == 0 {
		t.Error("expected a backup file to be created")
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 ||
		(func() bool {
			for i := 0; i <= len(s)-len(sub); i++ {
				if s[i:i+len(sub)] == sub {
					return true
				}
			}
			return false
		})())
}
