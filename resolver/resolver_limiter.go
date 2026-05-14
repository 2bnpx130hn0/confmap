package resolver

import (
	"fmt"

	"github.com/iamthe1whoknocks/confmap/limiter"
	"github.com/iamthe1whoknocks/confmap/loader"
	"github.com/iamthe1whoknocks/confmap/merger"
	"github.com/iamthe1whoknocks/confmap/validator"
)

// LimitedResolver wraps a Resolver and enforces structural limits via a Limiter.
type LimitedResolver struct {
	inner   *Resolver
	limiter *limiter.Limiter
}

// NewLimited creates a Resolver that enforces key-count and depth limits
// after merging and validating all config layers.
func NewLimited(loaders []loader.Loader, schema map[string]any, maxKeys, maxDepth int) *LimitedResolver {
	return &LimitedResolver{
		inner:   New(loaders, merger.New(), validator.New(schema)),
		limiter: limiter.New(maxKeys, maxDepth),
	}
}

// Resolve merges all layers, validates against the schema, and enforces limits.
func (lr *LimitedResolver) Resolve() (map[string]any, error) {
	cfg, err := lr.inner.Resolve()
	if err != nil {
		return nil, err
	}
	if err := lr.limiter.Check(cfg); err != nil {
		return nil, fmt.Errorf("resolver: limit check failed: %w", err)
	}
	return cfg, nil
}
