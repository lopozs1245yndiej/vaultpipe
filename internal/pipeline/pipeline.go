// Package pipeline provides a composable secret processing pipeline
// that chains multiple transformation stages together before writing output.
package pipeline

import (
	"context"
	"fmt"
)

// Stage is a function that processes a map of secrets and returns a new map.
type Stage func(ctx context.Context, secrets map[string]string) (map[string]string, error)

// Pipeline executes a sequence of stages over a set of secrets.
type Pipeline struct {
	stages []Stage
}

// New creates a new Pipeline with the given stages.
func New(stages ...Stage) *Pipeline {
	return &Pipeline{stages: stages}
}

// Add appends one or more stages to the pipeline.
func (p *Pipeline) Add(stages ...Stage) *Pipeline {
	p.stages = append(p.stages, stages...)
	return p
}

// Run executes all stages in order, passing the output of each stage as the
// input to the next. It returns the final transformed secrets or the first
// error encountered.
func (p *Pipeline) Run(ctx context.Context, secrets map[string]string) (map[string]string, error) {
	current := copyMap(secrets)
	for i, stage := range p.stages {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("pipeline cancelled at stage %d: %w", i, ctx.Err())
		default:
		}
		var err error
		current, err = stage(ctx, current)
		if err != nil {
			return nil, fmt.Errorf("pipeline stage %d: %w", i, err)
		}
	}
	return current, nil
}

// Len returns the number of stages in the pipeline.
func (p *Pipeline) Len() int {
	return len(p.stages)
}

func copyMap(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}
