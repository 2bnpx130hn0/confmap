package cloner_test

import (
	"testing"

	"github.com/example/confmap/cloner"
)

func TestClone_FlatMap(t *testing.T) {
	c := cloner.New()
	orig := map[string]any{"host": "localhost", "port": 8080}
	cloned, err := c.Clone(orig)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cloned["host"] != "localhost" || cloned["port"] != 8080 {
		t.Errorf("cloned values do not match original")
	}
	cloned["host"] = "changed"
	if orig["host"] == "changed" {
		t.Error("mutation of clone affected original")
	}
}

func TestClone_NestedMap(t *testing.T) {
	c := cloner.New()
	orig := map[string]any{
		"db": map[string]any{"host": "db.local", "port": 5432},
	}
	cloned, err := c.Clone(orig)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	inner := cloned["db"].(map[string]any)
	inner["host"] = "mutated"
	orig["db"].(map[string]any)["host"] == "mutated"
	if orig["db"].(map[string]any)["host"] == "mutated" {
		t.Error("nested mutation of clone affected original")
	}
}

func TestClone_SliceValues(t *testing.T) {
	c := cloner.New()
	orig := map[string]any{"tags": []any{"a", "b", "c"}}
	cloned, err := c.Clone(orig)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	slice := cloned["tags"].([]any)
	slice[0] = "z"
	if orig["tags"].([]any)[0] == "z" {
		t.Error("slice mutation of clone affected original")
	}
}

func TestClone_UnsupportedType(t *testing.T) {
	c := cloner.New()
	type custom struct{ X int }
	orig := map[string]any{"obj": custom{X: 1}}
	_, err := c.Clone(orig)
	if err == nil {
		t.Error("expected error for unsupported type, got nil")
	}
}

func TestMustClone_Panics(t *testing.T) {
	c := cloner.New()
	type bad struct{}
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for unsupported type")
		}
	}()
	c.MustClone(map[string]any{"x": bad{}})
}

func TestClone_EmptyMap(t *testing.T) {
	c := cloner.New()
	cloned, err := c.Clone(map[string]any{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cloned) != 0 {
		t.Errorf("expected empty map, got %v", cloned)
	}
}
