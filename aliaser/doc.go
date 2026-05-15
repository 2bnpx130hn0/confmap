// Package aliaser provides transparent key aliasing for config maps.
//
// It allows registering alternate names (aliases) for canonical config keys,
// so that consumers can use either name when specifying configuration values.
//
// Example:
//
//	a := aliaser.New()
//	_ = a.Register("db_host", "database.host")
//
//	cfg := map[string]any{"db_host": "localhost"}
//	resolved, _ := a.Apply(cfg)
//	// resolved == map[string]any{"database.host": "localhost"}
//
// Canonical keys already present in the config take precedence over alias values.
package aliaser
