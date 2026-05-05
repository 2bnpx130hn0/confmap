// Package transformer provides post-load transformation of config maps,
// allowing values to be coerced, renamed, or defaulted before validation.
package transformer

import "fmt"

// TransformFunc is a function that transforms a config map in place.
type TransformFunc func(data map[string]any) error

// Transformer applies a chain of TransformFuncs to a config map.
type Transformer struct {
	fns []TransformFunc
}

// New creates a new Transformer with the provided transform functions.
func New(fns ...TransformFunc) *Transformer {
	return &Transformer{fns: fns}
}

// Apply runs all registered transform functions against data in order.
// It returns the first error encountered, if any.
func (t *Transformer) Apply(data map[string]any) error {
	for i, fn := range t.fns {
		if err := fn(data); err != nil {
			return fmt.Errorf("transformer[%d]: %w", i, err)
		}
	}
	return nil
}

// SetDefault returns a TransformFunc that sets key to defaultVal
// if the key is absent from the config map.
func SetDefault(key string, defaultVal any) TransformFunc {
	return func(data map[string]any) error {
		if _, ok := data[key]; !ok {
			data[key] = defaultVal
		}
		return nil
	}
}

// Rename returns a TransformFunc that renames oldKey to newKey.
// If oldKey does not exist the operation is a no-op.
func Rename(oldKey, newKey string) TransformFunc {
	return func(data map[string]any) error {
		if v, ok := data[oldKey]; ok {
			data[newKey] = v
			delete(data, oldKey)
		}
		return nil
	}
}

// CoerceString returns a TransformFunc that converts the value at key
// to a string using fmt.Sprintf if the key exists and is not already a string.
func CoerceString(key string) TransformFunc {
	return func(data map[string]any) error {
		v, ok := data[key]
		if !ok {
			return nil
		}
		if _, already := v.(string); !already {
			data[key] = fmt.Sprintf("%v", v)
		}
		return nil
	}
}
