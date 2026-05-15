// Package matcher provides pattern-based key matching for config maps.
// It supports glob-style wildcards and exact matches against flat or nested keys.
package matcher

import (
	"path"
	"strings"
)

// Matcher holds a set of patterns and matches config keys against them.
type Matcher struct {
	patterns []string
}

// New creates a new Matcher with the given glob patterns.
// Patterns follow filepath.Match conventions (e.g. "db.*", "*.host").
func New(patterns ...string) *Matcher {
	return &Matcher{patterns: patterns}
}

// Match reports whether the given key matches any of the registered patterns.
func (m *Matcher) Match(key string) bool {
	for _, p := range m.patterns {
		if matched, err := path.Match(p, key); err == nil && matched {
			return true
		}
	}
	return false
}

// Filter returns a new config map containing only keys (flat dotted paths)
// that match at least one pattern.
func (m *Matcher) Filter(cfg map[string]any) map[string]any {
	result := make(map[string]any)
	for k, v := range flatten("", cfg) {
		if m.Match(k) {
			setNested(result, k, v)
		}
	}
	return result
}

// Keys returns all flat dotted keys from cfg that match at least one pattern.
func (m *Matcher) Keys(cfg map[string]any) []string {
	var out []string
	for k := range flatten("", cfg) {
		if m.Match(k) {
			out = append(out, k)
		}
	}
	return out
}

func flatten(prefix string, cfg map[string]any) map[string]any {
	out := make(map[string]any)
	for k, v := range cfg {
		full := k
		if prefix != "" {
			full = prefix + "." + k
		}
		if sub, ok := v.(map[string]any); ok {
			for fk, fv := range flatten(full, sub) {
				out[fk] = fv
			}
		} else {
			out[full] = v
		}
	}
	return out
}

func setNested(dst map[string]any, key string, value any) {
	parts := strings.SplitN(key, ".", 2)
	if len(parts) == 1 {
		dst[key] = value
		return
	}
	sub, ok := dst[parts[0]].(map[string]any)
	if !ok {
		sub = make(map[string]any)
		dst[parts[0]] = sub
	}
	setNested(sub, parts[1], value)
}
