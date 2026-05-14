// Package flattener provides utilities for converting nested config maps
// into flat dot-notation key-value pairs and restoring them.
package flattener

import (
	"fmt"
	"strings"
)

// Flattener converts nested map[string]any configs to flat dot-separated
// key-value maps and back.
type Flattener struct {
	sep string
}

// New creates a Flattener using the given separator (e.g. ".").
func New(sep string) *Flattener {
	if sep == "" {
		sep = "."
	}
	return &Flattener{sep: sep}
}

// Flatten converts a nested map into a flat map with compound keys.
// Example: {"a": {"b": 1}} → {"a.b": 1}
func (f *Flattener) Flatten(input map[string]any) map[string]any {
	result := make(map[string]any)
	f.flattenInto(input, "", result)
	return result
}

func (f *Flattener) flattenInto(m map[string]any, prefix string, out map[string]any) {
	for k, v := range m {
		key := k
		if prefix != "" {
			key = prefix + f.sep + k
		}
		if nested, ok := v.(map[string]any); ok {
			f.flattenInto(nested, key, out)
		} else {
			out[key] = v
		}
	}
}

// Expand restores a flat dot-notation map into a nested map.
// Example: {"a.b": 1} → {"a": {"b": 1}}
func (f *Flattener) Expand(input map[string]any) (map[string]any, error) {
	result := make(map[string]any)
	for k, v := range input {
		if err := f.setNested(result, strings.Split(k, f.sep), v); err != nil {
			return nil, fmt.Errorf("flattener: expand key %q: %w", k, err)
		}
	}
	return result, nil
}

func (f *Flattener) setNested(m map[string]any, parts []string, value any) error {
	if len(parts) == 1 {
		m[parts[0]] = value
		return nil
	}
	next, ok := m[parts[0]]
	if !ok {
		child := make(map[string]any)
		m[parts[0]] = child
		return f.setNested(child, parts[1:], value)
	}
	child, ok := next.(map[string]any)
	if !ok {
		return fmt.Errorf("key %q is not a map", parts[0])
	}
	return f.setNested(child, parts[1:], value)
}
