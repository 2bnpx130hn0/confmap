// Package caster provides type-safe accessors for values stored in a
// config map (map[string]interface{}). It is designed to work alongside
// the confmap resolver and merger packages, allowing callers to retrieve
// config values as concrete Go types without manual type assertions.
//
// Example usage:
//
//	c := caster.New(configMap)
//	port, err := c.Int("port")
//	debug, err := c.Bool("debug")
//	name, err := c.String("app.name")
package caster
