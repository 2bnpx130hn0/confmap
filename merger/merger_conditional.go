package merger

import "fmt"

// Condition is a predicate evaluated against the accumulated config map.
type Condition func(cfg map[string]any) bool

// ConditionalLayer pairs a config layer with a condition that gates its inclusion.
type ConditionalLayer struct {
	Layer     map[string]any
	Condition Condition
}

// ConditionalMerger merges only the layers whose conditions are satisfied.
type ConditionalMerger struct {
	base   *Merger
	layers []ConditionalLayer
}

// NewConditional returns a ConditionalMerger built on top of the standard Merger.
func NewConditional(layers []ConditionalLayer) *ConditionalMerger {
	return &ConditionalMerger{
		base:   New(),
		layers: layers,
	}
}

// Merge evaluates each layer's condition against the accumulated config and
// merges layers that pass into the result.
func (cm *ConditionalMerger) Merge() map[string]any {
	result := map[string]any{}
	for _, cl := range cm.layers {
		if cl.Condition != nil && !cl.Condition(result) {
			continue
		}
		result = cm.base.Merge(result, cl.Layer)
	}
	return result
}

// Always is a convenience Condition that always returns true.
func Always() Condition {
	return func(_ map[string]any) bool { return true }
}

// Never is a convenience Condition that always returns false.
func Never() Condition {
	return func(_ map[string]any) bool { return false }
}

// HasKey returns a Condition satisfied when the key exists in the accumulated config.
func HasKey(key string) Condition {
	return func(cfg map[string]any) bool {
		_, ok := cfg[key]
		return ok
	}
}

// ValueEquals returns a Condition satisfied when cfg[key] string-equals value.
func ValueEquals(key string, value any) Condition {
	return func(cfg map[string]any) bool {
		v, ok := cfg[key]
		if !ok {
			return false
		}
		return fmt.Sprintf("%v", v) == fmt.Sprintf("%v", value)
	}
}
