// Package merger provides several merge strategies for combining layered
// configuration maps.
//
// SchemaMerger
//
// SchemaMerger extends the basic merge behaviour with schema-version
// awareness. Each layer is tagged with a semantic version string. During
// Merge the merger checks that all layers share the same major version
// component. If any two layers disagree on the major version an error is
// returned, preventing accidentally mixing incompatible configuration
// schemas.
//
// Usage:
//
//	sm := merger.NewSchema()
//	sm.AddLayer("2.0.0", baseConfig)
//	sm.AddLayer("2.3.1", overrideConfig)
//	cfg, err := sm.Merge()
//
// Layers are applied in the order they were added; later layers win on
// key conflicts, matching the behaviour of the standard Merger.
package merger
