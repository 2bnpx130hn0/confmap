package merger_test

import (
	"testing"

	"github.com/patrickward/confmap/merger"
)

func TestOverlay_BasicOverride(t *testing.T) {
	base := map[string]any{"host": "localhost", "port": 5432}
	patch := map[string]any{"host": "remotehost"}

	result := merger.NewOverlay().AddLayer(base).AddLayer(patch).Merge()

	if result["host"] != "remotehost" {
		t.Errorf("expected remotehost, got %v", result["host"])
	}
	if result["port"] != 5432 {
		t.Errorf("expected 5432, got %v", result["port"])
	}
}

func TestOverlay_SkipsNilValues(t *testing.T) {
	base := map[string]any{"key": "original"}
	patch := map[string]any{"key": nil}

	result := merger.NewOverlay().AddLayer(base).AddLayer(patch).Merge()

	if result["key"] != "original" {
		t.Errorf("expected original, got %v", result["key"])
	}
}

func TestOverlay_SkipsEmptyStrings(t *testing.T) {
	base := map[string]any{"name": "alice"}
	patch := map[string]any{"name": ""}

	result := merger.NewOverlay().AddLayer(base).AddLayer(patch).Merge()

	if result["name"] != "alice" {
		t.Errorf("expected alice, got %v", result["name"])
	}
}

func TestOverlay_NilLayerSkipped(t *testing.T) {
	base := map[string]any{"x": 1}

	result := merger.NewOverlay().AddLayer(base).AddLayer(nil).Merge()

	if result["x"] != 1 {
		t.Errorf("expected 1, got %v", result["x"])
	}
}

func TestOverlay_NestedMerge(t *testing.T) {
	base := map[string]any{
		"db": map[string]any{"host": "localhost", "port": 5432},
	}
	patch := map[string]any{
		"db": map[string]any{"host": "remotehost"},
	}

	result := merger.NewOverlay().AddLayer(base).AddLayer(patch).Merge()

	db, ok := result["db"].(map[string]any)
	if !ok {
		t.Fatal("expected db to be a map")
	}
	if db["host"] != "remotehost" {
		t.Errorf("expected remotehost, got %v", db["host"])
	}
	if db["port"] != 5432 {
		t.Errorf("expected 5432, got %v", db["port"])
	}
}

func TestOverlay_EmptyLayers(t *testing.T) {
	result := merger.NewOverlay().Merge()
	if len(result) != 0 {
		t.Errorf("expected empty map, got %v", result)
	}
}

func TestOverlay_DoesNotMutateBase(t *testing.T) {
	base := map[string]any{"a": "original"}
	patch := map[string]any{"a": "changed"}

	merger.NewOverlay().AddLayer(base).AddLayer(patch).Merge()

	if base["a"] != "original" {
		t.Errorf("base was mutated: got %v", base["a"])
	}
}
