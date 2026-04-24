// Package notify provides post-sync notification hooks for vaultpipe,
// allowing callers to be alerted when secrets are written or changed.
package notify

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Event represents a sync notification payload.
type Event struct {
	Timestamp time.Time         `json:"timestamp"`
	SecretPath string           `json:"secret_path"`
	KeysChanged []string        `json:"keys_changed"`
	Meta        map[string]string `json:"meta,omitempty"`
}

// Notifier dispatches sync events to a webhook endpoint.
type Notifier struct {
	endpoint string
	client   *http.Client
	headers  map[string]string
}

// New creates a Notifier that posts events to the given endpoint.
// Optional headers (e.g. Authorization) can be supplied via headers map.
func New(endpoint string, headers map[string]string) (*Notifier, error) {
	if endpoint == "" {
		return nil, fmt.Errorf("notify: endpoint must not be empty")
	}
	return &Notifier{
		endpoint: endpoint,
		client:   &http.Client{Timeout: 10 * time.Second},
		headers:  headers,
	}, nil
}

// Send serialises the Event as JSON and POSTs it to the configured endpoint.
// A non-2xx response is treated as an error.
func (n *Notifier) Send(ctx context.Context, evt Event) error {
	if evt.Timestamp.IsZero() {
		evt.Timestamp = time.Now().UTC()
	}

	body, err := json.Marshal(evt)
	if err != nil {
		return fmt.Errorf("notify: marshal event: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, n.endpoint, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("notify: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	for k, v := range n.headers {
		req.Header.Set(k, v)
	}

	resp, err := n.client.Do(req)
	if err != nil {
		return fmt.Errorf("notify: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("notify: unexpected status %d from %s", resp.StatusCode, n.endpoint)
	}
	return nil
}
