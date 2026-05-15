package resolver_test

import (
	"errors"
	"testing"

	"github.com/your-org/confmap/resolver"
)

func TestMatchedResolver_FiltersKeys(t *testing.T) {
	loaders := []resolver.LoaderFunc{
		func() (map[string]any, error) {
			return map[string]any{
				"db": map[string]any{
					"host": "localhost",
					"port": 5432,
				},
				"app": map[string]any{
					"name": "myapp",
					"port": 8080,
				},
			}, nil
		},
	}

	mr, err := resolver.NewMatched(loaders, nil, "db.*")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	cfg, err := mr.Resolve()
	if err != nil {
		t.Fatalf("resolve error: %v", err)
	}

	db, ok := cfg["db"].(map[string]any)
	if !ok {
		t.Fatal("expected db map in result")
	}
	if db["host"] != "localhost" {
		t.Errorf("expected db.host=localhost, got %v", db["host"])
	}
	if _, ok := cfg["app"]; ok {
		t.Error("expected app to be filtered out")
	}
}

func TestMatchedResolver_LoaderError(t *testing.T) {
	loaders := []resolver.LoaderFunc{
		func() (map[string]any, error) {
			return nil, errors.New("load failure")
		},
	}

	mr, err := resolver.NewMatched(loaders, nil, "*")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = mr.Resolve()
	if err == nil {
		t.Fatal("expected error from loader")
	}
}

func TestMatchedResolver_NoPatternsError(t *testing.T) {
	_, err := resolver.NewMatched(nil, nil)
	if err == nil {
		t.Fatal("expected error when no patterns provided")
	}
}

func TestMatchedResolver_MultiplePatterns(t *testing.T) {
	loaders := []resolver.LoaderFunc{
		func() (map[string]any, error) {
			return map[string]any{
				"db":    map[string]any{"host": "db-host"},
				"cache": map[string]any{"ttl": 300},
				"log":   map[string]any{"level": "info"},
			}, nil
		},
	}

	mr, err := resolver.NewMatched(loaders, nil, "db.*", "cache.*")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	cfg, err := mr.Resolve()
	if err != nil {
		t.Fatalf("resolve error: %v", err)
	}

	if _, ok := cfg["db"]; !ok {
		t.Error("expected db in result")
	}
	if _, ok := cfg["cache"]; !ok {
		t.Error("expected cache in result")
	}
	if _, ok := cfg["log"]; ok {
		t.Error("expected log to be filtered out")
	}
}
