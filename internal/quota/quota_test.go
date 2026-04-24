package quota_test

import (
	"testing"
	"time"

	"github.com/yourusername/vaultpipe/internal/quota"
)

func TestNew_InvalidLimit(t *testing.T) {
	_, err := quota.New(0, time.Minute)
	if err == nil {
		t.Fatal("expected error for limit=0, got nil")
	}
}

func TestNew_InvalidWindow(t *testing.T) {
	_, err := quota.New(5, 0)
	if err == nil {
		t.Fatal("expected error for window=0, got nil")
	}
}

func TestNew_Valid(t *testing.T) {
	tr, err := quota.New(3, time.Minute)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tr == nil {
		t.Fatal("expected non-nil tracker")
	}
}

func TestCheck_WithinLimit(t *testing.T) {
	tr, _ := quota.New(3, time.Minute)
	for i := 0; i < 3; i++ {
		if err := tr.Check("secret/app"); err != nil {
			t.Fatalf("unexpected error on read %d: %v", i+1, err)
		}
	}
}

func TestCheck_ExceedsLimit(t *testing.T) {
	tr, _ := quota.New(2, time.Minute)
	_ = tr.Check("secret/db")
	_ = tr.Check("secret/db")
	err := tr.Check("secret/db")
	if err == nil {
		t.Fatal("expected quota exceeded error, got nil")
	}
	var qe *quota.ErrQuotaExceeded
	if ok := errorAs(err, &qe); !ok {
		t.Fatalf("expected *ErrQuotaExceeded, got %T", err)
	}
	if qe.Key != "secret/db" {
		t.Errorf("expected key %q, got %q", "secret/db", qe.Key)
	}
	if qe.Limit != 2 {
		t.Errorf("expected limit 2, got %d", qe.Limit)
	}
}

func TestCheck_WindowExpiry(t *testing.T) {
	tr, _ := quota.New(1, 50*time.Millisecond)
	if err := tr.Check("key"); err != nil {
		t.Fatalf("first check failed: %v", err)
	}
	if err := tr.Check("key"); err == nil {
		t.Fatal("expected quota exceeded before window reset")
	}
	time.Sleep(60 * time.Millisecond)
	if err := tr.Check("key"); err != nil {
		t.Fatalf("expected check to pass after window reset: %v", err)
	}
}

func TestReset_ClearsCounter(t *testing.T) {
	tr, _ := quota.New(1, time.Minute)
	_ = tr.Check("k")
	tr.Reset("k")
	if err := tr.Check("k"); err != nil {
		t.Fatalf("expected check to pass after reset: %v", err)
	}
}

func TestRemaining_ReturnsCorrectCount(t *testing.T) {
	tr, _ := quota.New(5, time.Minute)
	if r := tr.Remaining("x"); r != 5 {
		t.Fatalf("expected 5 remaining for unseen key, got %d", r)
	}
	_ = tr.Check("x")
	_ = tr.Check("x")
	if r := tr.Remaining("x"); r != 3 {
		t.Fatalf("expected 3 remaining, got %d", r)
	}
}

func TestFlush_ClearsAll(t *testing.T) {
	tr, _ := quota.New(1, time.Minute)
	_ = tr.Check("a")
	_ = tr.Check("b")
	tr.Flush()
	for _, k := range []string{"a", "b"} {
		if err := tr.Check(k); err != nil {
			t.Fatalf("expected check to pass after flush for key %q: %v", k, err)
		}
	}
}

// errorAs is a helper to avoid importing errors in test file directly.
func errorAs(err error, target **quota.ErrQuotaExceeded) bool {
	if qe, ok := err.(*quota.ErrQuotaExceeded); ok {
		*target = qe
		return true
	}
	return false
}
