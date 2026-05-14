package merger

import (
	"testing"
)

func TestMerge_Override(t *testing.T) {
	m := New(StrategyOverride)
	base := map[string]any{
		"host": "localhost",
		"port": 5432,
		"db": map[string]any{
			"name": "mydb",
			"pool": 5,
		},
	}
	override := map[string]any{
		"host": "prod.example.com",
		"db": map[string]any{
			"pool": 20,
		},
	}
	result, err := m.Merge(base, override)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["host"] != "prod.example.com" {
		t.Errorf("expected host=prod.example.com, got %v", result["host"])
	}
	if result["port"] != 5432 {
		t.Errorf("expected port=5432, got %v", result["port"])
	}
	db, ok := result["db"].(map[string]any)
	if !ok {
		t.Fatal("expected db to be a map")
	}
	if db["name"] != "mydb" {
		t.Errorf("expected db.name=mydb, got %v", db["name"])
	}
	if db["pool"] != 20 {
		t.Errorf("expected db.pool=20, got %v", db["pool"])
	}
}

func TestMerge_KeepExisting(t *testing.T) {
	m := New(StrategyKeepExisting)
	base := map[string]any{"timeout": 30}
	override := map[string]any{"timeout": 60, "retries": 3}
	result, err := m.Merge(base, override)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["timeout"] != 30 {
		t.Errorf("expected timeout=30, got %v", result["timeout"])
	}
	if result["retries"] != 3 {
		t.Errorf("expected retries=3, got %v", result["retries"])
	}
}

func TestMerge_NilLayer(t *testing.T) {
	m := New(StrategyOverride)
	_, err := m.Merge(map[string]any{"a": 1}, nil)
	if err == nil {
		t.Error("expected error for nil layer")
	}
}

func TestMerge_EmptyLayers(t *testing.T) {
	m := New(StrategyOverride)
	result, err := m.Merge()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected empty result, got %v", result)
	}
}

func TestMerge_SingleLayer(t *testing.T) {
	m := New(StrategyOverride)
	input := map[string]any{"key": "value", "num": 42}
	result, err := m.Merge(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["key"] != "value" {
		t.Errorf("expected key=value, got %v", result["key"])
	}
	if result["num"] != 42 {
		t.Errorf("expected num=42, got %v", result["num"])
	}
}
