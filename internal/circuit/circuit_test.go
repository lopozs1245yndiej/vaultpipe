package circuit_test

import (
	"testing"
	"time"

	"github.com/yourusername/vaultpipe/internal/circuit"
)

func TestNew_InvalidMaxFailures(t *testing.T) {
	_, err := circuit.New(0, time.Second)
	if err == nil {
		t.Fatal("expected error for zero maxFailures")
	}
}

func TestNew_InvalidResetTimeout(t *testing.T) {
	_, err := circuit.New(3, 0)
	if err == nil {
		t.Fatal("expected error for zero resetTimeout")
	}
}

func TestNew_Valid(t *testing.T) {
	b, err := circuit.New(3, time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b.State() != circuit.StateClosed {
		t.Errorf("expected StateClosed, got %v", b.State())
	}
}

func TestAllow_ClosedCircuit(t *testing.T) {
	b, _ := circuit.New(3, time.Second)
	if err := b.Allow(); err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}

func TestRecordFailure_OpensCircuitAtThreshold(t *testing.T) {
	b, _ := circuit.New(3, time.Second)
	b.RecordFailure()
	b.RecordFailure()
	if b.State() != circuit.StateClosed {
		t.Error("expected circuit to remain closed before threshold")
	}
	b.RecordFailure()
	if b.State() != circuit.StateOpen {
		t.Errorf("expected StateOpen after threshold, got %v", b.State())
	}
}

func TestAllow_OpenCircuit_ReturnsErrOpen(t *testing.T) {
	b, _ := circuit.New(1, 10*time.Second)
	b.RecordFailure()
	if err := b.Allow(); err != circuit.ErrOpen {
		t.Errorf("expected ErrOpen, got %v", err)
	}
}

func TestAllow_HalfOpenAfterTimeout(t *testing.T) {
	b, _ := circuit.New(1, 10*time.Millisecond)
	b.RecordFailure()
	time.Sleep(20 * time.Millisecond)
	if err := b.Allow(); err != nil {
		t.Errorf("expected nil after reset timeout, got %v", err)
	}
	if b.State() != circuit.StateHalfOpen {
		t.Errorf("expected StateHalfOpen, got %v", b.State())
	}
}

func TestRecordSuccess_ClosesCircuit(t *testing.T) {
	b, _ := circuit.New(1, 10*time.Millisecond)
	b.RecordFailure()
	time.Sleep(20 * time.Millisecond)
	_ = b.Allow() // transitions to half-open
	b.RecordSuccess()
	if b.State() != circuit.StateClosed {
		t.Errorf("expected StateClosed after success, got %v", b.State())
	}
	if b.Failures() != 0 {
		t.Errorf("expected 0 failures after success, got %d", b.Failures())
	}
}

func TestFailures_TracksCount(t *testing.T) {
	b, _ := circuit.New(5, time.Second)
	b.RecordFailure()
	b.RecordFailure()
	if b.Failures() != 2 {
		t.Errorf("expected 2 failures, got %d", b.Failures())
	}
}

func TestRecordFailure_HalfOpen_ReturnsToOpen(t *testing.T) {
	b, _ := circuit.New(1, 10*time.Millisecond)
	b.RecordFailure()
	time.Sleep(20 * time.Millisecond)
	_ = b.Allow() // transitions to half-open
	b.RecordFailure()
	if b.State() != circuit.StateOpen {
		t.Errorf("expected StateOpen after failure in half-open, got %v", b.State())
	}
}
