// Package deduper provides a Deduper type for removing redundant entries
// from configuration maps.
//
// It supports two deduplication strategies:
//
//   - DedupeKeys: removes keys that share the same flattened dot-path,
//     retaining only the first occurrence encountered during traversal.
//
//   - DedupeValues: removes sibling keys at the same map level whose
//     scalar values are identical, keeping only the last occurrence.
//
// Nested maps are handled recursively in both strategies. Neither method
// mutates the original configuration map.
//
// Example:
//
//	d := deduper.New()
//	clean, err := d.DedupeKeys(cfg)
//	if err != nil {
//		log.Fatal(err)
//	}
package deduper
