// Package matcher provides glob-pattern-based key matching for config maps.
//
// It allows callers to filter or inspect config entries by matching their
// dotted key paths against one or more wildcard patterns.
//
// Example usage:
//
//	m := matcher.New("db.*", "app.port")
//
//	cfg := map[string]any{
//		"db":  map[string]any{"host": "localhost", "port": 5432},
//		"app": map[string]any{"port": 8080, "name": "myapp"},
//	}
//
//	filtered := m.Filter(cfg)
//	// filtered contains db.host, db.port, and app.port
//
//	keys := m.Keys(cfg)
//	// keys: ["db.host", "db.port", "app.port"]
package matcher
