package fanout_test

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"

	"github.com/your-org/vaultpipe/internal/fanout"
)

// mockDest is a test destination that records calls and optionally returns an error.
type mockDest struct {
	name    string
	err     error
	called  atomic.Int32
	received map[string]string
}

func (m *mockDest) Write(_ context.Context, secrets map[string]string) error {
	m.called.Add(1)
	m.received = secrets
	return m.err
}

func (m *mockDest) Name() string { return m.name }

func TestNew_NoDest_ReturnsError(t *testing.T) {
	_, err := fanout.New()
	if err == nil {
		t.Fatal("expected error for zero destinations")
	}
}

func TestNew_WithDests_OK(t *testing.T) {
	d := &mockDest{name: "a"}
	_, err := fanout.New(d)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBroadcast_CallsAllDests(t *testing.T) {
	d1 := &mockDest{name: "d1"}
	d2 := &mockDest{name: "d2"}

	f, _ := fanout.New(d1, d2)
	secrets := map[string]string{"KEY": "val"}

	results := f.Broadcast(context.Background(), secrets)

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if d1.called.Load() != 1 {
		t.Error("d1 was not called exactly once")
	}
	if d2.called.Load() != 1 {
		t.Error("d2 was not called exactly once")
	}
}

func TestBroadcast_ReturnsPartialErrors(t *testing.T) {
	sentinel := errors.New("write failed")
	d1 := &mockDest{name: "ok"}
	d2 := &mockDest{name: "bad", err: sentinel}

	f, _ := fanout.New(d1, d2)
	results := f.Broadcast(context.Background(), map[string]string{})

	if err := fanout.Errors(results); err == nil {
		t.Fatal("expected combined error, got nil")
	}
}

func TestErrors_AllSuccess_ReturnsNil(t *testing.T) {
	d1 := &mockDest{name: "a"}
	d2 := &mockDest{name: "b"}

	f, _ := fanout.New(d1, d2)
	results := f.Broadcast(context.Background(), map[string]string{"X": "1"})

	if err := fanout.Errors(results); err != nil {
		t.Fatalf("expected nil, got: %v", err)
	}
}

func TestBroadcast_SecretsDeliveredToAllDests(t *testing.T) {
	d1 := &mockDest{name: "a"}
	d2 := &mockDest{name: "b"}
	secrets := map[string]string{"DB_PASS": "secret123"}

	f, _ := fanout.New(d1, d2)
	f.Broadcast(context.Background(), secrets)

	for _, d := range []*mockDest{d1, d2} {
		if d.received["DB_PASS"] != "secret123" {
			t.Errorf("%s: expected secret not delivered", d.name)
		}
	}
}
