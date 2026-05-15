package merger

import (
	"testing"
)

func TestVersionedMerge_AscendingOrder(t *testing.T) {
	vm := NewVersioned()
	vm.AddLayer("1.0.0", map[string]any{"host": "old", "port": 8080})
	vm.AddLayer("2.0.0", map[string]any{"host": "new"})

	got, err := vm.Merge()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["host"] != "new" {
		t.Errorf("expected host=new, got %v", got["host"])
	}
	if got["port"] != 8080 {
		t.Errorf("expected port=8080, got %v", got["port"])
	}
}

func TestVersionedMerge_InsertionOrderDoesNotMatter(t *testing.T) {
	vm := NewVersioned()
	// Add higher version first — lower version must not override it.
	vm.AddLayer("3.0.0", map[string]any{"key": "v3"})
	vm.AddLayer("1.0.0", map[string]any{"key": "v1"})
	vm.AddLayer("2.0.0", map[string]any{"key": "v2"})

	got, err := vm.Merge()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["key"] != "v3" {
		t.Errorf("expected key=v3, got %v", got["key"])
	}
}

func TestVersionedMerge_NilLayerSkipped(t *testing.T) {
	vm := NewVersioned()
	vm.AddLayer("1.0.0", nil)
	vm.AddLayer("2.0.0", map[string]any{"alive": true})

	got, err := vm.Merge()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["alive"] != true {
		t.Errorf("expected alive=true, got %v", got["alive"])
	}
}

func TestVersionedMerge_EmptyLayers(t *testing.T) {
	vm := NewVersioned()
	got, err := vm.Merge()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("expected empty map, got %v", got)
	}
}

func TestVersionedMerge_EmptyVersionError(t *testing.T) {
	vm := NewVersioned()
	vm.AddLayer("", map[string]any{"x": 1})

	_, err := vm.Merge()
	if err == nil {
		t.Fatal("expected error for empty version string, got nil")
	}
}

func TestVersionedMerge_SingleLayer(t *testing.T) {
	vm := NewVersioned()
	vm.AddLayer("1.2.3", map[string]any{"only": "value"})

	got, err := vm.Merge()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["only"] != "value" {
		t.Errorf("expected only=value, got %v", got["only"])
	}
}
