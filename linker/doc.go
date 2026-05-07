// Package linker provides cross-key reference resolution for config maps.
//
// Values may reference other keys using ${key.path} syntax. References are
// resolved recursively and in dependency order. Cycles are detected and
// reported as errors.
//
// Example:
//
//	cfg := map[string]any{
//		"base_url": "https://example.com",
//		"api_url":  "${base_url}/api",
//	}
//	l := linker.New(cfg)
//	resolved, err := l.Resolve()
//	// resolved["api_url"] == "https://example.com/api"
package linker
