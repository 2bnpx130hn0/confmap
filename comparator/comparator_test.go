package comparator_test

import (
	"testing"

	"github.com/iamando/confmap/comparator"
)

func TestCompare_EqualMaps(t *testing.T) {
	cmp := comparator.New()
	a := map[string]any{"host": "localhost", "port": 8080}
	b := map[string]any{"host": "localhost", "port": 8080}
	res := cmp.Compare(a, b)
	if !res.Equal {
		t.Errorf("expected maps to be equal")
	}
	if res.Similarity != 1.0 {
		t.Errorf("expected similarity 1.0, got %f", res.Similarity)
	}
}

func TestCompare_OnlyInA(t *testing.T) {
	cmp := comparator.New()
	a := map[string]any{"host": "localhost", "debug": true}
	b := map[string]any{"host": "localhost"}
	res := cmp.Compare(a, b)
	if res.Equal {
		t.Error("expected maps to differ")
	}
	if len(res.OnlyInA) != 1 || res.OnlyInA[0] != "debug" {
		t.Errorf("expected OnlyInA=[debug], got %v", res.OnlyInA)
	}
}

func TestCompare_OnlyInB(t *testing.T) {
	cmp := comparator.New()
	a := map[string]any{"host": "localhost"}
	b := map[string]any{"host": "localhost", "timeout": 30}
	res := cmp.Compare(a, b)
	if len(res.OnlyInB) != 1 || res.OnlyInB[0] != "timeout" {
		t.Errorf("expected OnlyInB=[timeout], got %v", res.OnlyInB)
	}
}

func TestCompare_DifferingValues(t *testing.T) {
	cmp := comparator.New()
	a := map[string]any{"port": 8080}
	b := map[string]any{"port": 9090}
	res := cmp.Compare(a, b)
	if res.Equal {
		t.Error("expected maps to differ")
	}
	if len(res.Differing) != 1 || res.Differing[0] != "port" {
		t.Errorf("expected Differing=[port], got %v", res.Differing)
	}
}

func TestCompare_NestedMaps(t *testing.T) {
	cmp := comparator.New()
	a := map[string]any{"db": map[string]any{"host": "localhost", "port": 5432}}
	b := map[string]any{"db": map[string]any{"host": "remotehost", "port": 5432}}
	res := cmp.Compare(a, b)
	if res.Equal {
		t.Error("expected nested maps to differ")
	}
	if len(res.Differing) != 1 || res.Differing[0] != "db.host" {
		t.Errorf("expected Differing=[db.host], got %v", res.Differing)
	}
}

func TestCompare_IgnoreKeys(t *testing.T) {
	cmp := comparator.New("updated_at")
	a := map[string]any{"name": "app", "updated_at": "2024-01-01"}
	b := map[string]any{"name": "app", "updated_at": "2024-06-01"}
	res := cmp.Compare(a, b)
	if !res.Equal {
		t.Errorf("expected maps equal after ignoring updated_at, differing: %v", res.Differing)
	}
}

func TestCompare_EmptyMaps(t *testing.T) {
	cmp := comparator.New()
	res := cmp.Compare(map[string]any{}, map[string]any{})
	if !res.Equal {
		t.Error("expected empty maps to be equal")
	}
}
