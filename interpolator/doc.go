// Package interpolator expands ${VAR} and $VAR style references within config
// map string values.
//
// Variables can be sourced from an explicit map or from the process environment.
//
// Example:
//
//	 interp := interpolator.New(map[string]string{"HOST": "localhost"}, true)
//	 out, err := interp.Apply(map[string]any{"addr": "${HOST}:8080"})
//	 // out["addr"] == "localhost:8080"
//
// Strict mode returns an error when a referenced variable is not defined.
// Non-strict mode silently replaces missing variables with an empty string.
package interpolator
