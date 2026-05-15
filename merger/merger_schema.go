package merger

import "fmt"

// SchemaLayer pairs a config layer with a semantic version string.
type SchemaLayer struct {
	Version string
	Data    map[string]any
}

// SchemaMerger merges layers while enforcing that every layer declares a
// compatible schema version (same major component).
type SchemaMerger struct {
	layers []SchemaLayer
}

// NewSchema returns a new SchemaMerger.
func NewSchema() *SchemaMerger {
	return &SchemaMerger{}
}

// AddLayer registers a versioned config layer.
func (s *SchemaMerger) AddLayer(version string, data map[string]any) {
	s.layers = append(s.layers, SchemaLayer{Version: version, Data: data})
}

// Merge combines all registered layers in insertion order.
// It returns an error if any two layers declare incompatible major versions.
func (s *SchemaMerger) Merge() (map[string]any, error) {
	if len(s.layers) == 0 {
		return map[string]any{}, nil
	}

	baseMajor, err := majorVersion(s.layers[0].Version)
	if err != nil {
		return nil, fmt.Errorf("merger/schema: base layer version %q invalid: %w", s.layers[0].Version, err)
	}

	result := map[string]any{}
	for _, layer := range s.layers {
		maj, err := majorVersion(layer.Version)
		if err != nil {
			return nil, fmt.Errorf("merger/schema: layer version %q invalid: %w", layer.Version, err)
		}
		if maj != baseMajor {
			return nil, fmt.Errorf("merger/schema: incompatible major versions %q vs %q", s.layers[0].Version, layer.Version)
		}
		if layer.Data == nil {
			continue
		}
		for k, v := range layer.Data {
			result[k] = v
		}
	}
	return result, nil
}

// majorVersion extracts the leading integer from a semver string like "2.1.0".
func majorVersion(v string) (string, error) {
	if v == "" {
		return "", fmt.Errorf("empty version string")
	}
	for i, ch := range v {
		if ch == '.' {
			return v[:i], nil
		}
	}
	return v, nil // no dot → treat whole string as major
}
