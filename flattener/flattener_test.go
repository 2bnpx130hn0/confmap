package flattener_test

import (
	"testing"

	"github.com/dev/confmap/flattener"
)

func TestFlatten_SimpleNested(t *testing.T) {
	f := flattener.New(".")
	input := map[string]any{
		"database": map[string]any{
			"host": "localhost",
			"port": 5432,
		},
	}
	out := f.Flatten(input)
	if out["database.host"] != "localhost" {
		t.Errorf("expected database.host=localhost, got %v", out["database.host"])
	}
	if out["database.port"] != 5432 {
		t.Errorf("expected database.port=5432, got %v", out["database.port"])
	}
}

func TestFlatten_FlatMap(t *testing.T) {
	f := flattener.New(".")
	input := map[string]any{"key": "value", "count": 3}
	out := f.Flatten(input)
	if len(out) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(out))
	}
	if out["key"] != "value" {
		t.Errorf("expected key=value, got %v", out["key"])
	}
}

func TestFlatten_DeepNesting(t *testing.T) {
	f := flattener.New(".")
	input := map[string]any{
		"a": map[string]any{
			"b": map[string]any{
				"c": true,
			},
		},
	}
	out := f.Flatten(input)
	if out["a.b.c"] != true {
		t.Errorf("expected a.b.c=true, got %v", out["a.b.c"])
	}
}

func TestExpand_SimpleFlat(t *testing.T) {
	f := flattener.New(".")
	input := map[string]any{
		"database.host": "localhost",
		"database.port": 5432,
	}
	out, err := f.Expand(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	db, ok := out["database"].(map[string]any)
	if !ok {
		t.Fatalf("expected database to be a map")
	}
	if db["host"] != "localhost" {
		t.Errorf("expected host=localhost, got %v", db["host"])
	}
}

func TestExpand_ConflictError(t *testing.T) {
	f := flattener.New(".")
	input := map[string]any{
		"a":   "scalar",
		"a.b": "nested",
	}
	_, err := f.Expand(input)
	if err == nil {
		t.Error("expected error when scalar conflicts with nested key")
	}
}

func TestFlattenExpand_RoundTrip(t *testing.T) {
	f := flattener.New(".")
	original := map[string]any{
		"server": map[string]any{
			"host": "example.com",
			"tls":  true,
		},
		"timeout": 30,
	}
	flat := f.Flatten(original)
	restored, err := f.Expand(flat)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	srv, ok := restored["server"].(map[string]any)
	if !ok {
		t.Fatalf("expected server to be a map after round-trip")
	}
	if srv["host"] != "example.com" {
		t.Errorf("expected host=example.com, got %v", srv["host"])
	}
	if restored["timeout"] != 30 {
		t.Errorf("expected timeout=30, got %v", restored["timeout"])
	}
}

func TestNew_DefaultSeparator(t *testing.T) {
	f := flattener.New("")
	input := map[string]any{"x": map[string]any{"y": 1}}
	out := f.Flatten(input)
	if _, ok := out["x.y"]; !ok {
		t.Error("expected default separator '.' to be used")
	}
}
