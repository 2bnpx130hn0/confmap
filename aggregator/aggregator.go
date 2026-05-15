// Package aggregator provides utilities for collecting and combining
// config maps by applying reduce-style operations across multiple sources.
package aggregator

import "fmt"

// Aggregator collects config layers and reduces them into a single map.
type Aggregator struct {
	layers []map[string]any
	reduceFn func(acc, layer map[string]any) (map[string]any, error)
}

// New returns an Aggregator using the provided reduce function.
// The reduce function is called for each layer in order, accumulating results.
func New(reduceFn func(acc, layer map[string]any) (map[string]any, error)) *Aggregator {
	return &Aggregator{reduceFn: reduceFn}
}

// Add appends a config layer to the aggregator. Nil layers are silently skipped.
func (a *Aggregator) Add(layer map[string]any) *Aggregator {
	if layer != nil {
		a.layers = append(a.layers, layer)
	}
	return a
}

// Reduce applies the reduce function across all added layers and returns the
// accumulated result. Returns an empty map if no layers have been added.
func (a *Aggregator) Reduce() (map[string]any, error) {
	if a.reduceFn == nil {
		return nil, fmt.Errorf("aggregator: reduce function must not be nil")
	}
	acc := map[string]any{}
	for i, layer := range a.layers {
		var err error
		acc, err = a.reduceFn(acc, layer)
		if err != nil {
			return nil, fmt.Errorf("aggregator: reduce error at layer %d: %w", i, err)
		}
	}
	return acc, nil
}

// Count returns the number of layers currently held by the aggregator.
func (a *Aggregator) Count() int {
	return len(a.layers)
}

// MergeReduce is a built-in reduce function that merges layers with later
// values overriding earlier ones for matching keys.
func MergeReduce(acc, layer map[string]any) (map[string]any, error) {
	result := make(map[string]any, len(acc))
	for k, v := range acc {
		result[k] = v
	}
	for k, v := range layer {
		result[k] = v
	}
	return result, nil
}

// CollectKeys is a built-in reduce function that accumulates all unique keys
// seen across layers, storing true as the value for each key.
func CollectKeys(acc, layer map[string]any) (map[string]any, error) {
	result := make(map[string]any, len(acc))
	for k, v := range acc {
		result[k] = v
	}
	for k := range layer {
		result[k] = true
	}
	return result, nil
}
