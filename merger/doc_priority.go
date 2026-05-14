// Package merger provides utilities for merging layered configuration maps.
//
// # Priority Merger
//
// NewPriority constructs a [PriorityMerger] that orders layers by an explicit
// integer priority before merging. Layers with a higher Priority value win over
// layers with a lower Priority value — regardless of insertion order.
//
// Example:
//
//	pm := merger.NewPriority([]merger.PriorityLayer{
//	    {Priority: 1,  Data: defaults},
//	    {Priority: 50, Data: fileConfig},
//	    {Priority: 99, Data: envOverrides},
//	})
//	result := pm.Merge()
//
// Nil Data maps are silently skipped. Layers with equal priority are merged in
// the order they were added (last-added wins), consistent with the base [Merger]
// behaviour.
package merger
