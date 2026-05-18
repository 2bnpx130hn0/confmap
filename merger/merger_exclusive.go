package merger

// ExclusiveMerger merges layers but only includes keys that appear in exactly
// one layer (i.e. keys shared across multiple layers are excluded from the result).
type ExclusiveMerger struct {
	layers []map[string]any
}

// NewExclusive creates an ExclusiveMerger.
func NewExclusive() *ExclusiveMerger {
	return &ExclusiveMerger{}
}

// AddLayer appends a config layer to the merger.
func (e *ExclusiveMerger) AddLayer(layer map[string]any) {
	if layer != nil {
		e.layers = append(e.layers, layer)
	}
}

// Merge returns a map containing only the keys that appear in exactly one layer.
// If a key appears in two or more layers it is omitted entirely.
func (e *ExclusiveMerger) Merge() map[string]any {
	if len(e.layers) == 0 {
		return map[string]any{}
	}

	// count how many layers contain each key
	count := map[string]int{}
	values := map[string]any{}

	for _, layer := range e.layers {
		for k, v := range layer {
			count[k]++
			values[k] = v // last-seen value kept for single-occurrence keys
		}
	}

	result := map[string]any{}
	for k, cnt := range count {
		if cnt == 1 {
			result[k] = values[k]
		}
	}
	return result
}
