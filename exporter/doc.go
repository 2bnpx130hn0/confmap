// Package exporter serializes a resolved config map into YAML, TOML, or JSON.
//
// It is intended to be used after config loading, merging, and validation to
// persist or display the final effective configuration.
//
// Example usage:
//
//	cfg := map[string]any{
//		"host": "localhost",
//		"port": 8080,
//	}
//
//	out, err := exporter.Export(cfg, exporter.FormatYAML)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(string(out))
package exporter
