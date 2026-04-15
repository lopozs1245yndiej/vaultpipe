package vault

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newMockVaultServer(t *testing.T, responseBody map[string]interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(responseBody); err != nil {
			t.Errorf("failed to encode mock response: %v", err)
		}
	}))
}

func TestNewClient_InvalidAddress(t *testing.T) {
	_, err := NewClient("://bad-address", "token", "")
	if err == nil {
		t.Fatal("expected error for invalid address, got nil")
	}
}

func TestReadSecrets_KVv2(t *testing.T) {
	payload := map[string]interface{}{
		"data": map[string]interface{}{
			"data": map[string]interface{}{
				"DB_PASSWORD": "supersecret",
				"API_KEY":     "abc123",
			},
		},
	}

	server := newMockVaultServer(t, payload)
	defer server.Close()

	client, err := NewClient(server.URL, "test-token", "")
	if err != nil {
		t.Fatalf("unexpected error creating client: %v", err)
	}

	secrets, err := client.ReadSecrets(context.Background(), "secret/data/myapp")
	if err != nil {
		t.Fatalf("unexpected error reading secrets: %v", err)
	}

	if secrets["DB_PASSWORD"] != "supersecret" {
		t.Errorf("expected DB_PASSWORD=supersecret, got %q", secrets["DB_PASSWORD"])
	}
	if secrets["API_KEY"] != "abc123" {
		t.Errorf("expected API_KEY=abc123, got %q", secrets["API_KEY"])
	}
}

func TestReadSecrets_WithPrefix(t *testing.T) {
	payload := map[string]interface{}{
		"data": map[string]interface{}{
			"TOKEN": "xyz",
		},
	}

	server := newMockVaultServer(t, payload)
	defer server.Close()

	client, err := NewClient(server.URL, "test-token", "prod/")
	if err != nil {
		t.Fatalf("unexpected error creating client: %v", err)
	}

	secrets, err := client.ReadSecrets(context.Background(), "myapp")
	if err != nil {
		t.Fatalf("unexpected error reading secrets: %v", err)
	}

	if secrets["TOKEN"] != "xyz" {
		t.Errorf("expected TOKEN=xyz, got %q", secrets["TOKEN"])
	}
}

func TestFlattenData(t *testing.T) {
	raw := map[string]interface{}{
		"KEY_A": "value_a",
		"KEY_B": 42,
	}

	result, err := flattenData(raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result["KEY_A"] != "value_a" {
		t.Errorf("expected KEY_A=value_a, got %q", result["KEY_A"])
	}
	if result["KEY_B"] != "42" {
		t.Errorf("expected KEY_B=42, got %q", result["KEY_B"])
	}
}
