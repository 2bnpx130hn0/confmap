// Package redactor provides functionality to mask or remove sensitive
// values from a config map based on key patterns or explicit key lists.
package redactor

import (
	"strings"
)

const redactedValue = "[REDACTED]"

// Redactor masks sensitive config values in place.
type Redactor struct {
	patterns []string
}

// New creates a Redactor that will mask values whose keys contain
// any of the provided patterns (case-insensitive substring match).
func New(patterns ...string) *Redactor {
	lower := make([]string, len(patterns))
	for i, p := range patterns {
		lower[i] = strings.ToLower(p)
	}
	return &Redactor{patterns: lower}
}

// Apply returns a deep copy of cfg with sensitive values replaced by
// the redacted placeholder string.
func (r *Redactor) Apply(cfg map[string]any) map[string]any {
	return r.redactMap(cfg)
}

func (r *Redactor) redactMap(m map[string]any) map[string]any {
	out := make(map[string]any, len(m))
	for k, v := range m {
		if r.isSensitive(k) {
			out[k] = redactedValue
			continue
		}
		switch val := v.(type) {
		case map[string]any:
			out[k] = r.redactMap(val)
		default:
			out[k] = v
		}
	}
	return out
}

func (r *Redactor) isSensitive(key string) bool {
	lk := strings.ToLower(key)
	for _, p := range r.patterns {
		if strings.Contains(lk, p) {
			return true
		}
	}
	return false
}

// Keys returns the list of keys (leaf keys only, dot-separated path)
// that were redacted in the last call to Apply.
func (r *Redactor) Keys(cfg map[string]any) []string {
	var found []string
	r.collectSensitive(cfg, "", &found)
	return found
}

func (r *Redactor) collectSensitive(m map[string]any, prefix string, out *[]string) {
	for k, v := range m {
		path := k
		if prefix != "" {
			path = prefix + "." + k
		}
		if r.isSensitive(k) {
			*out = append(*out, path)
			continue
		}
		if nested, ok := v.(map[string]any); ok {
			r.collectSensitive(nested, path, out)
		}
	}
}
