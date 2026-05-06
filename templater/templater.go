// Package templater provides variable interpolation for config values.
// It supports {{.KEY}} style placeholders that are resolved against a
// context map, enabling dynamic config composition at runtime.
package templater

import (
	"bytes"
	"fmt"
	"text/template"
)

// Templater resolves template expressions within config map values.
type Templater struct {
	ctx map[string]any
}

// New creates a Templater with the given context variables.
func New(ctx map[string]any) *Templater {
	return &Templater{ctx: ctx}
}

// Apply walks the config map and interpolates any string values that
// contain Go template expressions, returning a new resolved map.
func (t *Templater) Apply(cfg map[string]any) (map[string]any, error) {
	out := make(map[string]any, len(cfg))
	for k, v := range cfg {
		resolved, err := t.resolveValue(v)
		if err != nil {
			return nil, fmt.Errorf("templater: key %q: %w", k, err)
		}
		out[k] = resolved
	}
	return out, nil
}

func (t *Templater) resolveValue(v any) (any, error) {
	switch val := v.(type) {
	case string:
		return t.interpolate(val)
	case map[string]any:
		return t.Apply(val)
	case []any:
		result := make([]any, len(val))
		for i, item := range val {
			resolved, err := t.resolveValue(item)
			if err != nil {
				return nil, fmt.Errorf("index %d: %w", i, err)
			}
			result[i] = resolved
		}
		return result, nil
	default:
		return v, nil
	}
}

func (t *Templater) interpolate(s string) (string, error) {
	tmpl, err := template.New("").Option("missingkey=error").Parse(s)
	if err != nil {
		return "", fmt.Errorf("parse template: %w", err)
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, t.ctx); err != nil {
		return "", fmt.Errorf("execute template: %w", err)
	}
	return buf.String(), nil
}
