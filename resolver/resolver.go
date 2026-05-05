// Package resolver provides functionality to resolve a final merged
// configuration map from multiple layered sources with override precedence.
package resolver

import (
	"fmt"

	"github.com/yourorg/confmap/loader"
	"github.com/yourorg/confmap/merger"
	"github.com/yourorg/confmap/validator"
)

// Source represents a named configuration source with an associated loader.
type Source struct {
	Name   string
	Loader loader.Loader
}

// Resolver orchestrates loading, merging, and validating configuration layers.
type Resolver struct {
	sources []Source
	schema  map[string]interface{}
	m       *merger.Merger
}

// New creates a new Resolver with the given sources (in ascending priority order)
// and an optional JSON-Schema-like validation schema.
func New(schema map[string]interface{}, sources ...Source) *Resolver {
	return &Resolver{
		sources: sources,
		schema:  schema,
		m:       merger.New(),
	}
}

// Resolve loads all sources, merges them in order (later sources override earlier
// ones), validates the result against the schema, and returns the final config map.
func (r *Resolver) Resolve() (map[string]interface{}, error) {
	layers := make([]map[string]interface{}, 0, len(r.sources))

	for _, src := range r.sources {
		data, err := src.Loader.Load()
		if err != nil {
			return nil, fmt.Errorf("resolver: loading source %q: %w", src.Name, err)
		}
		layers = append(layers, data)
	}

	merged, err := r.m.Merge(layers...)
	if err != nil {
		return nil, fmt.Errorf("resolver: merging layers: %w", err)
	}

	if r.schema != nil {
		if err := validator.Validate(merged, r.schema); err != nil {
			return nil, fmt.Errorf("resolver: validation failed: %w", err)
		}
	}

	return merged, nil
}
