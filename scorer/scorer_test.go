package scorer_test

import (
	"testing"

	"github.com/your-org/confmap/scorer"
)

func TestEvaluate_FullCoverage(t *testing.T) {
	s := scorer.New([]string{"host", "port", "db.name"})
	cfg := map[string]any{
		"host": "localhost",
		"port": 5432,
		"db": map[string]any{
			"name": "mydb",
		},
	}
	sc, err := s.Evaluate(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sc.Total != 3 {
		t.Errorf("expected Total=3, got %d", sc.Total)
	}
	if sc.Coverage != 1.0 {
		t.Errorf("expected Coverage=1.0, got %f", sc.Coverage)
	}
}

func TestEvaluate_PartialCoverage(t *testing.T) {
	s := scorer.New([]string{"host", "port", "db.name", "db.user"})
	cfg := map[string]any{
		"host": "localhost",
		"port": 5432,
	}
	sc, err := s.Evaluate(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sc.Total != 2 {
		t.Errorf("expected Total=2, got %d", sc.Total)
	}
	if sc.Coverage != 0.5 {
		t.Errorf("expected Coverage=0.5, got %f", sc.Coverage)
	}
}

func TestEvaluate_EmptyValues(t *testing.T) {
	s := scorer.New([]string{"host", "port"})
	cfg := map[string]any{
		"host": "",
		"port": nil,
	}
	sc, err := s.Evaluate(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sc.EmptyKeys != 2 {
		t.Errorf("expected EmptyKeys=2, got %d", sc.EmptyKeys)
	}
}

func TestEvaluate_Depth(t *testing.T) {
	s := scorer.New(nil)
	cfg := map[string]any{
		"a": map[string]any{
			"b": map[string]any{
				"c": "deep",
			},
		},
	}
	sc, err := s.Evaluate(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sc.Depth != 2 {
		t.Errorf("expected Depth=2, got %d", sc.Depth)
	}
}

func TestEvaluate_NilConfig(t *testing.T) {
	s := scorer.New([]string{"host"})
	_, err := s.Evaluate(nil)
	if err == nil {
		t.Error("expected error for nil config, got nil")
	}
}

func TestEvaluate_NoExpectedKeys(t *testing.T) {
	s := scorer.New(nil)
	cfg := map[string]any{"foo": "bar"}
	sc, err := s.Evaluate(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sc.MaxScore != 1 {
		t.Errorf("expected MaxScore=1 (guard), got %d", sc.MaxScore)
	}
}
