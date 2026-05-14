// Package scoper provides scoped views into a configuration map.
//
// A Scoper is created from an existing config map and a prefix key.
// It exposes Get, Set, Keys, and Snapshot operations restricted to
// the subtree rooted at that prefix, preventing accidental access
// to unrelated parts of the configuration.
//
// Example:
//
//	data := map[string]any{
//		"database": map[string]any{
//			"host": "localhost",
//			"port": 5432,
//		},
//	}
//	sc, _ := scoper.New("database", data)
//	host, _ := sc.Get("host") // "localhost"
package scoper
