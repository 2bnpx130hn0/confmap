// Package coercer provides type coercion for config map values.
//
// It allows registering per-key coercion rules (ToString, ToInt, ToBool)
// and applying them to a config map in a single pass. Useful for normalising
// values loaded from environment variables or loosely-typed YAML/TOML files
// before validation or further processing.
//
// Example:
//
//	c := coercer.New().
//		ToInt("port").
//		ToBool("debug").
//		ToString("version")
//
//	if err := c.Apply(cfg); err != nil {
//		log.Fatal(err)
//	}
package coercer
