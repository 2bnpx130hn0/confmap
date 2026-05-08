// Package cloner provides deep-copy utilities for config maps.
// It ensures that cloned configs are fully independent from their
// source, preventing accidental mutation across layers.
package cloner

import "fmt"

// Cloner performs deep copies of config maps.
type Cloner struct{}

// New returns a new Cloner instance.
func New() *Cloner {
	return &Cloner{}
}

// Clone returns a deep copy of the provided config map.
// All nested maps and slices are recursively duplicated.
func (c *Cloner) Clone(config map[string]any) (map[string]any, error) {
	result, err := deepCopyMap(config)
	if err != nil {
		return nil, fmt.Errorf("cloner: %w", err)
	}
	return result, nil
}

// MustClone clones the config and panics on error.
func (c *Cloner) MustClone(config map[string]any) map[string]any {
	result, err := c.Clone(config)
	if err != nil {
		panic(err)
	}
	return result
}

func deepCopyMap(m map[string]any) (map[string]any, error) {
	out := make(map[string]any, len(m))
	for k, v := range m {
		copied, err := deepCopyValue(v)
		if err != nil {
			return nil, fmt.Errorf("key %q: %w", k, err)
		}
		out[k] = copied
	}
	return out, nil
}

func deepCopyValue(v any) (any, error) {
	switch val := v.(type) {
	case map[string]any:
		return deepCopyMap(val)
	case []any:
		return deepCopySlice(val)
	case string, int, int64, float64, bool, nil:
		return val, nil
	default:
		return nil, fmt.Errorf("unsupported value type %T", v)
	}
}

func deepCopySlice(s []any) ([]any, error) {
	out := make([]any, len(s))
	for i, v := range s {
		copied, err := deepCopyValue(v)
		if err != nil {
			return nil, fmt.Errorf("index %d: %w", i, err)
		}
		out[i] = copied
	}
	return out, nil
}
