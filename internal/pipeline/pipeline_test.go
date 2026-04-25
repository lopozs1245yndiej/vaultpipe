package pipeline_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/yourusername/vaultpipe/internal/pipeline"
)

func TestNew_EmptyPipeline(t *testing.T) {
	p := pipeline.New()
	if p.Len() != 0 {
		t.Fatalf("expected 0 stages, got %d", p.Len())
	}
}

func TestRun_PassesThroughWithNoStages(t *testing.T) {
	p := pipeline.New()
	input := map[string]string{"KEY": "value"}
	out, err := p.Run(context.Background(), input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["KEY"] != "value" {
		t.Errorf("expected value 'value', got %q", out["KEY"])
	}
}

func TestRun_AppliesStagesInOrder(t *testing.T) {
	uppercaseStage := func(_ context.Context, secrets map[string]string) (map[string]string, error) {
		out := make(map[string]string, len(secrets))
		for k, v := range secrets {
			out[k] = strings.ToUpper(v)
		}
		return out, nil
	}
	prefixStage := func(_ context.Context, secrets map[string]string) (map[string]string, error) {
		out := make(map[string]string, len(secrets))
		for k, v := range secrets {
			out[k] = "PREFIX_" + v
		}
		return out, nil
	}

	p := pipeline.New(uppercaseStage, prefixStage)
	out, err := p.Run(context.Background(), map[string]string{"k": "hello"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["k"] != "PREFIX_HELLO" {
		t.Errorf("expected 'PREFIX_HELLO', got %q", out["k"])
	}
}

func TestRun_StopsOnStageError(t *testing.T) {
	sentinel := errors.New("stage failed")
	failStage := func(_ context.Context, _ map[string]string) (map[string]string, error) {
		return nil, sentinel
	}
	neverStage := func(_ context.Context, _ map[string]string) (map[string]string, error) {
		t.Error("neverStage should not have been called")
		return nil, nil
	}

	p := pipeline.New(failStage, neverStage)
	_, err := p.Run(context.Background(), map[string]string{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, sentinel) {
		t.Errorf("expected sentinel error, got %v", err)
	}
}

func TestRun_CancelledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	called := false
	stage := func(_ context.Context, s map[string]string) (map[string]string, error) {
		called = true
		return s, nil
	}
	p := pipeline.New(stage)
	_, err := p.Run(ctx, map[string]string{"k": "v"})
	if err == nil {
		t.Fatal("expected cancellation error")
	}
	if called {
		t.Error("stage should not have been called after context cancellation")
	}
}

func TestAdd_AppendsStages(t *testing.T) {
	stage := func(_ context.Context, s map[string]string) (map[string]string, error) { return s, nil }
	p := pipeline.New(stage)
	p.Add(stage, stage)
	if p.Len() != 3 {
		t.Errorf("expected 3 stages, got %d", p.Len())
	}
}
