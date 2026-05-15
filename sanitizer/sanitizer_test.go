package sanitizer_test

import (
	"testing"

	"github.com/yourusername/confmap/sanitizer"
)

func TestApply_TrimSpace(t *testing.T) {
	s := sanitizer.New(sanitizer.TrimSpace)
	cfg := map[string]any{"key": "  hello  "}
	out := s.Apply(cfg)
	if out["key"] != "hello" {
		t.Errorf("expected 'hello', got %q", out["key"])
	}
}

func TestApply_ToLower(t *testing.T) {
	s := sanitizer.New(sanitizer.ToLower)
	cfg := map[string]any{"env": "PRODUCTION"}
	out := s.Apply(cfg)
	if out["env"] != "production" {
		t.Errorf("expected 'production', got %q", out["env"])
	}
}

func TestApply_ChainedRules(t *testing.T) {
	s := sanitizer.New(sanitizer.TrimSpace, sanitizer.ToUpper)
	cfg := map[string]any{"mode": "  debug  "}
	out := s.Apply(cfg)
	if out["mode"] != "DEBUG" {
		t.Errorf("expected 'DEBUG', got %q", out["mode"])
	}
}

func TestApply_NestedMap(t *testing.T) {
	s := sanitizer.New(sanitizer.TrimSpace)
	cfg := map[string]any{
		"db": map[string]any{
			"host": "  localhost  ",
		},
	}
	out := s.Apply(cfg)
	db, ok := out["db"].(map[string]any)
	if !ok {
		t.Fatal("expected nested map")
	}
	if db["host"] != "localhost" {
		t.Errorf("expected 'localhost', got %q", db["host"])
	}
}

func TestApply_SliceValues(t *testing.T) {
	s := sanitizer.New(sanitizer.ToLower)
	cfg := map[string]any{"tags": []any{"Alpha", "BETA"}}
	out := s.Apply(cfg)
	tags, ok := out["tags"].([]any)
	if !ok || len(tags) != 2 {
		t.Fatal("expected slice of length 2")
	}
	if tags[0] != "alpha" || tags[1] != "beta" {
		t.Errorf("unexpected tags: %v", tags)
	}
}

func TestApply_StripNull(t *testing.T) {
	s := sanitizer.New(sanitizer.StripNull)
	cfg := map[string]any{"val": "null", "other": "NULL"}
	out := s.Apply(cfg)
	if out["val"] != "" || out["other"] != "" {
		t.Errorf("expected empty strings, got %v", out)
	}
}

func TestApply_ReplaceRule(t *testing.T) {
	s := sanitizer.New(sanitizer.ReplaceRule("-", "_"))
	cfg := map[string]any{"key": "my-config-key"}
	out := s.Apply(cfg)
	if out["key"] != "my_config_key" {
		t.Errorf("expected 'my_config_key', got %q", out["key"])
	}
}

func TestApply_NonStringUnchanged(t *testing.T) {
	s := sanitizer.New(sanitizer.ToLower)
	cfg := map[string]any{"count": 42, "enabled": true}
	out := s.Apply(cfg)
	if out["count"] != 42 || out["enabled"] != true {
		t.Errorf("non-string values should be unchanged: %v", out)
	}
}

func TestApply_NilConfig(t *testing.T) {
	s := sanitizer.New(sanitizer.TrimSpace)
	out := s.Apply(nil)
	if out != nil {
		t.Errorf("expected nil, got %v", out)
	}
}

func TestApply_DoesNotMutateOriginal(t *testing.T) {
	s := sanitizer.New(sanitizer.ToUpper)
	cfg := map[string]any{"name": "alice"}
	_ = s.Apply(cfg)
	if cfg["name"] != "alice" {
		t.Error("original config was mutated")
	}
}
