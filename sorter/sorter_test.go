package sorter_test

import (
	"strings"
	"testing"

	"github.com/your-org/confmap/sorter"
)

func baseConfig() map[string]any {
	return map[string]any{
		"zebra": "last",
		"apple": "first",
		"mango": map[string]any{
			"z_sub": 99,
			"a_sub": 1,
		},
		"beta": "middle",
	}
}

func TestKeys_SortedOrder(t *testing.T) {
	s := sorter.New(baseConfig())
	keys := s.Keys()
	if !sorter.IsSorted(keys) {
		t.Errorf("expected sorted keys, got %v", keys)
	}
	if keys[0] != "apple" {
		t.Errorf("expected first key to be 'apple', got %s", keys[0])
	}
}

func TestFlatKeys_IncludesNested(t *testing.T) {
	s := sorter.New(baseConfig())
	flat := s.FlatKeys()
	if !sorter.IsSorted(flat) {
		t.Errorf("expected sorted flat keys, got %v", flat)
	}
	found := false
	for _, k := range flat {
		if k == "mango.a_sub" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected 'mango.a_sub' in flat keys, got %v", flat)
	}
}

func TestSortedLines_Format(t *testing.T) {
	s := sorter.New(map[string]any{
		"z": "last",
		"a": "first",
	})
	lines := s.SortedLines()
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "a=") {
		t.Errorf("expected first line to start with 'a=', got %s", lines[0])
	}
	if !strings.HasPrefix(lines[1], "z=") {
		t.Errorf("expected second line to start with 'z=', got %s", lines[1])
	}
}

func TestSorted_ReturnsShallowCopy(t *testing.T) {
	cfg := map[string]any{"x": 1, "y": 2}
	s := sorter.New(cfg)
	copy := s.Sorted()
	copy["z"] = 3
	if _, ok := cfg["z"]; ok {
		t.Error("modifying sorted copy should not affect original")
	}
}

func TestFilterByPrefix_ReturnsMatching(t *testing.T) {
	s := sorter.New(baseConfig())
	result := s.FilterByPrefix("mango")
	if len(result) != 2 {
		t.Fatalf("expected 2 keys with prefix 'mango', got %d: %v", len(result), result)
	}
	for _, k := range result {
		if !strings.HasPrefix(k, "mango") {
			t.Errorf("unexpected key %s does not start with 'mango'", k)
		}
	}
}

func TestIsSorted_TrueAndFalse(t *testing.T) {
	if !sorter.IsSorted([]string{"a", "b", "c"}) {
		t.Error("expected sorted slice to return true")
	}
	if sorter.IsSorted([]string{"c", "a", "b"}) {
		t.Error("expected unsorted slice to return false")
	}
}
