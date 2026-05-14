// Package grouper provides functionality to group config keys by a common prefix
// or a custom grouping function, returning sub-maps for each group.
package grouper

import "strings"

// Grouper partitions a flat or nested config map into named groups.
type Grouper struct {
	data map[string]any
}

// New creates a new Grouper wrapping the provided config map.
func New(data map[string]any) *Grouper {
	return &Grouper{data: data}
}

// ByPrefix groups top-level keys that share the given prefix (separated by sep).
// Keys matching "<prefix><sep>*" are collected under the prefix group with the
// prefix+sep stripped from their names. Keys not matching any prefix are placed
// under the empty-string group "".
func (g *Grouper) ByPrefix(prefixes []string, sep string) map[string]map[string]any {
	result := make(map[string]map[string]any)

	for key, val := range g.data {
		matched := false
		for _, p := range prefixes {
			token := p + sep
			if strings.HasPrefix(key, token) {
				if result[p] == nil {
					result[p] = make(map[string]any)
				}
				result[p][strings.TrimPrefix(key, token)] = val
				matched = true
				break
			}
		}
		if !matched {
			if result[""] == nil {
				result[""] = make(map[string]any)
			}
			result[""][key] = val
		}
	}
	return result
}

// ByFunc groups all top-level keys using a caller-supplied function that maps
// a key to a group name. Keys for which fn returns "" are placed in the
// ungrouped bucket under the empty-string key.
func (g *Grouper) ByFunc(fn func(key string) string) map[string]map[string]any {
	result := make(map[string]map[string]any)
	for key, val := range g.data {
		group := fn(key)
		if result[group] == nil {
			result[group] = make(map[string]any)
		}
		result[group][key] = val
	}
	return result
}

// Groups returns the distinct group names that would result from ByPrefix.
func (g *Grouper) Groups(prefixes []string, sep string) []string {
	seen := make(map[string]struct{})
	for _, p := range prefixes {
		token := p + sep
		for key := range g.data {
			if strings.HasPrefix(key, token) {
				seen[p] = struct{}{}
			}
		}
	}
	out := make([]string, 0, len(seen))
	for k := range seen {
		out = append(out, k)
	}
	return out
}
