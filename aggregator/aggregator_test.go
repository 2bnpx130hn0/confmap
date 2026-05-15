package aggregator_test

import (
	"errors"
	"testing"

	"github.com/nicholasgasior/confmap/aggregator"
)

func TestReduce_MergeReduce_Override(t *testing.T) {
	agg := aggregator.New(aggregator.MergeReduce)
	agg.Add(map[string]any{"host": "localhost", "port": 5432})
	agg.Add(map[string]any{"port": 9999})

	result, err := agg.Reduce()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["host"] != "localhost" {
		t.Errorf("expected host=localhost, got %v", result["host"])
	}
	if result["port"] != 9999 {
		t.Errorf("expected port=9999, got %v", result["port"])
	}
}

func TestReduce_CollectKeys(t *testing.T) {
	agg := aggregator.New(aggregator.CollectKeys)
	agg.Add(map[string]any{"alpha": 1, "beta": 2})
	agg.Add(map[string]any{"beta": 3, "gamma": 4})

	result, err := agg.Reduce()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, k := range []string{"alpha", "beta", "gamma"} {
		if result[k] != true {
			t.Errorf("expected key %q to be collected", k)
		}
	}
}

func TestReduce_NilLayerSkipped(t *testing.T) {
	agg := aggregator.New(aggregator.MergeReduce)
	agg.Add(map[string]any{"a": 1})
	agg.Add(nil)
	agg.Add(map[string]any{"b": 2})

	if agg.Count() != 2 {
		t.Errorf("expected 2 layers, got %d", agg.Count())
	}
	result, err := agg.Reduce()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["a"] != 1 || result["b"] != 2 {
		t.Errorf("unexpected result: %v", result)
	}
}

func TestReduce_EmptyLayers(t *testing.T) {
	agg := aggregator.New(aggregator.MergeReduce)
	result, err := agg.Reduce()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected empty result, got %v", result)
	}
}

func TestReduce_ReduceFnError(t *testing.T) {
	failing := func(acc, layer map[string]any) (map[string]any, error) {
		return nil, errors.New("boom")
	}
	agg := aggregator.New(failing)
	agg.Add(map[string]any{"x": 1})

	_, err := agg.Reduce()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestReduce_NilReduceFn(t *testing.T) {
	agg := aggregator.New(nil)
	_, err := agg.Reduce()
	if err == nil {
		t.Fatal("expected error for nil reduce function")
	}
}

func TestCount_ReflectsAddedLayers(t *testing.T) {
	agg := aggregator.New(aggregator.MergeReduce)
	if agg.Count() != 0 {
		t.Errorf("expected 0, got %d", agg.Count())
	}
	agg.Add(map[string]any{"k": "v"})
	if agg.Count() != 1 {
		t.Errorf("expected 1, got %d", agg.Count())
	}
}
