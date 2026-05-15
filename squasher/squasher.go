// Package squasher collapses a slice of config maps into a single flat map
// by applying each layer in order, with later layers taking precedence.
// Unlike the merger package, squasher produces a fully flattened dot-notation
// map suitable for serialisation or environment variable export.
package squasher

import "fmt"

// Squasher collapses layered config maps into a single flattened map.
type Squasher struct {
	separator string
}

// New returns a Squasher that uses sep as the key separator (e.g. ".").
func New(sep string) *Squasher {
	if sep == "" {
		sep = "."
	}
	return &Squasher{separator: sep}
}

// Squash merges all layers left-to-right (later layers win) and returns a
// single-level map with dot-separated keys.
func (s *Squasher) Squash(layers ...map[string]any) (map[string]any, error) {
	merged := make(map[string]any)
	for _, layer := range layers {
		if layer == nil {
			continue
		}
		flat, err := s.flatten("", layer)
		if err != nil {
			return nil, err
		}
		for k, v := range flat {
			merged[k] = v
		}
	}
	return merged, nil
}

// flatten recursively walks m and builds dot-separated keys.
func (s *Squasher) flatten(prefix string, m map[string]any) (map[string]any, error) {
	out := make(map[string]any)
	for k, v := range m {
		if k == "" {
			return nil, fmt.Errorf("squasher: empty key encountered under prefix %q", prefix)
		}
		fullKey := k
		if prefix != "" {
			fullKey = prefix + s.separator + k
		}
		switch child := v.(type) {
		case map[string]any:
			sub, err := s.flatten(fullKey, child)
			if err != nil {
				return nil, err
			}
			for sk, sv := range sub {
				out[sk] = sv
			}
		default:
			out[fullKey] = v
		}
	}
	return out, nil
}

// Expand converts a flat dot-separated map back into a nested map.
// It returns an error if a key segment collision is detected.
func (s *Squasher) Expand(flat map[string]any) (map[string]any, error) {
	root := make(map[string]any)
	for key, val := range flat {
		if err := insertNested(root, splitKey(key, s.separator), val); err != nil {
			return nil, fmt.Errorf("squasher: expand %q: %w", key, err)
		}
	}
	return root, nil
}

func splitKey(key, sep string) []string {
	var parts []string
	start := 0
	for i := 0; i <= len(key)-len(sep); i++ {
		if key[i:i+len(sep)] == sep {
			parts = append(parts, key[start:i])
			start = i + len(sep)
			i += len(sep) - 1
		}
	}
	parts = append(parts, key[start:])
	return parts
}

func insertNested(m map[string]any, parts []string, val any) error {
	if len(parts) == 1 {
		m[parts[0]] = val
		return nil
	}
	child, exists := m[parts[0]]
	if !exists {
		next := make(map[string]any)
		m[parts[0]] = next
		return insertNested(next, parts[1:], val)
	}
	next, ok := child.(map[string]any)
	if !ok {
		return fmt.Errorf("key segment %q already holds a scalar value", parts[0])
	}
	return insertNested(next, parts[1:], val)
}
