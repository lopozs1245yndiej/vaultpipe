// Package sync orchestrates reading secrets from Vault and writing them
// to a local .env file.
package sync

import (
	"fmt"

	"github.com/your-org/vaultpipe/internal/audit"
	"github.com/your-org/vaultpipe/internal/config"
	"github.com/your-org/vaultpipe/internal/diff"
	"github.com/your-org/vaultpipe/internal/envwriter"
)

// VaultReader is the interface satisfied by vault.Client.
type VaultReader interface {
	ReadSecrets(path string) (map[string]string, error)
}

// Syncer coordinates a single sync operation.
type Syncer struct {
	cfg    *config.Config
	vault  VaultReader
	logger *audit.Logger
}

// New creates a Syncer with the provided dependencies.
func New(cfg *config.Config, vault VaultReader, logger *audit.Logger) *Syncer {
	return &Syncer{cfg: cfg, vault: vault, logger: logger}
}

// Run executes the sync: reads secrets, computes a diff, writes the .env
// file, and logs an audit entry.
func (s *Syncer) Run() error {
	secrets, err := s.vault.ReadSecrets(s.cfg.SecretPath)
	if err != nil {
		s.logEntry(nil, "failure", err)
		return fmt.Errorf("sync: read secrets: %w", err)
	}

	w, err := envwriter.NewWriter(s.cfg.OutputFile)
	if err != nil {
		s.logEntry(nil, "failure", err)
		return fmt.Errorf("sync: open output: %w", err)
	}

	changes, err := diff.Compare(s.cfg.OutputFile, secrets)
	if err != nil {
		// Non-fatal: diff failure should not block the write.
		changes = map[string]string{}
	}

	if err := w.Write(secrets, s.cfg.Namespace); err != nil {
		s.logEntry(changes, "failure", err)
		return fmt.Errorf("sync: write env: %w", err)
	}

	s.logEntry(changes, "success", nil)
	return nil
}

func (s *Syncer) logEntry(changes map[string]string, status string, err error) {
	if s.logger == nil {
		return
	}
	entry := audit.Entry{
		SecretPath: s.cfg.SecretPath,
		OutputFile: s.cfg.OutputFile,
		Namespace:  s.cfg.Namespace,
		Changes:    changes,
		Status:     status,
	}
	if err != nil {
		entry.Error = err.Error()
	}
	_ = s.logger.Log(entry)
}
