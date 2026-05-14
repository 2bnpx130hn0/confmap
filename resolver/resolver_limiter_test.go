package resolver_test

import (
	"strings"
	"testing"

	"github.com/iamthe1whoknocks/confmap/loader"
	"github.com/iamthe1whoknocks/confmap/resolver"
)

type staticLoader struct {
	data map[string]any
}

func (s *staticLoader) Load() (map[string]any, error) {
	return s.data, nil
}

func TestLimitedResolver_WithinLimits(t *testing.T) {
	loaders := []loader.Loader{
		&staticLoader{data: map[string]any{"host": "localhost", "port": "8080"}},
	}
	r := resolver.NewLimited(loaders, nil, 10, 3)
	cfg, err := r.Resolve()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg["host"] != "localhost" {
		t.Errorf("expected host=localhost, got %v", cfg["host"])
	}
}

func TestLimitedResolver_ExceedsKeyCount(t *testing.T) {
	loaders := []loader.Loader{
		&staticLoader{data: map[string]any{
			"a": "1", "b": "2", "c": "3", "d": "4",
		}},
	}
	r := resolver.NewLimited(loaders, nil, 2, 0)
	_, err := r.Resolve()
	if err == nil {
		t.Fatal("expected error for key count exceeded")
	}
	if !strings.Contains(err.Error(), "limit check failed") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestLimitedResolver_ExceedsDepth(t *testing.T) {
	loaders := []loader.Loader{
		&staticLoader{data: map[string]any{
			"a": map[string]any{
				"b": map[string]any{
					"c": "deep",
				},
			},
		}},
	}
	r := resolver.NewLimited(loaders, nil, 0, 2)
	_, err := r.Resolve()
	if err == nil {
		t.Fatal("expected error for depth exceeded")
	}
	if !strings.Contains(err.Error(), "limit check failed") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestLimitedResolver_ZeroLimitsAlwaysPass(t *testing.T) {
	loaders := []loader.Loader{
		&staticLoader{data: map[string]any{
			"x": map[string]any{"y": map[string]any{"z": "v"}},
		}},
	}
	r := resolver.NewLimited(loaders, nil, 0, 0)
	_, err := r.Resolve()
	if err != nil {
		t.Fatalf("unexpected error with zero limits: %v", err)
	}
}
