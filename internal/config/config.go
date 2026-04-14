package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// Config holds all configuration for vaultpipe.
type Config struct {
	VaultAddr  string   `mapstructure:"vault_addr"`
	VaultToken string   `mapstructure:"vault_token"`
	Namespace  string   `mapstructure:"namespace"`
	SecretPath string   `mapstructure:"secret_path"`
	OutputFile string   `mapstructure:"output_file"`
	Filter     []string `mapstructure:"filter"`
}

// Load reads configuration from a config file and environment variables.
// Environment variables take precedence over config file values.
func Load(cfgFile string) (*Config, error) {
	v := viper.New()

	v.SetDefault("vault_addr", "http://127.0.0.1:8200")
	v.SetDefault("output_file", ".env")

	v.SetEnvPrefix("VAULTPIPE")
	v.AutomaticEnv()

	// Also respect the standard Vault env vars.
	v.BindEnv("vault_addr", "VAULT_ADDR")   //nolint:errcheck
	v.BindEnv("vault_token", "VAULT_TOKEN") //nolint:errcheck

	if cfgFile != "" {
		v.SetConfigFile(cfgFile)
	} else {
		v.SetConfigName(".vaultpipe")
		v.SetConfigType("yaml")
		v.AddConfigPath(".")
		v.AddConfigPath(os.ExpandEnv("$HOME"))
	}

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("reading config file: %w", err)
		}
		// Config file not found is acceptable; rely on env vars / flags.
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshalling config: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (c *Config) validate() error {
	if c.VaultAddr == "" {
		return fmt.Errorf("vault_addr must not be empty")
	}
	if c.VaultToken == "" {
		return fmt.Errorf("vault_token is required (set VAULT_TOKEN or vault_token in config)")
	}
	if c.SecretPath == "" {
		return fmt.Errorf("secret_path is required")
	}
	return nil
}
