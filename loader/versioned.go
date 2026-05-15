package loader

// VersionedLoader pairs a Loader implementation with a semantic version
// string that describes which schema version the loaded data conforms to.
// It is consumed by resolver.NewSchemaVersioned to enforce cross-layer
// version compatibility.
type VersionedLoader struct {
	// Version is a semver string such as "2.1.0". Only the major component
	// is compared for compatibility.
	Version string

	// Loader is the underlying source that produces the config map.
	Loader Loader
}

// NewVersionedLoader is a convenience constructor.
func NewVersionedLoader(version string, l Loader) VersionedLoader {
	return VersionedLoader{Version: version, Loader: l}
}
