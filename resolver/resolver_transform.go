package resolver

import (
	"fmt"

	"github.com/yourorg/confmap/transformer"
)

// WithTransformer returns a new Resolver that applies t to the merged
// config map before schema validation.
func (r *Resolver) WithTransformer(t *transformer.Transformer) *Resolver {
	return &Resolver{
		loaders:     r.loaders,
		merger:      r.merger,
		validator:   r.validator,
		transformer: t,
	}
}

// resolveWithTransform is the internal helper used by Resolve when a
// Transformer has been attached. It merges all layers, applies
// transformations, then validates.
func (r *Resolver) resolveWithTransform() (map[string]any, error) {
	merged, err := r.merge()
	if err != nil {
		return nil, fmt.Errorf("resolve: merge: %w", err)
	}

	if r.transformer != nil {
		if err := r.transformer.Apply(merged); err != nil {
			return nil, fmt.Errorf("resolve: transform: %w", err)
		}
	}

	if r.validator != nil {
		if err := r.validator.Validate(merged); err != nil {
			return nil, fmt.Errorf("resolve: validate: %w", err)
		}
	}

	return merged, nil
}
