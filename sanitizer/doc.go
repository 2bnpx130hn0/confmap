// Package sanitizer applies ordered transformation rules to string values
// within a config map, supporting nested maps and slices.
//
// # Overview
//
// A Sanitizer holds a list of Rule functions. When Apply is called, every
// string value in the config (including inside nested maps and slices) is
// passed through each rule in the order they were registered.
//
// # Built-in Rules
//
//   - TrimSpace  – strips leading/trailing whitespace
//   - ToLower    – lowercases the entire value
//   - ToUpper    – uppercases the entire value
//   - StripNull  – replaces the literal string "null" with ""
//   - ReplaceRule(old, new) – replaces all occurrences of old with new
//
// # Example
//
//	s := sanitizer.New(sanitizer.TrimSpace, sanitizer.ToLower)
//	clean := s.Apply(cfg)
package sanitizer
