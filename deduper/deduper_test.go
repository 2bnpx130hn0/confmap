package deduper_test

import (
	"testing"

	"github.com/nicholasgasior/confmap/deduper"
)

func baseConfig() map[string]any {
	return map[string]any{
		"host": "localhost",
		"port": 8080,
		"db": map[string]any{
			"host": "dbhost",
			"port": 5432,
		},
	}
}

func TestDedupeKeys_NilConfig(t *testing.T) {
	d := deduper.New()
	_, err := d.DedupeKeys(nil)
	if err == nil {
		t.Fatal("expected error for nil config")
	}
}

func TestDedupeKeys_RetainsAllUniqueKeys(t *testing.T) {
	d := deduper.New()
	cfg := baseConfig()
	result, err := d.DedupeKeys(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["host"] != "localhost" {
		t.Errorf("expected host=localhost, got %v", result["host"])
	}
	if result["port"] != 8080 {
		t.Errorf("expected port=8080, got %v", result["port"])
	}
	db, ok := result["db"].(map[string]any)
	if !ok {
		t.Fatal("expected db to be a map")
	}
	if db["host"] != "dbhost" {
		t.Errorf("expected db.host=dbhost, got %v", db["host"])
	}
}

func TestDedupeKeys_DoesNotMutateOriginal(t *testing.T) {
	d := deduper.New()
	cfg := baseConfig()
	_, err := d.DedupeKeys(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg["host"] != "localhost" {
		t.Error("original config was mutated")
	}
}

func TestDedupeValues_NilConfig(t *testing.T) {
	d := deduper.New()
	_, err := d.DedupeValues(nil)
	if err == nil {
		t.Fatal("expected error for nil config")
	}
}

func TestDedupeValues_RemovesDuplicateValues(t *testing.T) {
	d := deduper.New()
	cfg := map[string]any{
		"a": "same",
		"b": "same",
		"c": "different",
	}
	result, err := d.DedupeValues(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := result["a"]; ok {
		t.Error("expected key 'a' to be removed as duplicate value")
	}
	if result["c"] != "different" {
		t.Errorf("expected c=different, got %v", result["c"])
	}
}

func TestDedupeValues_NestedMapHandled(t *testing.T) {
	d := deduper.New()
	cfg := map[string]any{
		"db": map[string]any{
			"x": "val",
			"y": "val",
		},
	}
	result, err := d.DedupeValues(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	db, ok := result["db"].(map[string]any)
	if !ok {
		t.Fatal("expected db to be a map")
	}
	if len(db) != 1 {
		t.Errorf("expected 1 key after dedup, got %d", len(db))
	}
}

func TestCount_FlatConfig(t *testing.T) {
	d := deduper.New()
	cfg := map[string]any{"a": 1, "b": 2, "c": 3}
	if got := d.Count(cfg); got != 3 {
		t.Errorf("expected count=3, got %d", got)
	}
}

func TestCount_NestedConfig(t *testing.T) {
	d := deduper.New()
	cfg := baseConfig()
	// host, port, db.host, db.port = 4 leaf keys
	if got := d.Count(cfg); got != 4 {
		t.Errorf("expected count=4, got %d", got)
	}
}

func TestCount_NilConfig(t *testing.T) {
	d := deduper.New()
	if got := d.Count(nil); got != 0 {
		t.Errorf("expected count=0 for nil, got %d", got)
	}
}
