package namespaces_test

import (
	"testing"

	"github.com/your-org/confmap/namespaces"
)

func baseData() map[string]any {
	return map[string]any{
		"database": map[string]any{
			"host": "localhost",
			"port": 5432,
			"credentials": map[string]any{
				"user": "admin",
			},
		},
		"app": map[string]any{
			"debug": true,
		},
	}
}

func TestGet_ExistingKey(t *testing.T) {
	ns := namespaces.New("database", baseData())
	v, ok := ns.Get("host")
	if !ok {
		t.Fatal("expected key to exist")
	}
	if v != "localhost" {
		t.Fatalf("expected localhost, got %v", v)
	}
}

func TestGet_NestedKey(t *testing.T) {
	ns := namespaces.New("database", baseData())
	v, ok := ns.Get("credentials.user")
	if !ok {
		t.Fatal("expected nested key to exist")
	}
	if v != "admin" {
		t.Fatalf("expected admin, got %v", v)
	}
}

func TestGet_MissingKey(t *testing.T) {
	ns := namespaces.New("database", baseData())
	_, ok := ns.Get("missing")
	if ok {
		t.Fatal("expected key to be absent")
	}
}

func TestSet_NewKey(t *testing.T) {
	data := baseData()
	ns := namespaces.New("database", data)
	if err := ns.Set("timeout", 30); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	v, ok := ns.Get("timeout")
	if !ok || v != 30 {
		t.Fatalf("expected 30, got %v", v)
	}
}

func TestSet_CreatesIntermediateMaps(t *testing.T) {
	data := baseData()
	ns := namespaces.New("database", data)
	if err := ns.Set("pool.max", 10); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	v, ok := ns.Get("pool.max")
	if !ok || v != 10 {
		t.Fatalf("expected 10, got %v", v)
	}
}

func TestSet_ErrorOnNonMapIntermediate(t *testing.T) {
	data := baseData()
	ns := namespaces.New("database", data)
	// "host" is a string, not a map
	if err := ns.Set("host.sub", "x"); err == nil {
		t.Fatal("expected error when intermediate is not a map")
	}
}

func TestKeys_ReturnsLeaves(t *testing.T) {
	ns := namespaces.New("database", baseData())
	keys := ns.Keys()
	if len(keys) == 0 {
		t.Fatal("expected at least one key")
	}
	found := false
	for _, k := range keys {
		if k == "host" {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected 'host' in keys, got %v", keys)
	}
}

func TestEmptyPrefix_UsesRootMap(t *testing.T) {
	data := map[string]any{"foo": "bar"}
	ns := namespaces.New("", data)
	v, ok := ns.Get("foo")
	if !ok || v != "bar" {
		t.Fatalf("expected bar, got %v", v)
	}
}
