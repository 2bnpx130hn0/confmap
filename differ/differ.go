// Package differ computes the diff between two config maps,
// reporting added, removed, and changed keys.
package differ

import "fmt"

// ChangeType describes the kind of change detected.
type ChangeType string

const (
	Added   ChangeType = "added"
	Removed ChangeType = "removed"
	Changed ChangeType = "changed"
)

// Delta represents a single key-level change between two configs.
type Delta struct {
	Key      string
	Type     ChangeType
	OldValue interface{}
	NewValue interface{}
}

// Differ computes deltas between config snapshots.
type Differ struct{}

// New returns a new Differ.
func New() *Differ {
	return &Differ{}
}

// Diff returns the list of deltas between the old and new config maps.
// Both maps are flattened with dot-notation keys before comparison.
func (d *Differ) Diff(oldCfg, newCfg map[string]interface{}) []Delta {
	oldFlat := flatten("", oldCfg)
	newFlat := flatten("", newCfg)

	var deltas []Delta

	for k, ov := range oldFlat {
		nv, exists := newFlat[k]
		if !exists {
			deltas = append(deltas, Delta{Key: k, Type: Removed, OldValue: ov})
		} else if fmt.Sprintf("%v", ov) != fmt.Sprintf("%v", nv) {
			deltas = append(deltas, Delta{Key: k, Type: Changed, OldValue: ov, NewValue: nv})
		}
	}

	for k, nv := range newFlat {
		if _, exists := oldFlat[k]; !exists {
			deltas = append(deltas, Delta{Key: k, Type: Added, NewValue: nv})
		}
	}

	return deltas
}

// flatten recursively flattens a nested map into dot-notation keys.
func flatten(prefix string, m map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range m {
		fullKey := k
		if prefix != "" {
			fullKey = prefix + "." + k
		}
		if nested, ok := v.(map[string]interface{}); ok {
			for fk, fv := range flatten(fullKey, nested) {
				result[fk] = fv
			}
		} else {
			result[fullKey] = v
		}
	}
	return result
}
