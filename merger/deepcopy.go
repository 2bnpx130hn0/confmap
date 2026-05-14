package merger

// DeepCopy returns a deep copy of the given config map.
// Nested maps (map[string]any) and slices are recursively copied.
// All other values are copied by assignment (assumed to be value types
// or immutable references such as strings and numbers).
func DeepCopy(src map[string]any) map[string]any {
	if src == nil {
		return nil
	}
	dst := make(map[string]any, len(src))
	for k, v := range src {
		dst[k] = deepCopyValue(v)
	}
	return dst
}

func deepCopyValue(v any) any {
	switch val := v.(type) {
	case map[string]any:
		return DeepCopy(val)
	case []any:
		copy := make([]any, len(val))
		for i, item := range val {
			copy[i] = deepCopyValue(item)
		}
		return copy
	default:
		return v
	}
}
