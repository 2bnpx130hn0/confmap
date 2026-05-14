package resolver

import (
	"fmt"

	"github.com/your-org/confmap/pruner"
)

// PrunedResolver wraps an existing Resolver and runs the merged config through
// a Pruner before returning it to callers.
type PrunedResolver struct {
	inner  *Resolver
	pruner *pruner.Pruner
}

// NewPruned creates a PrunedResolver that applies the given pruner options to
// the resolved configuration.
func NewPruned(r *Resolver, opts pruner.Option) *PrunedResolver {
	return &PrunedResolver{
		inner:  r,
		pruner: pruner.New(opts),
	}
}

// Resolve returns the merged, validated, and pruned configuration map.
func (pr *PrunedResolver) Resolve() (map[string]any, error) {
	cfg, err := pr.inner.Resolve()
	if err != nil {
		return nil, fmt.Errorf("pruned resolver: %w", err)
	}
	clean, err := pr.pruner.Apply(cfg)
	if err != nil {
		return nil, fmt.Errorf("pruned resolver: %w", err)
	}
	return clean, nil
}
