// Package comparator provides deep comparison between two configuration maps.
//
// It flattens nested structures, supports key ignoring (e.g. timestamps or
// volatile fields), and returns a structured Result with keys only present
// in each side, keys whose values differ, and an overall similarity score.
//
// Example:
//
//	cmp := comparator.New("updated_at", "version")
//	result := cmp.Compare(configA, configB)
//	if !result.Equal {
//		fmt.Println("Differing keys:", result.Differing)
//	}
package comparator
