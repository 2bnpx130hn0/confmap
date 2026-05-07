package resolver

import (
	"fmt"

	"github.com/igorvisi/confmap/linker"
)

// LinkedResolver wraps a Resolver and resolves cross-key references
// in the merged config using the linker package.
type LinkedResolver struct {
	base *Resolver
}

// NewLinked creates a LinkedResolver that applies reference resolution
// after merging and optional validation.
func NewLinked(base *Resolver) *LinkedResolver {
	return &LinkedResolver{base: base}
}

// Resolve merges all layers, validates if a schema is set, then resolves
// all ${key} references in the resulting config map.
func (lr *LinkedResolver) Resolve() (map[string]any, error) {
	merged, err := lr.base.Resolve()
	if err != nil {
		return nil, fmt.Errorf("linked resolver: %w", err)
	}

	l := linker.New(merged)
	resolved, err := l.Resolve()
	if err != nil {
		return nil, fmt.Errorf("linked resolver: reference resolution failed: %w", err)
	}

	return resolved, nil
}
