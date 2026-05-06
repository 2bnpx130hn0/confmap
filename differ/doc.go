// Package differ provides utilities for computing the difference between
// two configuration maps.
//
// It flattens nested maps using dot-notation keys and identifies:
//   - Added keys: present in the new config but not in the old.
//   - Removed keys: present in the old config but not in the new.
//   - Changed keys: present in both but with different values.
//
// Example usage:
//
//	d := differ.New()
//	deltas := d.Diff(oldConfig, newConfig)
//	for _, delta := range deltas {
//		fmt.Printf("[%s] %s: %v -> %v\n", delta.Type, delta.Key, delta.OldValue, delta.NewValue)
//	}
package differ
