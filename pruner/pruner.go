// Package pruner provides utilities for removing nil, zero-value, or empty
// entries from a configuration map, producing a clean, compact representation.
package pruner

import "fmt"

// Option controls which values are considered prunable.
type Option int

const (
	// PruneNil removes keys whose value is nil.
	PruneNil Option = 1 << iota
	// PruneEmpty removes keys whose value is an empty string, slice, or map.
	PruneEmpty
	// PruneZero removes keys whose value is a numeric zero or false boolean.
	PruneZero
)

// Pruner removes unwanted entries from a config map.
type Pruner struct {
	opts Option
}

// New creates a Pruner with the given options OR-ed together.
// Pass 0 to use PruneNil only.
func New(opts Option) *Pruner {
	if opts == 0 {
		opts = PruneNil
	}
	return &Pruner{opts: opts}
}

// Apply returns a new map with prunable keys removed. The original is not
// mutated. Nested maps are pruned recursively.
func (p *Pruner) Apply(cfg map[string]any) (map[string]any, error) {
	if cfg == nil {
		return nil, fmt.Errorf("pruner: nil config")
	}
	return p.pruneMap(cfg), nil
}

// MustApply is like Apply but panics on error.
func (p *Pruner) MustApply(cfg map[string]any) map[string]any {
	out, err := p.Apply(cfg)
	if err != nil {
		panic(err)
	}
	return out
}

func (p *Pruner) pruneMap(m map[string]any) map[string]any {
	out := make(map[string]any, len(m))
	for k, v := range m {
		if nested, ok := v.(map[string]any); ok {
			pruned := p.pruneMap(nested)
			if len(pruned) > 0 || !p.has(PruneEmpty) {
				out[k] = pruned
			}
			continue
		}
		if p.shouldPrune(v) {
			continue
		}
		out[k] = v
	}
	return out
}

func (p *Pruner) shouldPrune(v any) bool {
	if v == nil && p.has(PruneNil) {
		return true
	}
	if p.has(PruneEmpty) {
		switch val := v.(type) {
		case string:
			if val == "" {
				return true
			}
		case []any:
			if len(val) == 0 {
				return true
			}
		}
	}
	if p.has(PruneZero) {
		switch val := v.(type) {
		case int:
			return val == 0
		case int64:
			return val == 0
		case float64:
			return val == 0
		case bool:
			return !val
		}
	}
	return false
}

func (p *Pruner) has(o Option) bool {
	return p.opts&o != 0
}
