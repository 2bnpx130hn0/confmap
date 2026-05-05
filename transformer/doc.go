// Package transformer provides composable transformation functions
// for mutating config maps before they are validated or consumed.
//
// # Overview
//
// A [Transformer] holds an ordered list of [TransformFunc] values.
// Calling Apply runs each function in sequence, returning on the
// first error.
//
// # Built-in helpers
//
//   - [SetDefault] — fills in a missing key with a fallback value.
//   - [Rename]     — renames one key to another.
//   - [CoerceString] — converts a non-string value to its string
//     representation via fmt.Sprintf.
//
// Custom transforms can be provided by implementing [TransformFunc].
package transformer
