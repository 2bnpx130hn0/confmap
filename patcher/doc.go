// Package patcher provides dot-notation key-path patching for configuration maps.
//
// It allows callers to apply targeted updates or deletions to nested
// configuration structures without replacing the entire map. Intermediate
// nodes are created automatically when applying a value to a deep path.
//
// Example usage:
//
//	cfg := map[string]any{"database": map[string]any{"host": "localhost"}}
//	p := patcher.New(cfg)
//	_ = p.Apply("database.port", 5432)
//	_ = p.Apply("feature.flags.dark_mode", true)
//	_ = p.Delete("database.host")
//	updated := p.Config()
package patcher
