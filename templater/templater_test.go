package templater_test

import (
	"testing"

	"github.com/your-org/confmap/templater"
)

func TestApply_SimpleInterpolation(t *testing.T) {
	tmplr := templater.New(map[string]any{"ENV": "production"})
	cfg := map[string]any{
		"mode": "{{.ENV}}",
	}
	out, err := tmplr.Apply(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["mode"] != "production" {
		t.Errorf("expected 'production', got %v", out["mode"])
	}
}

func TestApply_NoTemplate(t *testing.T) {
	tmplr := templater.New(map[string]any{})
	cfg := map[string]any{"key": "static-value", "count": 42}
	out, err := tmplr.Apply(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["key"] != "static-value" {
		t.Errorf("expected 'static-value', got %v", out["key"])
	}
	if out["count"] != 42 {
		t.Errorf("expected 42, got %v", out["count"])
	}
}

func TestApply_NestedMap(t *testing.T) {
	tmplr := templater.New(map[string]any{"HOST": "localhost", "PORT": "5432"})
	cfg := map[string]any{
		"database": map[string]any{
			"host": "{{.HOST}}",
			"port": "{{.PORT}}",
		},
	}
	out, err := tmplr.Apply(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	db, ok := out["database"].(map[string]any)
	if !ok {
		t.Fatal("expected nested map")
	}
	if db["host"] != "localhost" {
		t.Errorf("expected 'localhost', got %v", db["host"])
	}
	if db["port"] != "5432" {
		t.Errorf("expected '5432', got %v", db["port"])
	}
}

func TestApply_MissingContextKey(t *testing.T) {
	tmplr := templater.New(map[string]any{})
	cfg := map[string]any{"url": "http://{{.HOST}}/api"}
	_, err := tmplr.Apply(cfg)
	if err == nil {
		t.Error("expected error for missing context key, got nil")
	}
}

func TestApply_SliceValues(t *testing.T) {
	tmplr := templater.New(map[string]any{"REGION": "us-east-1"})
	cfg := map[string]any{
		"regions": []any{"{{.REGION}}", "eu-west-1"},
	}
	out, err := tmplr.Apply(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	regions, ok := out["regions"].([]any)
	if !ok {
		t.Fatal("expected slice")
	}
	if regions[0] != "us-east-1" {
		t.Errorf("expected 'us-east-1', got %v", regions[0])
	}
	if regions[1] != "eu-west-1" {
		t.Errorf("expected 'eu-west-1', got %v", regions[1])
	}
}

func TestApply_InvalidTemplate(t *testing.T) {
	tmplr := templater.New(map[string]any{})
	cfg := map[string]any{"bad": "{{.UNCLOSED"}
	_, err := tmplr.Apply(cfg)
	if err == nil {
		t.Error("expected parse error for invalid template, got nil")
	}
}
