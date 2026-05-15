package resolver

import (
	"fmt"

	"github.com/nicholasgasior/confmap/aggregator"
	"github.com/nicholasgasior/confmap/loader"
	"github.com/nicholasgasior/confmap/validator"
)

// NewAggregated builds a Resolver that loads each source via its loader,
// then reduces all resulting maps using the provided aggregator reduce function
// before running schema validation.
//
// This is useful when callers need custom accumulation semantics beyond simple
// last-write-wins merging.
func NewAggregated(
	loaders []loader.Loader,
	schema map[string]any,
	reduceFn func(acc, layer map[string]any) (map[string]any, error),
) *Resolver {
	return New(loaderFunc(func() (map[string]any, error) {
		agg := aggregator.New(reduceFn)
		for i, l := range loaders {
			cfg, err := l.Load()
			if err != nil {
				return nil, fmt.Errorf("resolver/aggregated: loader %d: %w", i, err)
			}
			agg.Add(cfg)
		}
		return agg.Reduce()
	}), validator.Schema(schema))
}

// loaderFunc is an adapter that lets a plain function satisfy loader.Loader.
type loaderFunc func() (map[string]any, error)

func (f loaderFunc) Load() (map[string]any, error) { return f() }
