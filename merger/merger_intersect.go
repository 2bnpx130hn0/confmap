package merger

// IntersectMerger merges only the keys that are present in ALL provided layers.
// Keys that do not appear in every layer are omitted from the result.
// When a key is present in all layers the value from the last layer wins.
type IntersectMerger struct {
	layers []map[string]any
}

// NewIntersect returns a new IntersectMerger.
func NewIntersect() *IntersectMerger {
	return &IntersectMerger{}
}

// AddLayer appends a configuration layer to the merger.
func (m *IntersectMerger) AddLayer(layer map[string]any) {
	if layer == nil {
		return
	}
	m.layers = append(m.layers, layer)
}

// Merge returns a map containing only the keys shared by every layer.
// Values are taken from the last layer that contains the key.
func (m *IntersectMerger) Merge() map[string]any {
	if len(m.layers) == 0 {
		return map[string]any{}
	}

	// Seed candidate keys from the first layer.
	candidate := make(map[string]struct{}, len(m.layers[0]))
	for k := range m.layers[0] {
		candidate[k] = struct{}{}
	}

	// Retain only keys present in every subsequent layer.
	for _, layer := range m.layers[1:] {
		for k := range candidate {
			if _, ok := layer[k]; !ok {
				delete(candidate, k)
			}
		}
	}

	// Build result: last-layer value wins for each surviving key.
	result := make(map[string]any, len(candidate))
	for k := range candidate {
		for _, layer := range m.layers {
			if v, ok := layer[k]; ok {
				result[k] = v
			}
		}
	}
	return result
}
