// Package freezer provides the ability to lock a config map to prevent
// further modifications, producing an immutable snapshot.
package freezer

import (
	"errors"
	"fmt"
)

// ErrFrozen is returned when a mutation is attempted on a frozen config.
var ErrFrozen = errors.New("config is frozen and cannot be modified")

// Freezer wraps a config map and enforces immutability once frozen.
type Freezer struct {
	data   map[string]any
	frozen bool
}

// New creates a new Freezer with a deep copy of the provided config.
func New(cfg map[string]any) *Freezer {
	return &Freezer{data: deepCopy(cfg)}
}

// Freeze locks the config, preventing any further writes.
func (f *Freezer) Freeze() {
	f.frozen = true
}

// IsFrozen reports whether the config has been frozen.
func (f *Freezer) IsFrozen() bool {
	return f.frozen
}

// Get retrieves a value by key. Works whether frozen or not.
func (f *Freezer) Get(key string) (any, bool) {
	v, ok := f.data[key]
	return v, ok
}

// Set sets a key/value pair. Returns ErrFrozen if the config is locked.
func (f *Freezer) Set(key string, value any) error {
	if f.frozen {
		return fmt.Errorf("%w: attempted to set key %q", ErrFrozen, key)
	}
	f.data[key] = value
	return nil
}

// Snapshot returns a deep copy of the current config map.
func (f *Freezer) Snapshot() map[string]any {
	return deepCopy(f.data)
}

// deepCopy recursively copies a map[string]any.
func deepCopy(src map[string]any) map[string]any {
	if src == nil {
		return nil
	}
	dst := make(map[string]any, len(src))
	for k, v := range src {
		if nested, ok := v.(map[string]any); ok {
			dst[k] = deepCopy(nested)
		} else {
			dst[k] = v
		}
	}
	return dst
}
