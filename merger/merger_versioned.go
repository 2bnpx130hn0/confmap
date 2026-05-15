package merger

import "fmt"

// VersionedLayer pairs a config map with a semantic version string.
type VersionedLayer struct {
	Version string
	Data    map[string]any
}

// VersionedMerger merges config layers in ascending version order,
// so higher versions always override lower ones regardless of insertion order.
type VersionedMerger struct {
	layers []VersionedLayer
}

// NewVersioned creates an empty VersionedMerger.
func NewVersioned() *VersionedMerger {
	return &VersionedMerger{}
}

// AddLayer registers a versioned config layer.
func (v *VersionedMerger) AddLayer(version string, data map[string]any) {
	v.layers = append(v.layers, VersionedLayer{Version: version, Data: data})
}

// Merge sorts layers by version (lexicographic) and merges them in ascending
// order so the highest version has final say on every key.
func (v *VersionedMerger) Merge() (map[string]any, error) {
	if len(v.layers) == 0 {
		return map[string]any{}, nil
	}

	sorted, err := sortedLayers(v.layers)
	if err != nil {
		return nil, err
	}

	result := map[string]any{}
	for _, layer := range sorted {
		if layer.Data == nil {
			continue
		}
		for k, val := range layer.Data {
			result[k] = val
		}
	}
	return result, nil
}

// sortedLayers returns layers ordered by ascending version string.
// Returns an error if any version string is empty.
func sortedLayers(layers []VersionedLayer) ([]VersionedLayer, error) {
	for _, l := range layers {
		if l.Version == "" {
			return nil, fmt.Errorf("merger: versioned layer has empty version string")
		}
	}
	out := make([]VersionedLayer, len(layers))
	copy(out, layers)
	// insertion sort — layer count is typically small
	for i := 1; i < len(out); i++ {
		for j := i; j > 0 && out[j].Version < out[j-1].Version; j-- {
			out[j], out[j-1] = out[j-1], out[j]
		}
	}
	return out, nil
}
