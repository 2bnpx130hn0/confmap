package merger

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestWeightedMerge_HigherWeightWins(t *testing.T) {
	wm := NewWeighted().
		Add(10, map[string]any{"host": "base", "port": 5432}).
		Add(20, map[string]any{"host": "override"})

	got := wm.Merge()
	want := map[string]any{"host": "override", "port": 5432}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("unexpected result (-want +got):\n%s", diff)
	}
}

func TestWeightedMerge_LowerWeightDoesNotOverride(t *testing.T) {
	wm := NewWeighted().
		Add(50, map[string]any{"debug": true}).
		Add(5, map[string]any{"debug": false})

	got := wm.Merge()
	if got["debug"] != true {
		t.Fatalf("expected debug=true, got %v", got["debug"])
	}
}

func TestWeightedMerge_EqualWeightLastAddedWins(t *testing.T) {
	wm := NewWeighted().
		Add(10, map[string]any{"key": "first"}).
		Add(10, map[string]any{"key": "second"})

	got := wm.Merge()
	if got["key"] != "second" {
		t.Fatalf("expected 'second', got %v", got["key"])
	}
}

func TestWeightedMerge_NilLayerSkipped(t *testing.T) {
	wm := NewWeighted().
		Add(1, nil).
		Add(2, map[string]any{"a": 1})

	got := wm.Merge()
	if got["a"] != 1 {
		t.Fatalf("expected a=1, got %v", got["a"])
	}
}

func TestWeightedMerge_EmptyMerger(t *testing.T) {
	wm := NewWeighted()
	got := wm.Merge()
	if len(got) != 0 {
		t.Fatalf("expected empty map, got %v", got)
	}
}

func TestWeightedMerge_Layers_ReturnsCopy(t *testing.T) {
	wm := NewWeighted().
		Add(1, map[string]any{"x": 1}).
		Add(2, map[string]any{"y": 2})

	layers := wm.Layers()
	if len(layers) != 2 {
		t.Fatalf("expected 2 layers, got %d", len(layers))
	}
	// Mutating the returned slice must not affect the merger.
	layers[0].Weight = 999
	if wm.Layers()[0].Weight == 999 {
		t.Fatal("Layers() returned a reference to internal state")
	}
}
