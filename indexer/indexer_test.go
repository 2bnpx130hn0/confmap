package indexer_test

import (
	"sort"
	"testing"

	"github.com/example/confmap/indexer"
)

func baseConfig() map[string]any {
	return map[string]any{
		"app": map[string]any{
			"name": "confmap",
			"version": "1.0",
		},
		"database": map[string]any{
			"host": "localhost",
			"port": 5432,
		},
		"debug": true,
	}
}

func TestGet_FlatKey(t *testing.T) {
	idx := indexer.New(baseConfig())
	v, ok := idx.Get("debug")
	if !ok {
		t.Fatal("expected key 'debug' to exist")
	}
	if v != true {
		t.Fatalf("expected true, got %v", v)
	}
}

func TestGet_NestedKey(t *testing.T) {
	idx := indexer.New(baseConfig())
	v, ok := idx.Get("database.host")
	if !ok {
		t.Fatal("expected key 'database.host' to exist")
	}
	if v != "localhost" {
		t.Fatalf("expected 'localhost', got %v", v)
	}
}

func TestGet_MissingKey(t *testing.T) {
	idx := indexer.New(baseConfig())
	_, ok := idx.Get("nonexistent.key")
	if ok {
		t.Fatal("expected key to be absent")
	}
}

func TestKeys_ContainsAllLeaves(t *testing.T) {
	idx := indexer.New(baseConfig())
	keys := idx.Keys()
	sort.Strings(keys)

	expected := []string{"app.name", "app.version", "database.host", "database.port", "debug"}
	if len(keys) != len(expected) {
		t.Fatalf("expected %d keys, got %d: %v", len(expected), len(keys), keys)
	}
	for i, k := range expected {
		if keys[i] != k {
			t.Errorf("expected key[%d]=%q, got %q", i, k, keys[i])
		}
	}
}

func TestMustGet_ExistingKey(t *testing.T) {
	idx := indexer.New(baseConfig())
	v := idx.MustGet("app.name")
	if v != "confmap" {
		t.Fatalf("expected 'confmap', got %v", v)
	}
}

func TestMustGet_MissingKey_Panics(t *testing.T) {
	idx := indexer.New(baseConfig())
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for missing key")
		}
	}()
	idx.MustGet("does.not.exist")
}

func TestRebuild_UpdatesIndex(t *testing.T) {
	idx := indexer.New(baseConfig())
	newCfg := map[string]any{"service": map[string]any{"port": 8080}}
	idx.Rebuild(newCfg)

	if _, ok := idx.Get("database.host"); ok {
		t.Error("old key should be gone after rebuild")
	}
	v, ok := idx.Get("service.port")
	if !ok {
		t.Fatal("expected 'service.port' after rebuild")
	}
	if v != 8080 {
		t.Fatalf("expected 8080, got %v", v)
	}
}
