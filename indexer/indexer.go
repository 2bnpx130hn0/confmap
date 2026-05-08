// Package indexer provides flat-key indexing of nested config maps,
// enabling fast dot-notation lookups (e.g. "database.host").
package indexer

import "fmt"

// Indexer holds a flat key→value index built from a nested config map.
type Indexer struct {
	index map[string]any
}

// New creates an Indexer from a nested config map.
// All nested keys are flattened using dot notation.
func New(cfg map[string]any) *Indexer {
	idx := make(map[string]any)
	flatten("", cfg, idx)
	return &Indexer{index: idx}
}

// Get returns the value for a dot-notation key and whether it was found.
func (ix *Indexer) Get(key string) (any, bool) {
	v, ok := ix.index[key]
	return v, ok
}

// Keys returns all indexed keys in an unspecified order.
func (ix *Indexer) Keys() []string {
	keys := make([]string, 0, len(ix.index))
	for k := range ix.index {
		keys = append(keys, k)
	}
	return keys
}

// MustGet returns the value for key or panics if the key is absent.
func (ix *Indexer) MustGet(key string) any {
	v, ok := ix.index[key]
	if !ok {
		panic(fmt.Sprintf("indexer: key %q not found", key))
	}
	return v
}

// Rebuild replaces the internal index with a freshly flattened version
// of the supplied config map.
func (ix *Indexer) Rebuild(cfg map[string]any) {
	next := make(map[string]any)
	flatten("", cfg, next)
	ix.index = next
}

// Has reports whether the given dot-notation key exists in the index.
func (ix *Indexer) Has(key string) bool {
	_, ok := ix.index[key]
	return ok
}

// flatten recursively walks src, writing dot-separated keys into dst.
func flatten(prefix string, src map[string]any, dst map[string]any) {
	for k, v := range src {
		fullKey := k
		if prefix != "" {
			fullKey = prefix + "." + k
		}
		if nested, ok := v.(map[string]any); ok {
			flatten(fullKey, nested, dst)
		} else {
			dst[fullKey] = v
		}
	}
}
