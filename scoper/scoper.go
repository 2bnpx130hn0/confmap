// Package scoper provides scoped views into a configuration map,
// allowing consumers to work with a subtree of config keys under
// a given prefix without exposing the full configuration.
package scoper

import "fmt"

// Scoper provides a prefix-scoped view of a flat or nested config map.
type Scoper struct {
	prefix string
	data   map[string]any
}

// New creates a new Scoper rooted at the given prefix within data.
// If prefix is empty the Scoper operates on the top-level map.
func New(prefix string, data map[string]any) (*Scoper, error) {
	if prefix == "" {
		return &Scoper{prefix: "", data: data}, nil
	}
	sub, ok := data[prefix]
	if !ok {
		return nil, fmt.Errorf("scoper: prefix %q not found in config", prefix)
	}
	subMap, ok := sub.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("scoper: prefix %q is not a map", prefix)
	}
	return &Scoper{prefix: prefix, data: subMap}, nil
}

// Get returns the value for key within the scoped subtree.
func (s *Scoper) Get(key string) (any, bool) {
	v, ok := s.data[key]
	return v, ok
}

// Set writes a value for key within the scoped subtree.
// It does not modify the original map passed to New.
func (s *Scoper) Set(key string, value any) {
	s.data[key] = value
}

// Keys returns all keys present in the scoped subtree.
func (s *Scoper) Keys() []string {
	keys := make([]string, 0, len(s.data))
	for k := range s.data {
		keys = append(keys, k)
	}
	return keys
}

// Prefix returns the prefix this Scoper is rooted at.
func (s *Scoper) Prefix() string {
	return s.prefix
}

// Snapshot returns a shallow copy of the scoped subtree.
func (s *Scoper) Snapshot() map[string]any {
	copy := make(map[string]any, len(s.data))
	for k, v := range s.data {
		copy[k] = v
	}
	return copy
}
