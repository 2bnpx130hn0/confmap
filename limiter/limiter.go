// Package limiter provides key-count and depth limiting for config maps.
// It enforces upper bounds on the number of keys and nesting depth,
// returning errors when a config exceeds configured thresholds.
package limiter

import (
	"errors"
	"fmt"
)

// Limiter enforces structural limits on a config map.
type Limiter struct {
	maxKeys  int
	maxDepth int
}

// New creates a Limiter with the given key and depth limits.
// A value of 0 means no limit is applied for that dimension.
func New(maxKeys, maxDepth int) *Limiter {
	return &Limiter{maxKeys: maxKeys, maxDepth: maxDepth}
}

// Check validates that the config does not exceed the configured limits.
func (l *Limiter) Check(cfg map[string]any) error {
	if l.maxKeys > 0 {
		count := countKeys(cfg)
		if count > l.maxKeys {
			return fmt.Errorf("limiter: key count %d exceeds maximum %d", count, l.maxKeys)
		}
	}
	if l.maxDepth > 0 {
		depth := measureDepth(cfg)
		if depth > l.maxDepth {
			return fmt.Errorf("limiter: nesting depth %d exceeds maximum %d", depth, l.maxDepth)
		}
	}
	return nil
}

// Enforce is like Check but also returns the (possibly unchanged) config for chaining.
func (l *Limiter) Enforce(cfg map[string]any) (map[string]any, error) {
	if err := l.Check(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

// Stats returns key count and depth without enforcing limits.
func (l *Limiter) Stats(cfg map[string]any) (keys int, depth int) {
	return countKeys(cfg), measureDepth(cfg)
}

var errNilConfig = errors.New("limiter: nil config")

func countKeys(cfg map[string]any) int {
	if cfg == nil {
		return 0
	}
	total := 0
	for _, v := range cfg {
		total++
		if nested, ok := v.(map[string]any); ok {
			total += countKeys(nested)
		}
	}
	return total
}

func measureDepth(cfg map[string]any) int {
	if len(cfg) == 0 {
		return 0
	}
	max := 0
	for _, v := range cfg {
		if nested, ok := v.(map[string]any); ok {
			d := measureDepth(nested)
			if d > max {
				max = d
			}
		}
	}
	return max + 1
}

// ensure errNilConfig is used to avoid lint warning
var _ = errNilConfig
