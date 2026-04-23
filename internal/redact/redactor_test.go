package redact_test

import (
	"testing"

	"github.com/yourusername/vaultpipe/internal/redact"
)

func TestNew_DefaultPlaceholder(t *testing.T) {
	r := redact.New("")
	if r.Placeholder() != "[REDACTED]" {
		t.Errorf("expected default placeholder, got %q", r.Placeholder())
	}
}

func TestNew_CustomPlaceholder(t *testing.T) {
	r := redact.New("***")
	if r.Placeholder() != "***" {
		t.Errorf("expected '***', got %q", r.Placeholder())
	}
}

func TestRedact_ReplacesSecretValue(t *testing.T) {
	r := redact.New("")
	r.Load(map[string]string{"DB_PASS": "supersecret"})

	got := r.Redact("connecting with password supersecret to db")
	want := "connecting with password [REDACTED] to db"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestRedact_MultipleValues(t *testing.T) {
	r := redact.New("[REDACTED]")
	r.Load(map[string]string{
		"TOKEN": "tok-abc123",
		"KEY":   "mykey",
	})

	got := r.Redact("token=tok-abc123 key=mykey")
	if got == "token=tok-abc123 key=mykey" {
		t.Error("expected values to be redacted")
	}
	if got != "token=[REDACTED] key=[REDACTED]" {
		t.Errorf("unexpected result: %q", got)
	}
}

func TestRedact_EmptyValueIgnored(t *testing.T) {
	r := redact.New("")
	r.Load(map[string]string{"EMPTY": ""})

	got := r.Redact("nothing should change here")
	want := "nothing should change here"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestRedact_NoSecretsLoaded(t *testing.T) {
	r := redact.New("")
	got := r.Redact("plain text")
	if got != "plain text" {
		t.Errorf("got %q, want 'plain text'", got)
	}
}

func TestRedactMap_ReplacesAllValues(t *testing.T) {
	r := redact.New("[REDACTED]")
	secrets := map[string]string{
		"API_KEY": "abc",
		"DB_PASS": "xyz",
	}

	got := r.RedactMap(secrets)
	for k, v := range got {
		if v != "[REDACTED]" {
			t.Errorf("key %q: expected [REDACTED], got %q", k, v)
		}
	}
	if len(got) != len(secrets) {
		t.Errorf("expected %d keys, got %d", len(secrets), len(got))
	}
}
