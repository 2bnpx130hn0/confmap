// Package filter provides utilities for selecting, excluding, and
// projecting keys from a config map.
package filter

import "sort"

// Filter holds the include/exclude rules applied to a config map.
type Filter struct {
	includes map[string]struct{}
	excludes map[string]struct{}
}

// New returns a new Filter with no rules applied.
func New() *Filter {
	return &Filter{
		includes: make(map[string]struct{}),
		excludes: make(map[string]struct{}),
	}
}

// Include registers top-level keys that should be kept.
// If at least one Include rule is registered, only those keys survive.
func (f *Filter) Include(keys ...string) *Filter {
	for _, k := range keys {
		f.includes[k] = struct{}{}
	}
	return f
}

// Exclude registers top-level keys that should be removed.
func (f *Filter) Exclude(keys ...string) *Filter {
	for _, k := range keys {
		f.excludes[k] = struct{}{}
	}
	return f
}

// Apply returns a new map with the filter rules applied.
// Include rules take precedence: if any includes are defined only those
// keys survive (minus any that are also explicitly excluded).
func (f *Filter) Apply(cfg map[string]any) map[string]any {
	out := make(map[string]any, len(cfg))
	for k, v := range cfg {
		if len(f.includes) > 0 {
			if _, ok := f.includes[k]; !ok {
				continue
			}
		}
		if _, ok := f.excludes[k]; ok {
			continue
		}
		out[k] = v
	}
	return out
}

// Keys returns a sorted slice of top-level keys present in cfg.
func Keys(cfg map[string]any) []string {
	keys := make([]string, 0, len(cfg))
	for k := range cfg {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
