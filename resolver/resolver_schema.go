package resolver

import (
	"fmt"

	"github.com/user/confmap/loader"
	"github.com/user/confmap/merger"
	"github.com/user/confmap/validator"
)

// SchemaVersioned is a Resolver that uses a SchemaMerger to guarantee all
// loaded layers share the same major schema version before validation.
type SchemaVersioned struct {
	loaders  []loader.VersionedLoader
	schema   map[string]any
}

// VersionedLoader pairs a Loader with the schema version it claims.
type VersionedLoader = loader.VersionedLoader

// NewSchemaVersioned constructs a resolver that enforces schema-version
// compatibility across all config layers.
func NewSchemaVersioned(schema map[string]any, loaders ...loader.VersionedLoader) *SchemaVersioned {
	return &SchemaVersioned{loaders: loaders, schema: schema}
}

// Resolve loads every layer, checks version compatibility via SchemaMerger,
// then validates the merged result against the provided schema.
func (sv *SchemaVersioned) Resolve() (map[string]any, error) {
	sm := merger.NewSchema()

	for _, vl := range sv.loaders {
		data, err := vl.Loader.Load()
		if err != nil {
			return nil, fmt.Errorf("resolver/schema: loader error: %w", err)
		}
		sm.AddLayer(vl.Version, data)
	}

	merged, err := sm.Merge()
	if err != nil {
		return nil, fmt.Errorf("resolver/schema: merge error: %w", err)
	}

	if sv.schema != nil {
		if err := validator.Validate(merged, sv.schema); err != nil {
			return nil, fmt.Errorf("resolver/schema: validation error: %w", err)
		}
	}

	return merged, nil
}
