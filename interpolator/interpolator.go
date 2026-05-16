// Package interpolator provides environment variable and value interpolation
// for config maps, expanding ${VAR} and $VAR style references.
package interpolator

import (
	"fmt"
	"os"
	"strings"
)

// Interpolator expands variable references in config string values.
type Interpolator struct {
	vars   map[string]string
	strict bool
}

// New creates an Interpolator using the provided variable map.
// If strict is true, missing variables return an error instead of an empty string.
func New(vars map[string]string, strict bool) *Interpolator {
	if vars == nil {
		vars = map[string]string{}
	}
	return &Interpolator{vars: vars, strict: strict}
}

// NewFromEnv creates an Interpolator seeded from the current process environment.
func NewFromEnv(strict bool) *Interpolator {
	vars := map[string]string{}
	for _, entry := range os.Environ() {
		parts := strings.SplitN(entry, "=", 2)
		if len(parts) == 2 {
			vars[parts[0]] = parts[1]
		}
	}
	return New(vars, strict)
}

// Apply walks the config map and expands variable references in all string values.
func (i *Interpolator) Apply(cfg map[string]any) (map[string]any, error) {
	result := make(map[string]any, len(cfg))
	for k, v := range cfg {
		expanded, err := i.expandValue(v)
		if err != nil {
			return nil, fmt.Errorf("interpolator: key %q: %w", k, err)
		}
		result[k] = expanded
	}
	return result, nil
}

func (i *Interpolator) expandValue(v any) (any, error) {
	switch val := v.(type) {
	case string:
		return i.expand(val)
	case map[string]any:
		return i.Apply(val)
	case []any:
		out := make([]any, len(val))
		for idx, item := range val {
			expanded, err := i.expandValue(item)
			if err != nil {
				return nil, err
			}
			out[idx] = expanded
		}
		return out, nil
	default:
		return v, nil
	}
}

func (i *Interpolator) expand(s string) (string, error) {
	var err error
	result := os.Expand(s, func(key string) string {
		if val, ok := i.vars[key]; ok {
			return val
		}
		if i.strict && err == nil {
			err = fmt.Errorf("undefined variable %q", key)
		}
		return ""
	})
	if err != nil {
		return "", err
	}
	return result, nil
}
