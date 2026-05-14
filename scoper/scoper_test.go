package scoper_test

import (
	"testing"

	"github.com/user/confmap/scoper"
)

func baseData() map[string]any {
	return map[string]any{
		"app": map[string]any{
			"name":    "confmap",
			"version": "1.0",
		},
		"db": map[string]any{
			"host": "localhost",
			"port": 5432,
		},
		"flat": "value",
	}
}

func TestNew_ValidPrefix(t *testing.T) {
	sc, err := scoper.New("app", baseData())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sc.Prefix() != "app" {
		t.Errorf("expected prefix 'app', got %q", sc.Prefix())
	}
}

func TestNew_MissingPrefix(t *testing.T) {
	_, err := scoper.New("missing", baseData())
	if err == nil {
		t.Fatal("expected error for missing prefix")
	}
}

func TestNew_NonMapPrefix(t *testing.T) {
	_, err := scoper.New("flat", baseData())
	if err == nil {
		t.Fatal("expected error when prefix value is not a map")
	}
}

func TestNew_EmptyPrefix(t *testing.T) {
	sc, err := scoper.New("", baseData())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, ok := sc.Get("app")
	if !ok {
		t.Error("expected to find 'app' key at root scope")
	}
}

func TestGet_ExistingKey(t *testing.T) {
	sc, _ := scoper.New("db", baseData())
	v, ok := sc.Get("host")
	if !ok {
		t.Fatal("expected key 'host' to exist")
	}
	if v != "localhost" {
		t.Errorf("expected 'localhost', got %v", v)
	}
}

func TestGet_MissingKey(t *testing.T) {
	sc, _ := scoper.New("db", baseData())
	_, ok := sc.Get("password")
	if ok {
		t.Error("expected missing key to return false")
	}
}

func TestSet_AddsKey(t *testing.T) {
	sc, _ := scoper.New("db", baseData())
	sc.Set("password", "secret")
	v, ok := sc.Get("password")
	if !ok || v != "secret" {
		t.Errorf("expected 'secret', got %v (ok=%v)", v, ok)
	}
}

func TestKeys_ContainsExpected(t *testing.T) {
	sc, _ := scoper.New("app", baseData())
	keys := sc.Keys()
	if len(keys) != 2 {
		t.Errorf("expected 2 keys, got %d", len(keys))
	}
}

func TestSnapshot_IsShallowCopy(t *testing.T) {
	sc, _ := scoper.New("db", baseData())
	snap := sc.Snapshot()
	snap["host"] = "changed"
	v, _ := sc.Get("host")
	if v == "changed" {
		t.Error("snapshot mutation should not affect scoper data")
	}
}
