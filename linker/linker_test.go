package linker_test

import (
	"testing"

	"github.com/igorvisi/confmap/linker"
)

func TestResolve_SimpleReference(t *testing.T) {
	cfg := map[string]any{
		"host":    "localhost",
		"address": "${host}:8080",
	}
	l := linker.New(cfg)
	out, err := l.Resolve()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["address"] != "localhost:8080" {
		t.Errorf("expected localhost:8080, got %v", out["address"])
	}
}

func TestResolve_ChainedReferences(t *testing.T) {
	cfg := map[string]any{
		"scheme":  "https",
		"host":    "example.com",
		"base":    "${scheme}://${host}",
		"api":     "${base}/v1",
	}
	l := linker.New(cfg)
	out, err := l.Resolve()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["api"] != "https://example.com/v1" {
		t.Errorf("expected https://example.com/v1, got %v", out["api"])
	}
}

func TestResolve_MissingReference(t *testing.T) {
	cfg := map[string]any{
		"url": "${missing_key}/path",
	}
	l := linker.New(cfg)
	_, err := l.Resolve()
	if err == nil {
		t.Fatal("expected error for missing reference, got nil")
	}
}

func TestResolve_CycleDetection(t *testing.T) {
	cfg := map[string]any{
		"a": "${b}",
		"b": "${a}",
	}
	l := linker.New(cfg)
	_, err := l.Resolve()
	if err == nil {
		t.Fatal("expected cycle error, got nil")
	}
}

func TestResolve_NoReferences(t *testing.T) {
	cfg := map[string]any{
		"port": "3000",
		"host": "127.0.0.1",
	}
	l := linker.New(cfg)
	out, err := l.Resolve()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["port"] != "3000" {
		t.Errorf("expected 3000, got %v", out["port"])
	}
}

func TestResolve_NestedReference(t *testing.T) {
	cfg := map[string]any{
		"db": map[string]any{
			"host": "db-host",
			"dsn":  "postgres://${db.host}/mydb",
		},
	}
	l := linker.New(cfg)
	out, err := l.Resolve()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	db, ok := out["db"].(map[string]any)
	if !ok {
		t.Fatal("expected nested map for 'db'")
	}
	if db["dsn"] != "postgres://db-host/mydb" {
		t.Errorf("expected postgres://db-host/mydb, got %v", db["dsn"])
	}
}
