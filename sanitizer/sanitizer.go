// Package sanitizer provides utilities for cleaning and normalizing config
// map values by applying user-defined sanitization rules to string fields.
package sanitizer

import (
	"strings"
)

// Rule is a function that transforms a string value.
type Rule func(string) string

// Sanitizer applies a chain of Rules to all string values in a config map.
type Sanitizer struct {
	rules []Rule
}

// New creates a new Sanitizer with the given rules applied in order.
func New(rules ...Rule) *Sanitizer {
	return &Sanitizer{rules: rules}
}

// Apply returns a new config map with all string values sanitized.
// Non-string and nested map values are handled recursively.
func (s *Sanitizer) Apply(cfg map[string]any) map[string]any {
	if cfg == nil {
		return nil
	}
	out := make(map[string]any, len(cfg))
	for k, v := range cfg {
		out[k] = s.sanitizeValue(v)
	}
	return out
}

func (s *Sanitizer) sanitizeValue(v any) any {
	switch val := v.(type) {
	case string:
		for _, r := range s.rules {
			val = r(val)
		}
		return val
	case map[string]any:
		return s.Apply(val)
	case []any:
		out := make([]any, len(val))
		for i, item := range val {
			out[i] = s.sanitizeValue(item)
		}
		return out
	default:
		return v
	}
}

// TrimSpace is a Rule that trims leading and trailing whitespace.
func TrimSpace(s string) string { return strings.TrimSpace(s) }

// ToLower is a Rule that converts a string to lowercase.
func ToLower(s string) string { return strings.ToLower(s) }

// ToUpper is a Rule that converts a string to uppercase.
func ToUpper(s string) string { return strings.ToUpper(s) }

// StripNull is a Rule that replaces the literal string "null" with an empty string.
func StripNull(s string) string {
	if strings.EqualFold(s, "null") {
		return ""
	}
	return s
}

// ReplaceRule returns a Rule that replaces all occurrences of old with new.
func ReplaceRule(old, new string) Rule {
	return func(s string) string {
		return strings.ReplaceAll(s, old, new)
	}
}
