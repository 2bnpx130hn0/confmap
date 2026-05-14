package resolver

import (
	"github.com/iamNilotpal/confmap/grouper"
	"github.com/iamNilotpal/confmap/loader"
	"github.com/iamNilotpal/confmap/merger"
	"github.com/iamNilotpal/confmap/validator"
)

// GroupedResolver resolves config and exposes it partitioned into prefix groups.
type GroupedResolver struct {
	*Resolver
	prefixes []string
	sep      string
}

// NewGrouped builds a Resolver that additionally exposes a ByPrefix view of
// the merged config. prefixes and sep are forwarded to grouper.Grouper.ByPrefix.
func NewGrouped(
	loaders []loader.Loader,
	schema map[string]validator.FieldSchema,
	prefixes []string,
	sep string,
) *GroupedResolver {
	return &GroupedResolver{
		Resolver: New(loaders, merger.New(), schema),
		prefixes: prefixes,
		sep:      sep,
	}
}

// Groups returns the merged config partitioned by the configured prefixes.
// It calls Resolve internally and returns an error if resolution fails.
func (gr *GroupedResolver) Groups() (map[string]map[string]any, error) {
	cfg, err := gr.Resolve()
	if err != nil {
		return nil, err
	}
	g := grouper.New(cfg)
	return g.ByPrefix(gr.prefixes, gr.sep), nil
}

// GroupByFunc resolves the config and partitions it using a custom function.
func (gr *GroupedResolver) GroupByFunc(fn func(string) string) (map[string]map[string]any, error) {
	cfg, err := gr.Resolve()
	if err != nil {
		return nil, err
	}
	g := grouper.New(cfg)
	return g.ByFunc(fn), nil
}
