// Package stringifier provides utilities for converting config maps into
// human-readable string representations.
//
// It supports:
//   - Flat and nested map rendering via dot-notation flattening
//   - Configurable key-value delimiters
//   - Optional value quoting
//   - Sorted key output for deterministic results
//   - Line prefix injection (e.g. for export or env-file generation)
//
// Example:
//
//	s := stringifier.New(stringifier.Options{
//		Delimiter:   "=",
//		QuoteValues: true,
//		SortKeys:    true,
//	})
//	lines := s.Render(cfg)
package stringifier
