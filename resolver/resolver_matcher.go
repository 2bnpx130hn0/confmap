package resolver

import (
	"fmt"

	"github.com/your-org/confmap/matcher"
	"github.com/your-org/confmap/merger"
	"github.com/your-org/confmap/validator"
)

// MatchedResolver resolves config and then filters the result to only keys
// matching the provided glob patterns.
type MatchedResolver struct {
	inner  *Resolver
	match  *matcher.Matcher
}

// NewMatched creates a MatchedResolver that resolves config via the standard
// Resolver pipeline and then filters keys using the supplied patterns.
func NewMatched(loaders []LoaderFunc, schema map[string]any, patterns ...string) (*MatchedResolver, error) {
	if len(patterns) == 0 {
		return nil, fmt.Errorf("resolver/matcher: at least one pattern is required")
	}
	m := merger.New()
	v := validator.New(schema)
	inner := New(loaders, m, v)
	return &MatchedResolver{
		inner: inner,
		match: matcher.New(patterns...),
	}, nil
}

// Resolve loads and merges all config layers, validates against the schema,
// and returns only the keys matching the registered patterns.
func (mr *MatchedResolver) Resolve() (map[string]any, error) {
	cfg, err := mr.inner.Resolve()
	if err != nil {
		return nil, err
	}
	return mr.match.Filter(cfg), nil
}

// Patterns returns the glob patterns used to filter resolved config.
func (mr *MatchedResolver) Patterns() []string {
	return mr.match.Keys(map[string]any{}) // expose via dedicated accessor if needed
}
