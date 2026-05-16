package interpolator_test

import (
	"os"
	"testing"

	"github.com/confmap/interpolator"
)

func TestApply_SimpleExpansion(t *testing.T) {
	ip := interpolator.New(map[string]string{"HOST": "localhost", "PORT": "9090"}, true)
	out, err := ip.Apply(map[string]any{"addr": "${HOST}:${PORT}"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["addr"] != "localhost:9090" {
		t.Errorf("expected localhost:9090, got %v", out["addr"])
	}
}

func TestApply_MissingVar_StrictError(t *testing.T) {
	ip := interpolator.New(map[string]string{}, true)
	_, err := ip.Apply(map[string]any{"key": "${MISSING}"})
	if err == nil {
		t.Fatal("expected error for missing variable in strict mode")
	}
}

func TestApply_MissingVar_NonStrictEmpty(t *testing.T) {
	ip := interpolator.New(map[string]string{}, false)
	out, err := ip.Apply(map[string]any{"key": "${MISSING}"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["key"] != "" {
		t.Errorf("expected empty string, got %v", out["key"])
	}
}

func TestApply_NestedMap(t *testing.T) {
	ip := interpolator.New(map[string]string{"ENV": "prod"}, true)
	out, err := ip.Apply(map[string]any{
		"db": map[string]any{"name": "mydb_${ENV}"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	db, ok := out["db"].(map[string]any)
	if !ok {
		t.Fatal("expected nested map")
	}
	if db["name"] != "mydb_prod" {
		t.Errorf("expected mydb_prod, got %v", db["name"])
	}
}

func TestApply_SliceValues(t *testing.T) {
	ip := interpolator.New(map[string]string{"X": "42"}, true)
	out, err := ip.Apply(map[string]any{"items": []any{"val=${X}", "plain"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	items, ok := out["items"].([]any)
	if !ok || len(items) != 2 {
		t.Fatal("expected slice of 2")
	}
	if items[0] != "val=42" {
		t.Errorf("expected val=42, got %v", items[0])
	}
}

func TestApply_NonStringPassthrough(t *testing.T) {
	ip := interpolator.New(map[string]string{}, false)
	out, err := ip.Apply(map[string]any{"count": 7, "flag": true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["count"] != 7 || out["flag"] != true {
		t.Error("non-string values should pass through unchanged")
	}
}

func TestNewFromEnv_ExpandsProcessEnv(t *testing.T) {
	t.Setenv("CONFMAP_TEST_VAR", "hello")
	_ = os.Getenv("CONFMAP_TEST_VAR") // ensure set
	ip := interpolator.NewFromEnv(true)
	out, err := ip.Apply(map[string]any{"msg": "${CONFMAP_TEST_VAR}"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["msg"] != "hello" {
		t.Errorf("expected hello, got %v", out["msg"])
	}
}
