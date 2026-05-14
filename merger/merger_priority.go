package merger

import "sort"

// PriorityLayer associates a config map with an integer priority.
// Higher priority values win during merge.
type PriorityLayer struct {
	Priority int
	Data     map[string]any
}

// PriorityMerger merges layers ordered by their declared priority.
type PriorityMerger struct {
	layers []PriorityLayer
}

// NewPriority creates a PriorityMerger with the given layers.
func NewPriority(layers []PriorityLayer) *PriorityMerger {
	return &PriorityMerger{layers: layers}
}

// AddLayer appends a new priority layer.
func (pm *PriorityMerger) AddLayer(priority int, data map[string]any) {
	pm.layers = append(pm.layers, PriorityLayer{Priority: priority, Data: data})
}

// Merge sorts layers by ascending priority and merges them so that
// higher-priority layers overwrite lower-priority ones.
func (pm *PriorityMerger) Merge() map[string]any {
	sorted := make([]PriorityLayer, len(pm.layers))
	copy(sorted, pm.layers)
	sort.SliceStable(sorted, func(i, j int) bool {
		return sorted[i].Priority < sorted[j].Priority
	})

	base := New()
	for _, layer := range sorted {
		if layer.Data == nil {
			continue
		}
		base.AddLayer(layer.Data)
	}
	return base.Merge()
}

// Layers returns a copy of the registered priority layers.
func (pm *PriorityMerger) Layers() []PriorityLayer {
	out := make([]PriorityLayer, len(pm.layers))
	copy(out, pm.layers)
	return out
}
