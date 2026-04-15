package sync

import (
	"context"
	"fmt"

	"github.com/user/vaultpipe/internal/audit"
	"github.com/user/vaultpipe/internal/config"
	"github.com/user/vaultpipe/internal/envwriter"
	"github.com/user/vaultpipe/internal/rotate"
	"github.com/user/vaultpipe/internal/vault"
)

// VaultReader is the interface for reading secrets from Vault.
type VaultReader interface {
	ReadSecrets(ctx context.Context, path string) (map[string]string, error)
}

// Syncer orchestrates reading from Vault and writing to an .env file.
type Syncer struct {
	cfg      *config.Config
	vault    VaultReader
	writer   *envwriter.Writer
	rotator  *rotate.Rotator
	auditor  *audit.Logger
}

// New creates a Syncer wired with a real Vault client and env writer.
func New(cfg *config.Config, auditor *audit.Logger) (*Syncer, error) {
	client, err := vault.NewClient(cfg.VaultAddr, cfg.VaultToken)
	if err != nil {
		return nil, fmt.Errorf("sync: create vault client: %w", err)
	}
	writer := envwriter.NewWriter(cfg.EnvFilePath, cfg.Namespace)
	rotator := rotate.New(cfg.BackupDir, cfg.MaxBackups)
	return &Syncer{
		cfg:     cfg,
		vault:   client,
		writer:  writer,
		rotator: rotator,
		auditor: auditor,
	}, nil
}

// Run reads secrets from Vault, rotates the existing .env file, and writes
// the new secrets. It logs the operation via the audit logger when provided.
func (s *Syncer) Run(ctx context.Context) error {
	secrets, err := s.vault.ReadSecrets(ctx, s.cfg.SecretPath)
	if err != nil {
		return fmt.Errorf("sync: read secrets: %w", err)
	}

	if err := s.rotator.Rotate(s.cfg.EnvFilePath); err != nil {
		return fmt.Errorf("sync: rotate env file: %w", err)
	}

	if err := s.writer.Write(secrets); err != nil {
		return fmt.Errorf("sync: write env file: %w", err)
	}

	if s.auditor != nil {
		_ = s.auditor.Log(audit.Entry{
			Operation:  "sync",
			SecretPath: s.cfg.SecretPath,
			EnvFile:    s.cfg.EnvFilePath,
			KeyCount:   len(secrets),
		})
	}

	return nil
}
