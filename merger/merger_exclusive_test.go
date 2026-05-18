package merger

import (
	"testing"
)

func TestExclusiveMerge_UniqueKeysRetained(t *testing.T) {
	m := NewExclusive()
	m.AddLayer(map[string]any{"a": 1, "shared": "x"})
	m.AddLayer(map[string]any{"b": 2, "shared": "y"})

	result := m.Merge()

	if _, ok := result["shared"]; ok {
		t.Error("expected 'shared' to be excluded but it was present")
	}
	if result["a"] != 1 {
		t.Errorf("expected a=1, got %v", result["a"])
	}
	if result["b"] != 2 {
		t.Errorf("expected b=2, got %v", result["b"])
	}
}

func TestExclusiveMerge_AllKeysShared(t *testing.T) {
	m := NewExclusive()
	m.AddLayer(map[string]any{"k": 1})
	m.AddLayer(map[string]any{"k": 2})

	result := m.Merge()

	if len(result) != 0 {
		t.Errorf("expected empty result, got %v", result)
	}
}

func TestExclusiveMerge_NilLayerSkipped(t *testing.T) {
	m := NewExclusive()
	m.AddLayer(nil)
	m.AddLayer(map[string]any{"only": true})

	result := m.Merge()

	if result["only"] != true {
		t.Errorf("expected only=true, got %v", result["only"])
	}
}

func TestExclusiveMerge_EmptyLayers(t *testing.T) {
	m := NewExclusive()

	result := m.Merge()

	if len(result) != 0 {
		t.Errorf("expected empty result, got %v", result)
	}
}

func TestExclusiveMerge_SingleLayer(t *testing.T) {
	m := NewExclusive()
	m.AddLayer(map[string]any{"x": 10, "y": 20})

	result := m.Merge()

	if result["x"] != 10 || result["y"] != 20 {
		t.Errorf("expected x=10 y=20, got %v", result)
	}
}

func TestExclusiveMerge_ThreeLayersPartialOverlap(t *testing.T) {
	m := NewExclusive()
	m.AddLayer(map[string]any{"a": 1, "overlap": "first"})
	m.AddLayer(map[string]any{"b": 2, "overlap": "second"})
	m.AddLayer(map[string]any{"c": 3})

	result := m.Merge()

	if _, ok := result["overlap"]; ok {
		t.Error("expected 'overlap' to be excluded")
	}
	for _, key := range []string{"a", "b", "c"} {
		if _, ok := result[key]; !ok {
			t.Errorf("expected key %q to be present", key)
		}
	}
}
