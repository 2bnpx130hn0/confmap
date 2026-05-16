package resolver

import (
	"fmt"

	"github.com/confmap/interpolator"
	"github.com/confmap/loader"
	"github.com/confmap/merger"
	"github.com/confmap/validator"
)

// InterpolatedResolver resolves config and expands variable references in all
// string values using the provided Interpolator before validation.
type InterpolatedResolver struct {
	loaders      []loader.Loader
	merger       *merger.Merger
	schema       map[string]any
	interpolator *interpolator.Interpolator
}

// NewInterpolated returns a Resolver that applies variable interpolation after
// merging layers and before schema validation.
func NewInterpolated(
	loaders []loader.Loader,
	schema map[string]any,
	interp *interpolator.Interpolator,
) *InterpolatedResolver {
	return &InterpolatedResolver{
		loaders:      loaders,
		merger:       merger.New(),
		schema:       schema,
		interpolator: interp,
	}
}

// Resolve loads, merges, interpolates, and validates the configuration.
func (r *InterpolatedResolver) Resolve() (map[string]any, error) {
	var layers []map[string]any
	for _, l := range r.loaders {
		cfg, err := l.Load()
		if err != nil {
			return nil, fmt.Errorf("interpolated resolver: load: %w", err)
		}
		layers = append(layers, cfg)
	}

	merged := r.merger.Merge(layers...)

	if r.interpolator != nil {
		var err error
		merged, err = r.interpolator.Apply(merged)
		if err != nil {
			return nil, fmt.Errorf("interpolated resolver: interpolate: %w", err)
		}
	}

	if r.schema != nil {
		if err := validator.Validate(merged, r.schema); err != nil {
			return nil, fmt.Errorf("interpolated resolver: validate: %w", err)
		}
	}

	return merged, nil
}
