package notify_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/your-org/vaultpipe/internal/notify"
)

func newTestServer(t *testing.T, statusCode int, received *notify.Event) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if received != nil {
			_ = json.NewDecoder(r.Body).Decode(received)
		}
		w.WriteHeader(statusCode)
	}))
}

func TestNew_EmptyEndpoint(t *testing.T) {
	_, err := notify.New("", nil)
	if err == nil {
		t.Fatal("expected error for empty endpoint")
	}
}

func TestSend_PostsJSONPayload(t *testing.T) {
	var received notify.Event
	srv := newTestServer(t, http.StatusOK, &received)
	defer srv.Close()

	n, err := notify.New(srv.URL, nil)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	evt := notify.Event{
		Timestamp:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		SecretPath:  "secret/app",
		KeysChanged: []string{"DB_PASS", "API_KEY"},
	}

	if err := n.Send(context.Background(), evt); err != nil {
		t.Fatalf("Send: %v", err)
	}
	if received.SecretPath != "secret/app" {
		t.Errorf("expected secret_path=secret/app, got %q", received.SecretPath)
	}
	if len(received.KeysChanged) != 2 {
		t.Errorf("expected 2 keys_changed, got %d", len(received.KeysChanged))
	}
}

func TestSend_SetsTimestampIfZero(t *testing.T) {
	var received notify.Event
	srv := newTestServer(t, http.StatusOK, &received)
	defer srv.Close()

	n, _ := notify.New(srv.URL, nil)
	_ = n.Send(context.Background(), notify.Event{SecretPath: "secret/x"})

	if received.Timestamp.IsZero() {
		t.Error("expected timestamp to be set automatically")
	}
}

func TestSend_Non2xxReturnsError(t *testing.T) {
	srv := newTestServer(t, http.StatusInternalServerError, nil)
	defer srv.Close()

	n, _ := notify.New(srv.URL, nil)
	err := n.Send(context.Background(), notify.Event{SecretPath: "secret/x"})
	if err == nil {
		t.Fatal("expected error for 500 response")
	}
}

func TestSend_ForwardsCustomHeaders(t *testing.T) {
	var gotHeader string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotHeader = r.Header.Get("X-Token")
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	n, _ := notify.New(srv.URL, map[string]string{"X-Token": "secret-token"})
	_ = n.Send(context.Background(), notify.Event{SecretPath: "secret/y"})

	if gotHeader != "secret-token" {
		t.Errorf("expected X-Token header, got %q", gotHeader)
	}
}
