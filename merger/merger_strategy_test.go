package merger

import (
	"testing"
)

func TestStrategicMerge_Override(t *testing.T) {
	sm := NewStrategic(StrategyOverride)
	base := map[string]any{"a": 1, "b": 2}
	overlay := map[string]any{"b": 99, "c": 3}

	result := sm.Merge(base, overlay)

	if result["a"] != 1 {
		t.Errorf("expected a=1, got %v", result["a"])
	}
	if result["b"] != 99 {
		t.Errorf("expected b=99 (overridden), got %v", result["b"])
	}
	if result["c"] != 3 {
		t.Errorf("expected c=3, got %v", result["c"])
	}
}

func TestStrategicMerge_KeepBase(t *testing.T) {
	sm := NewStrategic(StrategyKeepBase)
	base := map[string]any{"a": 1, "b": 2}
	overlay := map[string]any{"b": 99, "c": 3}

	result := sm.Merge(base, overlay)

	if result["b"] != 2 {
		t.Errorf("expected b=2 (base kept), got %v", result["b"])
	}
	if result["c"] != 3 {
		t.Errorf("expected c=3 (new key added), got %v", result["c"])
	}
}

func TestStrategicMerge_AppendSlice(t *testing.T) {
	sm := NewStrategic(StrategyAppendSlice)
	base := map[string]any{"tags": []any{"a", "b"}}
	overlay := map[string]any{"tags": []any{"c", "d"}}

	result := sm.Merge(base, overlay)

	tags, ok := result["tags"].([]any)
	if !ok {
		t.Fatal("expected tags to be []any")
	}
	if len(tags) != 4 {
		t.Errorf("expected 4 tags, got %d", len(tags))
	}
}

func TestStrategicMerge_AppendSlice_FallbackOnNonSlice(t *testing.T) {
	sm := NewStrategic(StrategyAppendSlice)
	base := map[string]any{"x": "hello"}
	overlay := map[string]any{"x": "world"}

	result := sm.Merge(base, overlay)

	if result["x"] != "world" {
		t.Errorf("expected x=world (override fallback), got %v", result["x"])
	}
}

func TestStrategicMerge_DoesNotMutateBase(t *testing.T) {
	sm := NewStrategic(StrategyOverride)
	base := map[string]any{"a": 1}
	overlay := map[string]any{"a": 2}

	_ = sm.Merge(base, overlay)

	if base["a"] != 1 {
		t.Errorf("base was mutated: expected a=1, got %v", base["a"])
	}
}

func TestStrategicMerge_NilOverlay(t *testing.T) {
	sm := NewStrategic(StrategyOverride)
	base := map[string]any{"a": 1}

	result := sm.Merge(base, nil)

	if result["a"] != 1 {
		t.Errorf("expected a=1, got %v", result["a"])
	}
}
