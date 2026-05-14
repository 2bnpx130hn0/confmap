package merger

import "sort"

// WeightedLayer pairs a config map with a numeric priority weight.
// Higher weight wins when two layers define the same key.
type WeightedLayer struct {
	Weight int
	Data   map[string]any
}

// WeightedMerger merges layers according to their declared weights rather
// than their slice order.  Layers with equal weight are merged in the order
// they were added (stable sort).
type WeightedMerger struct {
	layers []WeightedLayer
}

// NewWeighted returns an empty WeightedMerger.
func NewWeighted() *WeightedMerger {
	return &WeightedMerger{}
}

// Add appends a weighted layer to the merger.
func (w *WeightedMerger) Add(weight int, data map[string]any) *WeightedMerger {
	w.layers = append(w.layers, WeightedLayer{Weight: weight, Data: data})
	return w
}

// Merge returns a single map produced by applying all layers in ascending
// weight order so that higher-weight values overwrite lower-weight ones.
func (w *WeightedMerger) Merge() map[string]any {
	sorted := make([]WeightedLayer, len(w.layers))
	copy(sorted, w.layers)

	// Stable sort: ascending weight so the highest weight is applied last
	// (and therefore wins).
	sort.SliceStable(sorted, func(i, j int) bool {
		return sorted[i].Weight < sorted[j].Weight
	})

	result := map[string]any{}
	for _, layer := range sorted {
		if layer.Data == nil {
			continue
		}
		for k, v := range layer.Data {
			result[k] = v
		}
	}
	return result
}

// Layers returns the registered layers in the order they were added.
func (w *WeightedMerger) Layers() []WeightedLayer {
	out := make([]WeightedLayer, len(w.layers))
	copy(out, w.layers)
	return out
}
