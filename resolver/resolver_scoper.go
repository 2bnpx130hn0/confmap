package resolver

import (
	"fmt"

	"github.com/user/confmap/scoper"
)

// ScopedResolver wraps a Resolver and exposes a scoped view of the
// resolved configuration under a fixed prefix.
type ScopedResolver struct {
	base   *Resolver
	prefix string
}

// NewScoped creates a ScopedResolver that, after resolving, returns
// only the subtree rooted at prefix.
func NewScoped(r *Resolver, prefix string) *ScopedResolver {
	return &ScopedResolver{base: r, prefix: prefix}
}

// Resolve runs the underlying resolver and returns a Scoper rooted
// at the configured prefix.
func (sr *ScopedResolver) Resolve() (*scoper.Scoper, error) {
	cfg, err := sr.base.Resolve()
	if err != nil {
		return nil, fmt.Errorf("scoped resolver: %w", err)
	}
	sc, err := scoper.New(sr.prefix, cfg)
	if err != nil {
		return nil, fmt.Errorf("scoped resolver: %w", err)
	}
	return sc, nil
}

// Prefix returns the prefix this ScopedResolver scopes to.
func (sr *ScopedResolver) Prefix() string {
	return sr.prefix
}
