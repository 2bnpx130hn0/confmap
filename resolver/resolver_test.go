package resolver_test

import (
	"errors"
	"testing"

	"github.com/yourorg/confmap/resolver"
)

// stubLoader is a test double that returns a preset map or error.
type stubLoader struct {
	data map[string]interface{}
	err  error
}

func (s *stubLoader) Load() (map[string]interface{}, error) {
	return s.data, s.err
}

func TestResolver_MergesAndValidates(t *testing.T) {
	schema := map[string]interface{}{
		"fields": map[string]interface{}{
			"host": map[string]interface{}{"type": "string", "required": true},
			"port": map[string]interface{}{"type": "int", "required": true},
		},
	}

	r := resolver.New(
		schema,
		resolver.Source{Name: "base", Loader: &stubLoader{data: map[string]interface{}{"host": "localhost", "port": 8080}}},
		resolver.Source{Name: "override", Loader: &stubLoader{data: map[string]interface{}{"port": 9090}}},
	)

	cfg, err := r.Resolve()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg["host"] != "localhost" {
		t.Errorf("expected host=localhost, got %v", cfg["host"])
	}
	if cfg["port"] != 9090 {
		t.Errorf("expected port=9090, got %v", cfg["port"])
	}
}

func TestResolver_LoaderError(t *testing.T) {
	r := resolver.New(
		nil,
		resolver.Source{Name: "bad", Loader: &stubLoader{err: errors.New("read error")}},
	)

	_, err := r.Resolve()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestResolver_ValidationFailure(t *testing.T) {
	schema := map[string]interface{}{
		"fields": map[string]interface{}{
			"host": map[string]interface{}{"type": "string", "required": true},
		},
	}

	// "host" is missing — validation should fail.
	r := resolver.New(
		schema,
		resolver.Source{Name: "base", Loader: &stubLoader{data: map[string]interface{}{"port": 8080}}},
	)

	_, err := r.Resolve()
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}
}

func TestResolver_NoSchema(t *testing.T) {
	r := resolver.New(
		nil,
		resolver.Source{Name: "base", Loader: &stubLoader{data: map[string]interface{}{"key": "value"}}},
	)

	cfg, err := r.Resolve()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg["key"] != "value" {
		t.Errorf("expected key=value, got %v", cfg["key"])
	}
}
