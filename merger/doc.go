// Package merger provides utilities for merging layered configuration maps
// with configurable override precedence.
//
// # Overview
//
// When building applications with multiple config sources (files, env vars,
// defaults), it is common to need a way to combine them with a clear
// precedence order. The merger package handles this by accepting an ordered
// slice of config maps and combining them according to the chosen strategy.
//
// # Strategies
//
// StrategyOverride (default): later layers overwrite earlier values.
// Nested maps are merged recursively so that only conflicting leaf keys
// are replaced.
//
// StrategyKeepExisting: the first value set for a key wins. Useful when
// you want defaults to act as fallbacks that cannot be overridden.
//
// # Example
//
//	m := merger.New(merger.StrategyOverride)
//	result, err := m.Merge(defaults, fileConfig, envConfig)
package merger
