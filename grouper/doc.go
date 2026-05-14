// Package grouper partitions a config map into named sub-maps based on key
// prefixes or a custom grouping function.
//
// Example usage:
//
//	g := grouper.New(map[string]any{
//		"db_host":    "localhost",
//		"db_port":    5432,
//		"cache_host": "redis",
//		"debug":      true,
//	})
//
//	groups := g.ByPrefix([]string{"db", "cache"}, "_")
//	// groups["db"]    => {"host": "localhost", "port": 5432}
//	// groups["cache"] => {"host": "redis"}
//	// groups[""]      => {"debug": true}
package grouper
