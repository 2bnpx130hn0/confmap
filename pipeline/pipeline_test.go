package pipeline_test

import (
	"errors"
	"testing"

	"github.com/example/confmap/pipeline"
)

func TestRun_EmptyPipeline(t *testing.T) {
	input := map[string]any{"key": "value"}
	pl := pipeline.New("empty")
	out, err := pl.Run(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["key"] != "value" {
		t.Errorf("expected key=value, got %v", out["key"])
	}
}

func TestRun_StagesAppliedInOrder(t *testing.T) {
	order := []int{}
	makeStage := func(n int) pipeline.Stage {
		return func(cfg map[string]any) (map[string]any, error) {
			order = append(order, n)
			return cfg, nil
		}
	}
	pl := pipeline.New("order").Use(makeStage(1), makeStage(2), makeStage(3))
	_, _ = pl.Run(map[string]any{})
	if len(order) != 3 || order[0] != 1 || order[1] != 2 || order[2] != 3 {
		t.Errorf("unexpected stage order: %v", order)
	}
}

func TestRun_StageError_HaltsPipeline(t *testing.T) {
	called := false
	pl := pipeline.New("err").
		Use(func(_ map[string]any) (map[string]any, error) {
			return nil, errors.New("boom")
		}).
		Use(func(cfg map[string]any) (map[string]any, error) {
			called = true
			return cfg, nil
		})
	_, err := pl.Run(map[string]any{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if called {
		t.Error("second stage should not have been called")
	}
}

func TestRun_DoesNotMutateInput(t *testing.T) {
	input := map[string]any{"a": 1}
	pl := pipeline.New("mutate").Use(func(cfg map[string]any) (map[string]any, error) {
		cfg["b"] = 2
		return cfg, nil
	})
	_, _ = pl.Run(input)
	if _, ok := input["b"]; ok {
		t.Error("original input was mutated")
	}
}

func TestSetDefaults_FillsMissing(t *testing.T) {
	pl := pipeline.New("defaults").Use(pipeline.SetDefaults(map[string]any{"timeout": 30, "host": "localhost"}))
	out, err := pl.Run(map[string]any{"host": "prod.example.com"})
	if err != nil {
		t.Fatal(err)
	}
	if out["host"] != "prod.example.com" {
		t.Errorf("existing key overwritten: %v", out["host"])
	}
	if out["timeout"] != 30 {
		t.Errorf("default not applied: %v", out["timeout"])
	}
}

func TestRequireKeys_MissingKey(t *testing.T) {
	pl := pipeline.New("req").Use(pipeline.RequireKeys("host", "port"))
	_, err := pl.Run(map[string]any{"host": "localhost"})
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestFilterKeys_RemovesMatchingPrefixes(t *testing.T) {
	pl := pipeline.New("filter").Use(pipeline.FilterKeys("internal_", "debug_"))
	out, err := pl.Run(map[string]any{
		"host":          "localhost",
		"internal_id":   "x",
		"debug_verbose": true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := out["internal_id"]; ok {
		t.Error("internal_id should have been filtered")
	}
	if _, ok := out["debug_verbose"]; ok {
		t.Error("debug_verbose should have been filtered")
	}
	if out["host"] != "localhost" {
		t.Error("host should be retained")
	}
}

func TestPipeline_NameAndLen(t *testing.T) {
	pl := pipeline.New("meta").Use(pipeline.RequireKeys("x"), pipeline.FilterKeys("y"))
	if pl.Name() != "meta" {
		t.Errorf("unexpected name: %s", pl.Name())
	}
	if pl.Len() != 2 {
		t.Errorf("expected 2 stages, got %d", pl.Len())
	}
}
