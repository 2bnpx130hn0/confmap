package merger

// FallbackMerger merges a primary config and falls back to a default config
// for any keys that are missing or have zero/nil values in the primary.
//
// This is useful when you have a "base defaults" layer and a "user overrides"
// layer, and you want to guarantee all keys from the defaults appear in the
// final result even if the user did not specify them.
type FallbackMerger struct {
	primary  map[string]any
	fallback map[string]any
}

// NewFallback creates a FallbackMerger. primary values take precedence; any
// key absent from primary is filled from fallback.
func NewFallback(primary, fallback map[string]any) *FallbackMerger {
	return &FallbackMerger{
		primary:  primary,
		fallback: fallback,
	}
}

// Merge returns a new map that contains all keys from fallback, overridden by
// any keys present in primary. Neither input map is mutated.
func (f *FallbackMerger) Merge() map[string]any {
	result := make(map[string]any)

	// Seed with fallback values.
	for k, v := range f.fallback {
		result[k] = deepCopyValue(v)
	}

	// Override / extend with primary values.
	for k, v := range f.primary {
		if v == nil {
			continue
		}
		// If both sides are maps, recurse so nested defaults are preserved.
		if primaryMap, ok := v.(map[string]any); ok {
			if fallbackMap, fb := result[k].(map[string]any); fb {
				result[k] = NewFallback(primaryMap, fallbackMap).Merge()
				continue
			}
		}
		result[k] = deepCopyValue(v)
	}

	return result
}
