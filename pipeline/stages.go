package pipeline

import (
	"fmt"
	"strings"
)

// StageFunc wraps a plain function as a Stage with a descriptive label (for debugging).
func StageFunc(label string, fn func(map[string]any) (map[string]any, error)) Stage {
	return func(cfg map[string]any) (map[string]any, error) {
		result, err := fn(cfg)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", label, err)
		}
		return result, nil
	}
}

// FilterKeys returns a Stage that removes keys whose names match any of the
// provided prefixes (case-insensitive).
func FilterKeys(prefixes ...string) Stage {
	return StageFunc("filter-keys", func(cfg map[string]any) (map[string]any, error) {
		out := make(map[string]any, len(cfg))
		for k, v := range cfg {
			matched := false
			for _, p := range prefixes {
				if strings.HasPrefix(strings.ToLower(k), strings.ToLower(p)) {
					matched = true
					break
				}
			}
			if !matched {
				out[k] = v
			}
		}
		return out, nil
	})
}

// SetDefaults returns a Stage that applies default values for missing keys.
func SetDefaults(defaults map[string]any) Stage {
	return StageFunc("set-defaults", func(cfg map[string]any) (map[string]any, error) {
		out := make(map[string]any, len(cfg))
		for k, v := range cfg {
			out[k] = v
		}
		for k, v := range defaults {
			if _, exists := out[k]; !exists {
				out[k] = v
			}
		}
		return out, nil
	})
}

// RequireKeys returns a Stage that fails if any of the specified keys are absent.
func RequireKeys(keys ...string) Stage {
	return StageFunc("require-keys", func(cfg map[string]any) (map[string]any, error) {
		for _, k := range keys {
			if _, ok := cfg[k]; !ok {
				return nil, fmt.Errorf("required key %q is missing", k)
			}
		}
		return cfg, nil
	})
}
