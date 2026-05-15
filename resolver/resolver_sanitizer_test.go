package resolver_test

import (
	"errors"
	"testing"

	"github.com/yourusername/confmap/loader"
	"github.com/yourusername/confmap/merger"
	"github.com/yourusername/confmap/resolver"
	"github.com/yourusername/confmap/sanitizer"
)

func makeSanitizedResolver(data map[string]any, schema map[string]any, rules ...sanitizer.Rule) *resolver.SanitizedResolver {
	l := loader.NewStaticLoader(data)
	m := merger.New()
	r := resolver.New([]loader.Loader{l}, m, nil)
	return resolver.NewSanitized(r, schema, rules...)
}

func TestSanitizedResolver_TrimsStrings(t *testing.T) {
	data := map[string]any{"host": "  localhost  "}
	sr := makeSanitizedResolver(data, nil, sanitizer.TrimSpace)
	out, err := sr.Resolve()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["host"] != "localhost" {
		t.Errorf("expected 'localhost', got %q", out["host"])
	}
}

func TestSanitizedResolver_LowercasesStrings(t *testing.T) {
	data := map[string]any{"env": "STAGING"}
	sr := makeSanitizedResolver(data, nil, sanitizer.ToLower)
	out, err := sr.Resolve()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["env"] != "staging" {
		t.Errorf("expected 'staging', got %q", out["env"])
	}
}

func TestSanitizedResolver_ValidationPasses(t *testing.T) {
	data := map[string]any{"host": "  db  "}
	schema := map[string]any{"host": map[string]any{"type": "string", "required": true}}
	sr := makeSanitizedResolver(data, schema, sanitizer.TrimSpace)
	_, err := sr.Resolve()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSanitizedResolver_LoaderError(t *testing.T) {
	l := loader.NewErrorLoader(errors.New("load failure"))
	m := merger.New()
	r := resolver.New([]loader.Loader{l}, m, nil)
	sr := resolver.NewSanitized(r, nil, sanitizer.TrimSpace)
	_, err := sr.Resolve()
	if err == nil {
		t.Fatal("expected error from failing loader")
	}
}

func TestSanitizedResolver_NoRulesPassthrough(t *testing.T) {
	data := map[string]any{"key": "  value  "}
	sr := makeSanitizedResolver(data, nil)
	out, err := sr.Resolve()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["key"] != "  value  " {
		t.Errorf("expected untouched value, got %q", out["key"])
	}
}
