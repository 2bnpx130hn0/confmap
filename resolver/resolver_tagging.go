package resolver

import "github.com/example/confmap/tagging"

// TaggedResolver wraps a Resolver and applies key-level tag annotations
// to the resolved config, exposing helpers for tag-aware access.
type TaggedResolver struct {
	resolver *Resolver
	tagger   *tagging.Tagger
}

// NewTagged creates a TaggedResolver from an existing Resolver and Tagger.
func NewTagged(r *Resolver, t *tagging.Tagger) *TaggedResolver {
	return &TaggedResolver{resolver: r, tagger: t}
}

// Resolve returns the merged, validated config map.
func (tr *TaggedResolver) Resolve() (map[string]any, error) {
	return tr.resolver.Resolve()
}

// ResolveAnnotated returns the resolved config with an additional
// "__tags__" key that maps each tagged key to its tag slice.
func (tr *TaggedResolver) ResolveAnnotated() (map[string]any, error) {
	cfg, err := tr.resolver.Resolve()
	if err != nil {
		return nil, err
	}
	return tr.tagger.Annotate(cfg), nil
}

// ResolveFiltered returns only the config keys that carry the given tag.
func (tr *TaggedResolver) ResolveFiltered(tag string) (map[string]any, error) {
	cfg, err := tr.resolver.Resolve()
	if err != nil {
		return nil, err
	}
	return tr.tagger.FilterByTag(cfg, tag), nil
}

// Tags returns the set of distinct tags present across all keys in the
// resolved config, which can be used to enumerate available filter values.
func (tr *TaggedResolver) Tags() ([]string, error) {
	cfg, err := tr.resolver.Resolve()
	if err != nil {
		return nil, err
	}
	return tr.tagger.ListTags(cfg), nil
}
