package resolver

import (
	"testing"

	"github.com/your-org/confmap/merger"
)

func TestConditionalResolver_MergesConditionally(t *testing.T) {
	layers := []merger.ConditionalLayer{
		{Layer: map[string]any{"env": "staging", "port": 8080}, Condition: merger.Always()},
		{Layer: map[string]any{"debug": true}, Condition: merger.ValueEquals("env", "staging")},
		{Layer: map[string]any{"debug": false}, Condition: merger.ValueEquals("env", "prod")},
	}
	r := NewConditional(layers, nil)
	cfg, err := r.Resolve()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg["debug"] != true {
		t.Errorf("expected debug=true for staging, got %v", cfg["debug"])
	}
	if cfg["port"] != 8080 {
		t.Errorf("expected port=8080, got %v", cfg["port"])
	}
}

func TestConditionalResolver_ValidationFailure(t *testing.T) {
	layers := []merger.ConditionalLayer{
		{Layer: map[string]any{"port": 9090}, Condition: merger.Always()},
	}
	schema := map[string]any{
		"required": []any{"host"},
	}
	r := NewConditional(layers, schema)
	_, err := r.Resolve()
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}
}

func TestConditionalResolver_NoSchema(t *testing.T) {
	layers := []merger.ConditionalLayer{
		{Layer: map[string]any{"key": "val"}, Condition: merger.Always()},
	}
	r := NewConditional(layers, nil)
	cfg, err := r.Resolve()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg["key"] != "val" {
		t.Errorf("expected key=val, got %v", cfg["key"])
	}
}

func TestConditionalResolver_SkippedLayerAbsent(t *testing.T) {
	layers := []merger.ConditionalLayer{
		{Layer: map[string]any{"base": 1}, Condition: merger.Always()},
		{Layer: map[string]any{"secret": "hidden"}, Condition: merger.Never()},
	}
	r := NewConditional(layers, nil)
	cfg, err := r.Resolve()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := cfg["secret"]; ok {
		t.Error("expected secret key to be absent")
	}
}
