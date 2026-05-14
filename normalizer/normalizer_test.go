package normalizer_test

import (
	"sort"
	"testing"

	"github.com/your-org/confmap/normalizer"
)

func TestApply_Lowercase(t *testing.T) {
	n := normalizer.New(normalizer.Options{Lowercase: true})
	input := map[string]any{"HOST": "localhost", "PORT": 8080}
	out := n.Apply(input)
	if _, ok := out["host"]; !ok {
		t.Error("expected key 'host'")
	}
	if _, ok := out["port"]; !ok {
		t.Error("expected key 'port'")
	}
}

func TestApply_TrimSpace(t *testing.T) {
	n := normalizer.New(normalizer.Options{TrimSpace: true})
	input := map[string]any{"  key  ": "value"}
	out := n.Apply(input)
	if _, ok := out["key"]; !ok {
		t.Error("expected trimmed key 'key'")
	}
}

func TestApply_SeparatorReplacement(t *testing.T) {
	n := normalizer.New(normalizer.Options{OldSep: "-", NewSep: "_"})
	input := map[string]any{"my-key": "val", "another-key": 1}
	out := n.Apply(input)
	if _, ok := out["my_key"]; !ok {
		t.Error("expected key 'my_key'")
	}
	if _, ok := out["another_key"]; !ok {
		t.Error("expected key 'another_key'")
	}
}

func TestApply_Nested(t *testing.T) {
	n := normalizer.New(normalizer.Options{Lowercase: true})
	input := map[string]any{
		"DATABASE": map[string]any{
			"HOST": "db.local",
			"PORT": 5432,
		},
	}
	out := n.Apply(input)
	db, ok := out["database"].(map[string]any)
	if !ok {
		t.Fatal("expected nested map under 'database'")
	}
	if db["host"] != "db.local" {
		t.Errorf("expected 'host' = 'db.local', got %v", db["host"])
	}
}

func TestApply_DoesNotMutateOriginal(t *testing.T) {
	n := normalizer.New(normalizer.Options{Lowercase: true})
	input := map[string]any{"KEY": "value"}
	_ = n.Apply(input)
	if _, ok := input["KEY"]; !ok {
		t.Error("original map should not be mutated")
	}
}

func TestApply_SliceValues(t *testing.T) {
	n := normalizer.New(normalizer.Options{Lowercase: true})
	input := map[string]any{
		"ITEMS": []any{
			map[string]any{"NAME": "alpha"},
			map[string]any{"NAME": "beta"},
		},
	}
	out := n.Apply(input)
	items, ok := out["items"].([]any)
	if !ok || len(items) != 2 {
		t.Fatal("expected slice with 2 items under 'items'")
	}
	first, ok := items[0].(map[string]any)
	if !ok {
		t.Fatal("expected map as first slice element")
	}
	if first["name"] != "alpha" {
		t.Errorf("expected 'name' = 'alpha', got %v", first["name"])
	}
}

func TestKeys_ReturnNormalizedKeys(t *testing.T) {
	n := normalizer.New(normalizer.Options{Lowercase: true})
	input := map[string]any{"HOST": "x", "PORT": 1, "DEBUG": true}
	keys := n.Keys(input)
	sort.Strings(keys)
	expected := []string{"debug", "host", "port"}
	for i, k := range keys {
		if k != expected[i] {
			t.Errorf("key[%d]: expected %q, got %q", i, expected[i], k)
		}
	}
}
