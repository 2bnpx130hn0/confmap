package resolver

import (
	"fmt"

	"github.com/iamBijoyKar/confmap/aliaser"
	"github.com/iamBijoyKar/confmap/merger"
	"github.com/iamBijoyKar/confmap/validator"
)

// AliasedResolver resolves config and expands registered key aliases
// before returning the final merged config map.
type AliasedResolver struct {
	base    *Resolver
	aliaser *aliaser.Aliaser
}

// NewAliased creates a Resolver that applies key aliasing after merging.
// Aliases are expanded so callers always see canonical key names.
func NewAliased(
	loaders []Loader,
	schema map[string]any,
	a *aliaser.Aliaser,
) (*AliasedResolver, error) {
	if a == nil {
		return nil, fmt.Errorf("resolver: aliaser must not be nil")
	}
	m := merger.New()
	v := validator.New(schema)
	return &AliasedResolver{
		base:    New(loaders, m, v),
		aliaser: a,
	}, nil
}

// Resolve merges all loader layers, expands aliases, and validates the result.
func (ar *AliasedResolver) Resolve() (map[string]any, error) {
	cfg, err := ar.base.Resolve()
	if err != nil {
		return nil, err
	}
	return ar.aliaser.Apply(cfg)
}
