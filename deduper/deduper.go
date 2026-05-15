// Package deduper provides utilities for removing duplicate keys and values
// from configuration maps, with support for nested structures.
package deduper

import "fmt"

// Deduper removes duplicate entries from a configuration map.
type Deduper struct {
	seenKeys map[string]struct{}
}

// New creates a new Deduper instance.
func New() *Deduper {
	return &Deduper{}
}

// DedupeKeys returns a new config map with only the first occurrence of each
// flattened key path retained. Subsequent duplicate keys are dropped.
func (d *Deduper) DedupeKeys(config map[string]any) (map[string]any, error) {
	if config == nil {
		return nil, fmt.Errorf("deduper: config must not be nil")
	}
	d.seenKeys = make(map[string]struct{})
	return d.dedupeMap(config, ""), nil
}

func (d *Deduper) dedupeMap(m map[string]any, prefix string) map[string]any {
	result := make(map[string]any, len(m))
	for k, v := range m {
		fullKey := k
		if prefix != "" {
			fullKey = prefix + "." + k
		}
		if _, seen := d.seenKeys[fullKey]; seen {
			continue
		}
		d.seenKeys[fullKey] = struct{}{}
		if nested, ok := v.(map[string]any); ok {
			result[k] = d.dedupeMap(nested, fullKey)
		} else {
			result[k] = v
		}
	}
	return result
}

// DedupeValues returns a new config map where keys whose values are duplicated
// across siblings at the same level are collapsed, keeping only the last seen value.
func (d *Deduper) DedupeValues(config map[string]any) (map[string]any, error) {
	if config == nil {
		return nil, fmt.Errorf("deduper: config must not be nil")
	}
	return dedupeValuesMap(config), nil
}

func dedupeValuesMap(m map[string]any) map[string]any {
	seen := make(map[string]string) // value string -> first key
	result := make(map[string]any, len(m))
	for k, v := range m {
		if nested, ok := v.(map[string]any); ok {
			result[k] = dedupeValuesMap(nested)
			continue
		}
		valStr := fmt.Sprintf("%v", v)
		if firstKey, exists := seen[valStr]; exists {
			delete(result, firstKey)
		}
		seen[valStr] = k
		result[k] = v
	}
	return result
}

// Count returns the number of unique flattened keys in the config.
func (d *Deduper) Count(config map[string]any) int {
	if config == nil {
		return 0
	}
	return countKeys(config)
}

func countKeys(m map[string]any) int {
	count := 0
	for _, v := range m {
		if nested, ok := v.(map[string]any); ok {
			count += countKeys(nested)
		} else {
			count++
		}
	}
	return count
}
