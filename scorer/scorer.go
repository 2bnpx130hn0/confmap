// Package scorer evaluates a config map and assigns a quality score
// based on key coverage, value completeness, and depth metrics.
package scorer

import "fmt"

// Score holds the result of evaluating a config.
type Score struct {
	Total     int
	MaxScore  int
	Coverage  float64 // 0.0 – 1.0
	Depth     int
	EmptyKeys int
}

// Scorer evaluates config maps.
type Scorer struct {
	expected []string
}

// New creates a Scorer that treats expectedKeys as the full set of
// keys a "complete" config should contain.
func New(expectedKeys []string) *Scorer {
	return &Scorer{expected: expectedKeys}
}

// Evaluate scores the provided config map.
func (s *Scorer) Evaluate(cfg map[string]any) (Score, error) {
	if cfg == nil {
		return Score{}, fmt.Errorf("scorer: config must not be nil")
	}

	flat := flatten(cfg, "")

	present := 0
	for _, k := range s.expected {
		if _, ok := flat[k]; ok {
			present++
		}
	}

	empty := 0
	for _, v := range flat {
		if v == nil || v == "" {
			empty++
		}
	}

	maxScore := len(s.expected)
	if maxScore == 0 {
		maxScore = 1
	}

	coverage := float64(present) / float64(maxScore)

	return Score{
		Total:     present,
		MaxScore:  maxScore,
		Coverage:  coverage,
		Depth:     maxDepth(cfg, 0),
		EmptyKeys: empty,
	}, nil
}

// flatten converts a nested map into dot-separated flat keys.
func flatten(m map[string]any, prefix string) map[string]any {
	out := make(map[string]any)
	for k, v := range m {
		full := k
		if prefix != "" {
			full = prefix + "." + k
		}
		if nested, ok := v.(map[string]any); ok {
			for nk, nv := range flatten(nested, full) {
				out[nk] = nv
			}
		} else {
			out[full] = v
		}
	}
	return out
}

// maxDepth returns the maximum nesting depth of the map.
func maxDepth(m map[string]any, current int) int {
	max := current
	for _, v := range m {
		if nested, ok := v.(map[string]any); ok {
			if d := maxDepth(nested, current+1); d > max {
				max = d
			}
		}
	}
	return max
}
