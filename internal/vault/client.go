package vault

import (
	"context"
	"fmt"
	"strings"

	vaultapi "github.com/hashicorp/vault/api"
)

// Client wraps the Vault API client with helper methods.
type Client struct {
	api    *vaultapi.Client
	prefix string
}

// NewClient creates a new Vault client using the provided address and token.
func NewClient(address, token, prefix string) (*Client, error) {
	cfg := vaultapi.DefaultConfig()
	cfg.Address = address

	api, err := vaultapi.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create vault client: %w", err)
	}

	api.SetToken(token)

	return &Client{
		api:    api,
		prefix: strings.TrimSuffix(prefix, "/"),
	}, nil
}

// ReadSecrets reads key-value secrets from the given path in Vault.
// It supports KV v2 by looking under the "data" sub-key.
func (c *Client) ReadSecrets(ctx context.Context, path string) (map[string]string, error) {
	fullPath := path
	if c.prefix != "" {
		fullPath = c.prefix + "/" + strings.TrimPrefix(path, "/")
	}

	secret, err := c.api.Logical().ReadWithContext(ctx, fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read secret at %q: %w", fullPath, err)
	}
	if secret == nil {
		return nil, fmt.Errorf("no secret found at path %q", fullPath)
	}

	data, ok := secret.Data["data"]
	if !ok {
		// KV v1 — data is directly in secret.Data
		return flattenData(secret.Data)
	}

	kvData, ok := data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected data format at path %q", fullPath)
	}

	return flattenData(kvData)
}

func flattenData(raw map[string]interface{}) (map[string]string, error) {
	out := make(map[string]string, len(raw))
	for k, v := range raw {
		switch val := v.(type) {
		case string:
			out[k] = val
		case fmt.Stringer:
			out[k] = val.String()
		default:
			out[k] = fmt.Sprintf("%v", val)
		}
	}
	return out, nil
}
