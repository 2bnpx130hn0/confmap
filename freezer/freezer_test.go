package freezer_test

import (
	"errors"
	"testing"

	"github.com/iamkevla/confmap/freezer"
)

func TestNew_StoresData(t *testing.T) {
	cfg := map[string]any{"host": "localhost", "port": 8080}
	f := freezer.New(cfg)

	v, ok := f.Get("host")
	if !ok || v != "localhost" {
		t.Fatalf("expected host=localhost, got %v", v)
	}
}

func TestSet_BeforeFreeze(t *testing.T) {
	f := freezer.New(map[string]any{})
	if err := f.Set("key", "value"); err != nil {
		t.Fatalf("unexpected error before freeze: %v", err)
	}
	v, ok := f.Get("key")
	if !ok || v != "value" {
		t.Fatalf("expected key=value, got %v", v)
	}
}

func TestSet_AfterFreeze_ReturnsError(t *testing.T) {
	f := freezer.New(map[string]any{"a": 1})
	f.Freeze()

	err := f.Set("a", 2)
	if err == nil {
		t.Fatal("expected error after freeze, got nil")
	}
	if !errors.Is(err, freezer.ErrFrozen) {
		t.Fatalf("expected ErrFrozen, got %v", err)
	}
}

func TestIsFrozen(t *testing.T) {
	f := freezer.New(map[string]any{})
	if f.IsFrozen() {
		t.Fatal("should not be frozen initially")
	}
	f.Freeze()
	if !f.IsFrozen() {
		t.Fatal("should be frozen after Freeze()")
	}
}

func TestSnapshot_IsDeepCopy(t *testing.T) {
	orig := map[string]any{"db": map[string]any{"host": "127.0.0.1"}}
	f := freezer.New(orig)

	snap := f.Snapshot()
	// Mutate the snapshot; original inside freezer should be unaffected.
	snap["db"].(map[string]any)["host"] = "mutated"

	v, _ := f.Get("db")
	dbMap, ok := v.(map[string]any)
	if !ok || dbMap["host"] != "127.0.0.1" {
		t.Fatalf("snapshot mutation leaked into freezer: %v", dbMap["host"])
	}
}

func TestNew_OriginalMutationDoesNotAffectFreezer(t *testing.T) {
	cfg := map[string]any{"env": "dev"}
	f := freezer.New(cfg)
	cfg["env"] = "prod" // mutate original

	v, _ := f.Get("env")
	if v != "dev" {
		t.Fatalf("expected dev, got %v", v)
	}
}
