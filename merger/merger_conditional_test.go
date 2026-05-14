package merger

import (
	"testing"
)

func TestConditionalMerge_AllPass(t *testing.T) {
	layers := []ConditionalLayer{
		{Layer: map[string]any{"env": "prod"}, Condition: Always()},
		{Layer: map[string]any{"debug": false}, Condition: Always()},
	}
	cm := NewConditional(layers)
	result := cm.Merge()
	if result["env"] != "prod" {
		t.Errorf("expected env=prod, got %v", result["env"])
	}
	if result["debug"] != false {
		t.Errorf("expected debug=false, got %v", result["debug"])
	}
}

func TestConditionalMerge_NeverSkipsLayer(t *testing.T) {
	layers := []ConditionalLayer{
		{Layer: map[string]any{"base": "yes"}, Condition: Always()},
		{Layer: map[string]any{"extra": "no"}, Condition: Never()},
	}
	cm := NewConditional(layers)
	result := cm.Merge()
	if _, ok := result["extra"]; ok {
		t.Error("expected extra key to be absent")
	}
	if result["base"] != "yes" {
		t.Errorf("expected base=yes, got %v", result["base"])
	}
}

func TestConditionalMerge_HasKey(t *testing.T) {
	layers := []ConditionalLayer{
		{Layer: map[string]any{"feature_flags": true}, Condition: Always()},
		{Layer: map[string]any{"flag_detail": "verbose"}, Condition: HasKey("feature_flags")},
		{Layer: map[string]any{"should_not": "appear"}, Condition: HasKey("nonexistent")},
	}
	cm := NewConditional(layers)
	result := cm.Merge()
	if result["flag_detail"] != "verbose" {
		t.Errorf("expected flag_detail=verbose, got %v", result["flag_detail"])
	}
	if _, ok := result["should_not"]; ok {
		t.Error("expected should_not to be absent")
	}
}

func TestConditionalMerge_ValueEquals(t *testing.T) {
	layers := []ConditionalLayer{
		{Layer: map[string]any{"env": "staging"}, Condition: Always()},
		{Layer: map[string]any{"log_level": "debug"}, Condition: ValueEquals("env", "staging")},
		{Layer: map[string]any{"log_level": "error"}, Condition: ValueEquals("env", "prod")},
	}
	cm := NewConditional(layers)
	result := cm.Merge()
	if result["log_level"] != "debug" {
		t.Errorf("expected log_level=debug, got %v", result["log_level"])
	}
}

func TestConditionalMerge_NilConditionAlwaysMerges(t *testing.T) {
	layers := []ConditionalLayer{
		{Layer: map[string]any{"key": "value"}, Condition: nil},
	}
	cm := NewConditional(layers)
	result := cm.Merge()
	if result["key"] != "value" {
		t.Errorf("expected key=value, got %v", result["key"])
	}
}

func TestConditionalMerge_EmptyLayers(t *testing.T) {
	cm := NewConditional([]ConditionalLayer{})
	result := cm.Merge()
	if len(result) != 0 {
		t.Errorf("expected empty result, got %v", result)
	}
}

func TestConditionalMerge_LaterLayerOverrides(t *testing.T) {
	layers := []ConditionalLayer{
		{Layer: map[string]any{"timeout": 30}, Condition: Always()},
		{Layer: map[string]any{"timeout": 60}, Condition: Always()},
	}
	cm := NewConditional(layers)
	result := cm.Merge()
	if result["timeout"] != 60 {
		t.Errorf("expected timeout=60, got %v", result["timeout"])
	}
}
