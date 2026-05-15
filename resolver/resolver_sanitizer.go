package resolver

import (
	"fmt"

	"github.com/yourusername/confmap/sanitizer"
	"github.com/yourusername/confmap/validator"
)

// SanitizedResolver resolves config and applies sanitization rules before
// optional schema validation.
type SanitizedResolver struct {
	inner  *Resolver
	san    *sanitizer.Sanitizer
	schema map[string]any
}

// NewSanitized returns a SanitizedResolver that runs sanitization rules on the
// merged config. Pass a nil schema to skip validation.
func NewSanitized(r *Resolver, schema map[string]any, rules ...sanitizer.Rule) *SanitizedResolver {
	return &SanitizedResolver{
		inner:  r,
		san:    sanitizer.New(rules...),
		schema: schema,
	}
}

// Resolve merges all layers, sanitizes string values, then validates against
// the schema (if provided).
func (sr *SanitizedResolver) Resolve() (map[string]any, error) {
	raw, err := sr.inner.Resolve()
	if err != nil {
		return nil, fmt.Errorf("sanitized resolver: %w", err)
	}

	clean := sr.san.Apply(raw)

	if sr.schema != nil {
		if err := validator.Validate(clean, sr.schema); err != nil {
			return nil, fmt.Errorf("sanitized resolver validation: %w", err)
		}
	}

	return clean, nil
}
