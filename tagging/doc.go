// Package tagging provides lightweight key-level annotation for config maps.
//
// A Tagger stores arbitrary string tags against config key names. Tags can
// be used to drive downstream behaviour such as:
//
//   - Redacting or masking sensitive values before logging
//   - Identifying required vs optional keys for documentation
//   - Filtering subsets of a config for specific consumers
//
// Tags are case-sensitive and order-preserving. A key may carry any number
// of tags; duplicate tags on the same key are silently deduplicated.
//
// Example:
//
//	tr := tagging.New()
//	tr.Tag("db.password", "secret", "required")
//	tr.Tag("app.debug",   "optional")
//
//	// Retrieve all tags attached to a key.
//	tags := tr.Tags("db.password") // ["secret", "required"]
//
//	// Check whether a key carries a specific tag.
//	if tr.HasTag("db.password", "secret") {
//		// redact before logging
//	}
//
//	// Collect all keys that share a tag.
//	secrets := tr.FilterByTag(cfg, "secret")
package tagging
