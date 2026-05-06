// Package profiler provides utilities for inspecting and summarizing
// a merged config map, including key counts, depth analysis, and type profiles.
package profiler

import "fmt"

// Profile holds summary statistics about a config map.
type Profile struct {
	TotalKeys  int
	MaxDepth   int
	TypeCounts map[string]int
}

// Profiler inspects a config map and produces a Profile.
type Profiler struct{}

// New returns a new Profiler.
func New() *Profiler {
	return &Profiler{}
}

// Analyze walks the config map and returns a Profile.
func (p *Profiler) Analyze(cfg map[string]any) Profile {
	result := Profile{
		TypeCounts: make(map[string]int),
	}
	walkConfig(cfg, 1, &result)
	return result
}

func walkConfig(cfg map[string]any, depth int, prof *Profile) {
	for _, v := range cfg {
		prof.TotalKeys++
		if depth > prof.MaxDepth {
			prof.MaxDepth = depth
		}
		switch val := v.(type) {
		case map[string]any:
			prof.TypeCounts["map"]++
			walkConfig(val, depth+1, prof)
		case string:
			prof.TypeCounts["string"]++
		case int, int64, int32:
			prof.TypeCounts["int"]++
		case float64, float32:
			prof.TypeCounts["float"]++
		case bool:
			prof.TypeCounts["bool"]++
		case nil:
			prof.TypeCounts["nil"]++
		default:
			prof.TypeCounts[fmt.Sprintf("%T", v)]++
		}
	}
}

// Summary returns a human-readable summary string for a Profile.
func Summary(prof Profile) string {
	s := fmt.Sprintf("TotalKeys: %d, MaxDepth: %d, Types: {", prof.TotalKeys, prof.MaxDepth)
	first := true
	for k, v := range prof.TypeCounts {
		if !first {
			s += ", "
		}
		s += fmt.Sprintf("%s:%d", k, v)
		first = false
	}
	s += "}"
	return s
}
