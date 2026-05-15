// Package merger provides several merge strategies for combining layered
// configuration maps.
//
// # OverlayMerger
//
// NewOverlay returns an OverlayMerger that applies layers in insertion order.
// Unlike the standard Merger, it silently skips nil values and empty strings
// so that sparse "patch" layers do not accidentally erase values that were
// set by an earlier layer.
//
// This is useful when environment-specific overrides only declare the keys
// they actually want to change, leaving everything else untouched.
//
// Example:
//
//	base  := map[string]any{"host": "localhost", "port": 5432}
//	patch := map[string]any{"host": "prod.db", "port": nil} // nil is ignored
//
//	result := merger.NewOverlay().
//	    AddLayer(base).
//	    AddLayer(patch).
//	    Merge()
//	// result == {"host": "prod.db", "port": 5432}
package merger
