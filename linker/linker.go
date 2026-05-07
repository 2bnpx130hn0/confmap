// Package linker provides cross-key reference resolution for config maps.
// It allows config values to reference other keys using a ${key.path} syntax,
// resolving them in dependency order.
package linker

import (
	"fmt"
	"regexp"
	"strings"
)

var refPattern = regexp.MustCompile(`\$\{([^}]+)\}`)

// Linker resolves internal references within a config map.
type Linker struct {
	config map[string]any
}

// New creates a new Linker for the given config.
func New(config map[string]any) *Linker {
	return &Linker{config: config}
}

// Resolve returns a new config map with all ${key} references substituted.
// Returns an error if a referenced key is missing or a cycle is detected.
func (l *Linker) Resolve() (map[string]any, error) {
	flat := flatten(l.config, "")
	resolved := make(map[string]string, len(flat))
	visiting := make(map[string]bool)

	var resolve func(key string) (string, error)
	resolve = func(key string) (string, error) {
		if v, ok := resolved[key]; ok {
			return v, nil
		}
		if visiting[key] {
			return "", fmt.Errorf("linker: cycle detected at key %q", key)
		}
		raw, ok := flat[key]
		if !ok {
			return "", fmt.Errorf("linker: key %q not found", key)
		}
		visiting[key] = true
		result := refPattern.ReplaceAllStringFunc(raw, func(match string) string {
			ref := match[2 : len(match)-1]
			v, err := resolve(ref)
			if err != nil {
				return match
			}
			return v
		})
		visiting[key] = false
		resolved[key] = result
		return result, nil
	}

	for key := range flat {
		if _, err := resolve(key); err != nil {
			return nil, err
		}
	}

	return unflatten(resolved), nil
}

func flatten(m map[string]any, prefix string) map[string]string {
	out := make(map[string]string)
	for k, v := range m {
		fullKey := k
		if prefix != "" {
			fullKey = prefix + "." + k
		}
		switch val := v.(type) {
		case map[string]any:
			for fk, fv := range flatten(val, fullKey) {
				out[fk] = fv
			}
		default:
			out[fullKey] = fmt.Sprintf("%v", val)
		}
	}
	return out
}

func unflatten(flat map[string]string) map[string]any {
	out := make(map[string]any)
	for k, v := range flat {
		parts := strings.SplitN(k, ".", 2)
		if len(parts) == 1 {
			out[k] = v
		} else {
			nested, ok := out[parts[0]]
			if !ok {
				nested = make(map[string]any)
				out[parts[0]] = nested
			}
			sub := unflatten(map[string]string{parts[1]: v})
			for sk, sv := range sub {
				nested.(map[string]any)[sk] = sv
			}
		}
	}
	return out
}
