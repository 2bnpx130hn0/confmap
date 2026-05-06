// Package tagging provides lightweight key-level annotation for config maps.
//
// A Tagger stores arbitrary string tags against config key names. Tags can
// be used to drive downstream behaviour such as:
//
//   - Redacting or masking sensitive values before logging
//   - Identifying required vs optional keys for documentation
//   - Filtering subsets of a config for specific consumers
//
// Example:
//
//	tr := tagging.New()
//	tr.Tag("db.password", "secret", "required")
//	tr.Tag("app.debug",   "optional")
//
//	secrets := tr.FilterByTag(cfg, "secret")
package tagging
