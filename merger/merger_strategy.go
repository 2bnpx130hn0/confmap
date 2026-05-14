package merger

// Strategy defines how two values are merged when a key exists in both layers.
type Strategy int

const (
	// StrategyOverride replaces the base value with the overlay value (default).
	StrategyOverride Strategy = iota
	// StrategyKeepBase retains the base value when the key already exists.
	StrategyKeepBase
	// StrategyAppendSlice appends overlay slice elements to the base slice;
	// falls back to StrategyOverride for non-slice values.
	StrategyAppendSlice
)

// StrategicMerger merges config layers using a configurable Strategy.
type StrategicMerger struct {
	strategy Strategy
}

// NewStrategic returns a StrategicMerger that applies the given Strategy.
func NewStrategic(s Strategy) *StrategicMerger {
	return &StrategicMerger{strategy: s}
}

// Merge combines base and overlay maps according to the configured strategy.
// It never mutates base; it returns a new map.
func (sm *StrategicMerger) Merge(base, overlay map[string]any) map[string]any {
	result := make(map[string]any, len(base))
	for k, v := range base {
		result[k] = v
	}

	for k, overlayVal := range overlay {
		baseVal, exists := result[k]

		switch sm.strategy {
		case StrategyKeepBase:
			if !exists {
				result[k] = overlayVal
			}

		case StrategyAppendSlice:
			if exists {
				baseSlice, baseOK := baseVal.([]any)
				overlaySlice, overlayOK := overlayVal.([]any)
				if baseOK && overlayOK {
					result[k] = append(baseSlice, overlaySlice...)
					continue
				}
			}
			result[k] = overlayVal

		default: // StrategyOverride
			result[k] = overlayVal
		}
	}

	return result
}
