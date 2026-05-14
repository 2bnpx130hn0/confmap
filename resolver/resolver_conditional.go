package resolver

import (
	"fmt"

	"github.com/your-org/confmap/merger"
	"github.com/your-org/confmap/validator"
)

// ConditionalResolver resolves a config by conditionally merging layers and
// then validating the result against an optional schema.
type ConditionalResolver struct {
	cm     *merger.ConditionalMerger
	schema map[string]any
}

// NewConditional builds a ConditionalResolver from a set of conditional layers
// and an optional validation schema (may be nil).
func NewConditional(layers []merger.ConditionalLayer, schema map[string]any) *ConditionalResolver {
	return &ConditionalResolver{
		cm:     merger.NewConditional(layers),
		schema: schema,
	}
}

// Resolve merges the conditional layers and validates the result.
// Returns an error if validation fails.
func (cr *ConditionalResolver) Resolve() (map[string]any, error) {
	cfg := cr.cm.Merge()
	if cr.schema != nil {
		if err := validator.Validate(cfg, cr.schema); err != nil {
			return nil, fmt.Errorf("conditional resolver: validation failed: %w", err)
		}
	}
	return cfg, nil
}
