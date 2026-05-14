// Package merger provides utilities for merging layered configuration maps.
//
// # Conditional Merging
//
// ConditionalMerger extends the standard Merger with predicate-gated layers.
// Each ConditionalLayer carries a Condition — a function evaluated against the
// accumulated config at the moment the layer is about to be applied. If the
// condition returns false the layer is skipped entirely.
//
// Built-in conditions:
//
//	 Always()             — always include the layer
//	 Never()              — always skip the layer
//	 HasKey(key)          — include when the accumulated config contains key
//	 ValueEquals(key, v)  — include when cfg[key] string-equals v
//
// A nil Condition is treated as Always.
//
// Example:
//
//	cm := merger.NewConditional([]merger.ConditionalLayer{
//	    {Layer: base,    Condition: merger.Always()},
//	    {Layer: staging, Condition: merger.ValueEquals("env", "staging")},
//	    {Layer: prod,    Condition: merger.ValueEquals("env", "prod")},
//	})
//	result := cm.Merge()
package merger
