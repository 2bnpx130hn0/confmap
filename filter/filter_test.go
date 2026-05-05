package filter_test

import (
	"testing"

	"github.com/yourorg/confmap/filter"
)

func baseConfig() map[string]any {
	return map[string]any{
		"host":    "localhost",
		"port":    8080,
		"debug":   true,
		"timeout": 30,
	}
}

func TestApply_NoRules(t *testing.T) {
	cfg := baseConfig()
	out := filter.New().Apply(cfg)
	if len(out) != len(cfg) {
		t.Fatalf("expected %d keys, got %d", len(cfg), len(out))
	}
}

func TestApply_IncludeOnly(t *testing.T) {
	out := filter.New().Include("host", "port").Apply(baseConfig())
	if len(out) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(out))
	}
	if _, ok := out["host"]; !ok {
		t.Error("expected 'host' to be present")
	}
	if _, ok := out["debug"]; ok {
		t.Error("expected 'debug' to be absent")
	}
}

func TestApply_ExcludeOnly(t *testing.T) {
	out := filter.New().Exclude("debug", "timeout").Apply(baseConfig())
	if len(out) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(out))
	}
	if _, ok := out["debug"]; ok {
		t.Error("expected 'debug' to be absent")
	}
}

func TestApply_IncludeAndExclude(t *testing.T) {
	// include host+port, but then exclude port — only host should remain
	out := filter.New().Include("host", "port").Exclude("port").Apply(baseConfig())
	if len(out) != 1 {
		t.Fatalf("expected 1 key, got %d", len(out))
	}
	if _, ok := out["host"]; !ok {
		t.Error("expected 'host' to be present")
	}
}

func TestApply_EmptyConfig(t *testing.T) {
	out := filter.New().Include("host").Apply(map[string]any{})
	if len(out) != 0 {
		t.Fatalf("expected empty map, got %d keys", len(out))
	}
}

func TestKeys_Sorted(t *testing.T) {
	cfg := baseConfig()
	keys := filter.Keys(cfg)
	expected := []string{"debug", "host", "port", "timeout"}
	if len(keys) != len(expected) {
		t.Fatalf("expected %d keys, got %d", len(expected), len(keys))
	}
	for i, k := range keys {
		if k != expected[i] {
			t.Errorf("position %d: expected %q, got %q", i, expected[i], k)
		}
	}
}
