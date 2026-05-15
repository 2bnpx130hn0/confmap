package merger

// Overlay merges a set of layers where each successive layer only overrides
// keys that are explicitly present (non-nil) in that layer. Unlike the base
// Merger, Overlay skips zero-value strings and nil values so that sparse
// "patch" layers do not accidentally blank out existing values.

// OverlayMerger applies layers in order, skipping blank/nil values.
type OverlayMerger struct {
	layers []map[string]any
}

// NewOverlay creates an OverlayMerger.
func NewOverlay() *OverlayMerger {
	return &OverlayMerger{}
}

// AddLayer appends a layer to the overlay stack.
func (o *OverlayMerger) AddLayer(layer map[string]any) *OverlayMerger {
	if layer != nil {
		o.layers = append(o.layers, layer)
	}
	return o
}

// Merge applies all layers in order, skipping nil and empty-string values.
func (o *OverlayMerger) Merge() map[string]any {
	result := map[string]any{}
	for _, layer := range o.layers {
		overlayCopy(result, layer)
	}
	return result
}

// overlayCopy recursively merges src into dst, skipping nil and "" values.
func overlayCopy(dst, src map[string]any) {
	for k, v := range src {
		if v == nil {
			continue
		}
		if s, ok := v.(string); ok && s == "" {
			continue
		}
		if srcMap, ok := v.(map[string]any); ok {
			if dstMap, ok := dst[k].(map[string]any); ok {
				merged := make(map[string]any, len(dstMap))
				for dk, dv := range dstMap {
					merged[dk] = dv
				}
				overlayCopy(merged, srcMap)
				dst[k] = merged
				continue
			}
		}
		dst[k] = v
	}
}
