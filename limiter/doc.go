// Package limiter enforces structural constraints on configuration maps.
//
// It supports two types of limits:
//
//   - MaxKeys: the total number of keys (including nested) must not exceed this value.
//   - MaxDepth: the nesting depth of the map must not exceed this value.
//
// A limit value of 0 means that dimension is unconstrained.
//
// Example:
//
//	l := limiter.New(100, 5)
//	if err := l.Check(cfg); err != nil {
//		log.Fatal(err)
//	}
package limiter
