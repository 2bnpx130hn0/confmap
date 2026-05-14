package resolver

import (
	"github.com/yourusername/confmap/stringifier"
)

// StringifyingResolver wraps a Resolver and exposes helpers to render the
// resolved configuration as formatted strings.
type StringifyingResolver struct {
	inner     *Resolver
	stringify *stringifier.Stringifier
}

// NewStringifying creates a StringifyingResolver using the provided Resolver
// and stringifier Options.
func NewStringifying(r *Resolver, opts stringifier.Options) *StringifyingResolver {
	return &StringifyingResolver{
		inner:     r,
		stringify: stringifier.New(opts),
	}
}

// Resolve returns the merged and validated config map.
func (sr *StringifyingResolver) Resolve() (map[string]any, error) {
	return sr.inner.Resolve()
}

// Lines resolves the config and returns it as a slice of formatted strings.
func (sr *StringifyingResolver) Lines() ([]string, error) {
	cfg, err := sr.inner.Resolve()
	if err != nil {
		return nil, err
	}
	return sr.stringify.Render(cfg), nil
}

// Text resolves the config and returns it as a single formatted string.
func (sr *StringifyingResolver) Text() (string, error) {
	cfg, err := sr.inner.Resolve()
	if err != nil {
		return "", err
	}
	return sr.stringify.String(cfg), nil
}
