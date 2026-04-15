// Package sync orchestrates reading secrets from Vault and writing them to a .env file.
package sync

import (
	"fmt"
	"log"

	"github.com/yourusername/vaultpipe/internal/config"
	"github.com/yourusername/vaultpipe/internal/envwriter"
	"github.com/yourusername/vaultpipe/internal/vault"
)

// Syncer coordinates fetching secrets from Vault and persisting them locally.
type Syncer struct {
	cfg    *config.Config
	client *vault.Client
	writer *envwriter.Writer
}

// New creates a Syncer wired up with the provided config.
func New(cfg *config.Config) (*Syncer, error) {
	client, err := vault.NewClient(cfg.VaultAddr, cfg.VaultToken)
	if err != nil {
		return nil, fmt.Errorf("syncer: initialising vault client: %w", err)
	}

	writer, err := envwriter.NewWriter(cfg.OutputFile)
	if err != nil {
		return nil, fmt.Errorf("syncer: initialising env writer: %w", err)
	}

	return &Syncer{
		cfg:    cfg,
		client: client,
		writer: writer,
	}, nil
}

// Run performs a single sync cycle: read secrets then write the .env file.
func (s *Syncer) Run() error {
	log.Printf("syncer: reading secrets from %s", s.cfg.SecretPath)

	secrets, err := s.client.ReadSecrets(s.cfg.SecretPath)
	if err != nil {
		return fmt.Errorf("syncer: reading secrets: %w", err)
	}

	log.Printf("syncer: fetched %d secret(s)", len(secrets))

	if err := s.writer.Write(secrets, s.cfg.Namespace); err != nil {
		return fmt.Errorf("syncer: writing env file: %w", err)
	}

	log.Printf("syncer: wrote %s (namespace filter: %q)", s.cfg.OutputFile, s.cfg.Namespace)
	return nil
}
