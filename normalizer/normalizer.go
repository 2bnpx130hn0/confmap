// Package normalizer provides key normalization for config maps,
// converting keys to a canonical form (e.g. lowercase, trimmed, separator-unified).
package normalizer

import (
	"strings"
)

// Options controls how keys are normalized.
type Options struct {
	// Lowercase converts all keys to lowercase.
	Lowercase bool
	// TrimSpace removes leading/trailing whitespace from keys.
	TrimSpace bool
	// Separator replaces any occurrence of the OldSep with NewSep in keys.
	OldSep string
	NewSep string
}

// Normalizer applies key normalization to a config map.
type Normalizer struct {
	opts Options
}

// New returns a Normalizer with the given Options.
func New(opts Options) *Normalizer {
	return &Normalizer{opts: opts}
}

// Apply returns a new map with all keys (at every nesting level) normalized
// according to the Normalizer's Options. The original map is not mutated.
func (n *Normalizer) Apply(cfg map[string]any) map[string]any {
	return n.normalizeMap(cfg)
}

func (n *Normalizer) normalizeMap(m map[string]any) map[string]any {
	out := make(map[string]any, len(m))
	for k, v := range m {
		nk := n.normalizeKey(k)
		out[nk] = n.normalizeValue(v)
	}
	return out
}

func (n *Normalizer) normalizeKey(k string) string {
	if n.opts.TrimSpace {
		k = strings.TrimSpace(k)
	}
	if n.opts.Lowercase {
		k = strings.ToLower(k)
	}
	if n.opts.OldSep != "" && n.opts.NewSep != n.opts.OldSep {
		k = strings.ReplaceAll(k, n.opts.OldSep, n.opts.NewSep)
	}
	return k
}

func (n *Normalizer) normalizeValue(v any) any {
	switch val := v.(type) {
	case map[string]any:
		return n.normalizeMap(val)
	case []any:
		out := make([]any, len(val))
		for i, item := range val {
			out[i] = n.normalizeValue(item)
		}
		return out
	default:
		return v
	}
}

// Keys returns the normalized top-level keys of the given config map.
func (n *Normalizer) Keys(cfg map[string]any) []string {
	normalized := n.Apply(cfg)
	keys := make([]string, 0, len(normalized))
	for k := range normalized {
		keys = append(keys, k)
	}
	return keys
}
