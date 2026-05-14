// Package pipeline provides a composable config processing pipeline
// that chains loaders, transformers, validators, and exporters together
// in a single fluent API.
package pipeline

import (
	"fmt"

	"github.com/example/confmap/merger"
)

// Stage represents a single processing step applied to a config map.
type Stage func(cfg map[string]any) (map[string]any, error)

// Pipeline holds an ordered list of stages to apply sequentially.
type Pipeline struct {
	stages []Stage
	name   string
}

// New creates a new named Pipeline.
func New(name string) *Pipeline {
	return &Pipeline{name: name}
}

// Use appends one or more stages to the pipeline.
func (p *Pipeline) Use(stages ...Stage) *Pipeline {
	p.stages = append(p.stages, stages...)
	return p
}

// Run executes all stages in order, passing the output of each stage
// as the input to the next. The initial config is not mutated.
func (p *Pipeline) Run(cfg map[string]any) (map[string]any, error) {
	current := merger.DeepCopy(cfg)
	for i, stage := range p.stages {
		result, err := stage(current)
		if err != nil {
			return nil, fmt.Errorf("pipeline %q stage %d: %w", p.name, i+1, err)
		}
		current = result
	}
	return current, nil
}

// Name returns the pipeline's identifier.
func (p *Pipeline) Name() string {
	return p.name
}

// Len returns the number of registered stages.
func (p *Pipeline) Len() int {
	return len(p.stages)
}
