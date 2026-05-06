// Package patcher provides functionality to apply partial updates
// (patches) to a configuration map using dot-notation key paths.
package patcher

import (
	"fmt"
	"strings"
)

// Patcher applies key-path patches to a configuration map.
type Patcher struct {
	config map[string]any
}

// New creates a new Patcher wrapping the given config map.
func New(config map[string]any) *Patcher {
	return &Patcher{config: config}
}

// Apply sets the value at the given dot-notation key path.
// Intermediate maps are created as needed.
// Returns an error if a non-map value blocks traversal.
func (p *Patcher) Apply(keyPath string, value any) error {
	keys := strings.Split(keyPath, ".")
	current := p.config

	for i, key := range keys[:len(keys)-1] {
		val, exists := current[key]
		if !exists {
			next := make(map[string]any)
			current[key] = next
			current = next
			continue
		}
		next, ok := val.(map[string]any)
		if !ok {
			return fmt.Errorf("patcher: key %q at segment %d is not a map", strings.Join(keys[:i+1], "."), i)
		}
		current = next
	}

	current[keys[len(keys)-1]] = value
	return nil
}

// Delete removes the value at the given dot-notation key path.
// Returns an error if the path cannot be traversed.
func (p *Patcher) Delete(keyPath string) error {
	keys := strings.Split(keyPath, ".")
	current := p.config

	for i, key := range keys[:len(keys)-1] {
		val, exists := current[key]
		if !exists {
			return fmt.Errorf("patcher: key %q not found", strings.Join(keys[:i+1], "."))
		}
		next, ok := val.(map[string]any)
		if !ok {
			return fmt.Errorf("patcher: key %q at segment %d is not a map", strings.Join(keys[:i+1], "."), i)
		}
		current = next
	}

	delete(current, keys[len(keys)-1])
	return nil
}

// Config returns the underlying (mutated) configuration map.
func (p *Patcher) Config() map[string]any {
	return p.config
}
