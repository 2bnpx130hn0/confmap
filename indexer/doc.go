// Package indexer builds a flat dot-notation index over a nested
// config map produced by confmap loaders or the merger.
//
// Example usage:
//
//	cfg := map[string]any{
//		"database": map[string]any{
//			"host": "localhost",
//			"port": 5432,
//		},
//	}
//
//	idx := indexer.New(cfg)
//	host, _ := idx.Get("database.host") // "localhost"
//	port, _ := idx.Get("database.port") // 5432
package indexer
