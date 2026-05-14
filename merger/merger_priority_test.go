package merger_test

import (
	"testing"

	"github.com/user/confmap/merger"
)

func TestPriorityMerge_HigherWins(t *testing.T) {
	layers := []merger.PriorityLayer{
		{Priority: 10, Data: map[string]any{"env": "prod", "debug": false}},
		{Priority: 1, Data: map[string]any{"env": "dev", "verbose": true}},
	}
	pm := merger.NewPriority(layers)
	result := pm.Merge()

	if result["env"] != "prod" {
		t.Errorf("expected env=prod, got %v", result["env"])
	}
	if result["verbose"] != true {
		t.Errorf("expected verbose=true, got %v", result["verbose"])
	}
	if result["debug"] != false {
		t.Errorf("expected debug=false, got %v", result["debug"])
	}
}

func TestPriorityMerge_LowerDoesNotOverride(t *testing.T) {
	layers := []merger.PriorityLayer{
		{Priority: 5, Data: map[string]any{"key": "base"}},
		{Priority: 1, Data: map[string]any{"key": "low"}},
	}
	pm := merger.NewPriority(layers)
	result := pm.Merge()

	if result["key"] != "base" {
		t.Errorf("expected key=base, got %v", result["key"])
	}
}

func TestPriorityMerge_NilLayerSkipped(t *testing.T) {
	layers := []merger.PriorityLayer{
		{Priority: 10, Data: nil},
		{Priority: 1, Data: map[string]any{"x": 42}},
	}
	pm := merger.NewPriority(layers)
	result := pm.Merge()

	if result["x"] != 42 {
		t.Errorf("expected x=42, got %v", result["x"])
	}
}

func TestPriorityMerge_AddLayer(t *testing.T) {
	pm := merger.NewPriority(nil)
	pm.AddLayer(1, map[string]any{"a": "low"})
	pm.AddLayer(9, map[string]any{"a": "high"})
	result := pm.Merge()

	if result["a"] != "high" {
		t.Errorf("expected a=high, got %v", result["a"])
	}
}

func TestPriorityMerge_EmptyLayers(t *testing.T) {
	pm := merger.NewPriority(nil)
	result := pm.Merge()
	if len(result) != 0 {
		t.Errorf("expected empty result, got %v", result)
	}
}

func TestPriorityMerge_Layers_ReturnsCopy(t *testing.T) {
	original := []merger.PriorityLayer{
		{Priority: 1, Data: map[string]any{"k": "v"}},
	}
	pm := merger.NewPriority(original)
	copy := pm.Layers()
	copy[0].Priority = 999

	if pm.Layers()[0].Priority == 999 {
		t.Error("Layers() should return a copy, not a reference")
	}
}
