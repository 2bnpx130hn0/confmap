package merger

import "fmt"

// MergeStrategy defines how values are merged when keys conflict.
type MergeStrategy int

const (
	// StrategyOverride replaces existing values with new ones (default).
	StrategyOverride MergeStrategy = iota
	// StrategyKeepExisting keeps the first value encountered.
	StrategyKeepExisting
)

// Merger merges multiple config maps with a defined precedence.
type Merger struct {
	strategy MergeStrategy
}

// New creates a new Merger with the given strategy.
func New(strategy MergeStrategy) *Merger {
	return &Merger{strategy: strategy}
}

// Merge combines layers of config maps in order.
// Later layers have higher precedence when StrategyOverride is used.
func (m *Merger) Merge(layers ...map[string]any) (map[string]any, error) {
	result := make(map[string]any)
	for i, layer := range layers {
		if layer == nil {
			return nil, fmt.Errorf("merger: layer %d is nil", i)
		}
		if err := m.mergeInto(result, layer); err != nil {
			return nil, fmt.Errorf("merger: layer %d: %w", i, err)
		}
	}
	return result, nil
}

func (m *Merger) mergeInto(dst, src map[string]any) error {
	for k, srcVal := range src {
		dstVal, exists := dst[k]
		if !exists {
			dst[k] = srcVal
			continue
		}
		srcMap, srcIsMap := srcVal.(map[string]any)
		dstMap, dstIsMap := dstVal.(map[string]any)
		if srcIsMap && dstIsMap {
			if err := m.mergeInto(dstMap, srcMap); err != nil {
				return err
			}
			continue
		}
		switch m.strategy {
		case StrategyOverride:
			dst[k] = srcVal
		case StrategyKeepExisting:
			// keep dstVal, do nothing
		}
	}
	return nil
}
