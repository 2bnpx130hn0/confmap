// Package pruner provides a configurable Pruner that strips unwanted entries
// from a config map.
//
// Three pruning modes are available and can be combined with bitwise OR:
//
//	- PruneNil   – removes keys whose value is nil (default when opts == 0)
//	- PruneEmpty – removes keys whose value is "", [], or {}
//	- PruneZero  – removes keys whose value is 0, 0.0, or false
//
// Example:
//
//	p := pruner.New(pruner.PruneNil | pruner.PruneEmpty)
//	clean, err := p.Apply(cfg)
package pruner
