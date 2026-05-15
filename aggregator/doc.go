// Package aggregator provides a flexible reduce-style mechanism for combining
// multiple config layers into a single map.
//
// Unlike the merger package which has fixed merge semantics, aggregator allows
// callers to supply a custom reduce function, enabling use-cases such as key
// collection, value counting, or conditional accumulation.
//
// Basic usage:
//
//	agg := aggregator.New(aggregator.MergeReduce)
//	agg.Add(baseConfig)
//	agg.Add(overrideConfig)
//	result, err := agg.Reduce()
//
// Built-in reduce functions:
//   - MergeReduce  – later layers override earlier keys.
//   - CollectKeys  – accumulates all unique keys as boolean flags.
package aggregator
