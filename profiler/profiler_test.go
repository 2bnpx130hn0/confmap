package profiler_test

import (
	"strings"
	"testing"

	"github.com/yourorg/confmap/profiler"
)

func TestAnalyze_FlatConfig(t *testing.T) {
	cfg := map[string]any{
		"host": "localhost",
		"port": 8080,
		"debug": true,
	}
	p := profiler.New()
	prof := p.Analyze(cfg)

	if prof.TotalKeys != 3 {
		t.Errorf("expected 3 total keys, got %d", prof.TotalKeys)
	}
	if prof.MaxDepth != 1 {
		t.Errorf("expected max depth 1, got %d", prof.MaxDepth)
	}
	if prof.TypeCounts["string"] != 1 {
		t.Errorf("expected 1 string, got %d", prof.TypeCounts["string"])
	}
	if prof.TypeCounts["int"] != 1 {
		t.Errorf("expected 1 int, got %d", prof.TypeCounts["int"])
	}
	if prof.TypeCounts["bool"] != 1 {
		t.Errorf("expected 1 bool, got %d", prof.TypeCounts["bool"])
	}
}

func TestAnalyze_NestedConfig(t *testing.T) {
	cfg := map[string]any{
		"database": map[string]any{
			"host": "db.local",
			"port": 5432,
		},
		"app": "myapp",
	}
	p := profiler.New()
	prof := p.Analyze(cfg)

	if prof.TotalKeys != 4 {
		t.Errorf("expected 4 total keys, got %d", prof.TotalKeys)
	}
	if prof.MaxDepth != 2 {
		t.Errorf("expected max depth 2, got %d", prof.MaxDepth)
	}
	if prof.TypeCounts["map"] != 1 {
		t.Errorf("expected 1 map, got %d", prof.TypeCounts["map"])
	}
}

func TestAnalyze_EmptyConfig(t *testing.T) {
	p := profiler.New()
	prof := p.Analyze(map[string]any{})

	if prof.TotalKeys != 0 {
		t.Errorf("expected 0 keys, got %d", prof.TotalKeys)
	}
	if prof.MaxDepth != 0 {
		t.Errorf("expected max depth 0, got %d", prof.MaxDepth)
	}
}

func TestSummary_ContainsExpectedFields(t *testing.T) {
	cfg := map[string]any{
		"name": "test",
		"value": 42,
	}
	p := profiler.New()
	prof := p.Analyze(cfg)
	summary := profiler.Summary(prof)

	if !strings.Contains(summary, "TotalKeys: 2") {
		t.Errorf("summary missing TotalKeys: %s", summary)
	}
	if !strings.Contains(summary, "MaxDepth: 1") {
		t.Errorf("summary missing MaxDepth: %s", summary)
	}
}
