// Package comparator provides utilities for comparing two config maps
// and determining equality, similarity scores, and structural differences.
package comparator

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
)

// Comparator holds options for config comparison.
type Comparator struct {
	ignoreKeys map[string]struct{}
}

// Result holds the outcome of a comparison.
type Result struct {
	Equal      bool
	Similarity float64 // 0.0 to 1.0
	OnlyInA    []string
	OnlyInB    []string
	Differing  []string
}

// New creates a new Comparator.
func New(ignoreKeys ...string) *Comparator {
	set := make(map[string]struct{}, len(ignoreKeys))
	for _, k := range ignoreKeys {
		set[strings.ToLower(k)] = struct{}{}
	}
	return &Comparator{ignoreKeys: set}
}

// Compare performs a deep comparison between two config maps.
func (c *Comparator) Compare(a, b map[string]any) Result {
	flatA := flatten(a, "")
	flatB := flatten(b, "")

	for k := range flatA {
		if _, ignored := c.ignoreKeys[strings.ToLower(k)]; ignored {
			delete(flatA, k)
		}
	}
	for k := range flatB {
		if _, ignored := c.ignoreKeys[strings.ToLower(k)]; ignored {
			delete(flatB, k)
		}
	}

	var onlyA, onlyB, differing []string

	for k, va := range flatA {
		if vb, ok := flatB[k]; !ok {
			onlyA = append(onlyA, k)
		} else if !reflect.DeepEqual(va, vb) {
			differing = append(differing, k)
		}
	}
	for k := range flatB {
		if _, ok := flatA[k]; !ok {
			onlyB = append(onlyB, k)
		}
	}

	sort.Strings(onlyA)
	sort.Strings(onlyB)
	sort.Strings(differing)

	total := len(flatA) + len(onlyB)
	var similarity float64
	if total > 0 {
		matching := len(flatA) - len(onlyA) - len(differing)
		similarity = float64(matching*2) / float64(total+len(flatA)-len(onlyA)-len(differing)+len(onlyB))
		if similarity < 0 {
			similarity = 0
		}
	} else {
		similarity = 1.0
	}

	return Result{
		Equal:      len(onlyA) == 0 && len(onlyB) == 0 && len(differing) == 0,
		Similarity: similarity,
		OnlyInA:    onlyA,
		OnlyInB:    onlyB,
		Differing:  differing,
	}
}

func flatten(m map[string]any, prefix string) map[string]any {
	out := make(map[string]any)
	for k, v := range m {
		key := k
		if prefix != "" {
			key = fmt.Sprintf("%s.%s", prefix, k)
		}
		if nested, ok := v.(map[string]any); ok {
			for nk, nv := range flatten(nested, key) {
				out[nk] = nv
			}
		} else {
			out[key] = v
		}
	}
	return out
}
