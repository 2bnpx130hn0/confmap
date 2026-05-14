package limiter_test

import (
	"strings"
	"testing"

	"github.com/iamthe1whoknocks/confmap/limiter"
)

func TestCheck_WithinLimits(t *testing.T) {
	l := limiter.New(10, 3)
	cfg := map[string]any{
		"a": "1",
		"b": map[string]any{"c": "2"},
	}
	if err := l.Check(cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCheck_ExceedsKeyCount(t *testing.T) {
	l := limiter.New(2, 0)
	cfg := map[string]any{
		"a": "1",
		"b": "2",
		"c": "3",
	}
	err := l.Check(cfg)
	if err == nil {
		t.Fatal("expected error for exceeded key count")
	}
	if !strings.Contains(err.Error(), "key count") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestCheck_ExceedsDepth(t *testing.T) {
	l := limiter.New(0, 2)
	cfg := map[string]any{
		"a": map[string]any{
			"b": map[string]any{
				"c": "deep",
			},
		},
	}
	err := l.Check(cfg)
	if err == nil {
		t.Fatal("expected error for exceeded depth")
	}
	if !strings.Contains(err.Error(), "nesting depth") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestEnforce_ReturnsConfig(t *testing.T) {
	l := limiter.New(5, 2)
	cfg := map[string]any{"x": "y"}
	out, err := l.Enforce(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["x"] != "y" {
		t.Errorf("expected config to be returned unchanged")
	}
}

func TestEnforce_ReturnsNilOnViolation(t *testing.T) {
	l := limiter.New(1, 0)
	cfg := map[string]any{"a": "1", "b": "2"}
	out, err := l.Enforce(cfg)
	if err == nil {
		t.Fatal("expected error")
	}
	if out != nil {
		t.Errorf("expected nil config on error")
	}
}

func TestStats_CountsCorrectly(t *testing.T) {
	l := limiter.New(0, 0)
	cfg := map[string]any{
		"a": "1",
		"b": map[string]any{
			"c": "2",
			"d": "3",
		},
	}
	keys, depth := l.Stats(cfg)
	if keys != 4 {
		t.Errorf("expected 4 keys, got %d", keys)
	}
	if depth != 2 {
		t.Errorf("expected depth 2, got %d", depth)
	}
}

func TestStats_EmptyConfig(t *testing.T) {
	l := limiter.New(0, 0)
	keys, depth := l.Stats(map[string]any{})
	if keys != 0 || depth != 0 {
		t.Errorf("expected 0 keys and 0 depth, got %d/%d", keys, depth)
	}
}

func TestCheck_ZeroLimits_AlwaysPasses(t *testing.T) {
	l := limiter.New(0, 0)
	cfg := map[string]any{
		"a": map[string]any{"b": map[string]any{"c": map[string]any{"d": "v"}}},
	}
	if err := l.Check(cfg); err != nil {
		t.Fatalf("unexpected error with zero limits: %v", err)
	}
}
