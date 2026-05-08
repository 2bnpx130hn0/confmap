// Package sorter provides utilities for sorting config map keys
// and producing deterministic ordered representations of config data.
package sorter

import (
	"fmt"
	"sort"
	"strings"
)

// Sorter holds a config map and provides sorted access.
type Sorter struct {
	data map[string]any
}

// New creates a new Sorter from the given config map.
func New(data map[string]any) *Sorter {
	return &Sorter{data: data}
}

// Keys returns all top-level keys in sorted order.
func (s *Sorter) Keys() []string {
	keys := make([]string, 0, len(s.data))
	for k := range s.data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// FlatKeys returns all leaf key paths in sorted order using dot notation.
func (s *Sorter) FlatKeys() []string {
	var keys []string
	collectKeys(s.data, "", &keys)
	sort.Strings(keys)
	return keys
}

// Sorted returns a new map with the same data (shallow copy).
// Primarily useful for range iteration after calling Keys().
func (s *Sorter) Sorted() map[string]any {
	out := make(map[string]any, len(s.data))
	for k, v := range s.data {
		out[k] = v
	}
	return out
}

// SortedLines returns key=value lines in sorted key order for flat configs.
func (s *Sorter) SortedLines() []string {
	keys := s.FlatKeys()
	lines := make([]string, 0, len(keys))
	flat := flattenMap(s.data, "")
	for _, k := range keys {
		lines = append(lines, fmt.Sprintf("%s=%v", k, flat[k]))
	}
	return lines
}

func collectKeys(m map[string]any, prefix string, out *[]string) {
	for k, v := range m {
		full := k
		if prefix != "" {
			full = prefix + "." + k
		}
		if nested, ok := v.(map[string]any); ok {
			collectKeys(nested, full, out)
		} else {
			*out = append(*out, full)
		}
	}
}

func flattenMap(m map[string]any, prefix string) map[string]any {
	out := make(map[string]any)
	for k, v := range m {
		full := k
		if prefix != "" {
			full = prefix + "." + k
		}
		if nested, ok := v.(map[string]any); ok {
			for fk, fv := range flattenMap(nested, full) {
				out[fk] = fv
			}
		} else {
			out[full] = v
		}
	}
	return out
}

// IsSorted reports whether the given slice of strings is already sorted.
func IsSorted(keys []string) bool {
	return sort.StringsAreSorted(keys)
}

// FilterByPrefix returns flat keys that start with the given prefix.
func (s *Sorter) FilterByPrefix(prefix string) []string {
	var result []string
	for _, k := range s.FlatKeys() {
		if strings.HasPrefix(k, prefix) {
			result = append(result, k)
		}
	}
	return result
}
