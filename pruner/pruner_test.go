package pruner_test

import (
	"testing"

	"github.com/your-org/confmap/pruner"
)

func TestApply_RemovesNilValues(t *testing.T) {
	p := pruner.New(pruner.PruneNil)
	cfg := map[string]any{
		"host": "localhost",
		"port": nil,
	}
	out, err := p.Apply(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["port"]; ok {
		t.Error("expected 'port' to be pruned")
	}
	if out["host"] != "localhost" {
		t.Error("expected 'host' to be retained")
	}
}

func TestApply_RemovesEmptyStrings(t *testing.T) {
	p := pruner.New(pruner.PruneEmpty)
	cfg := map[string]any{
		"name": "",
		"env":  "prod",
	}
	out, err := p.Apply(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["name"]; ok {
		t.Error("expected 'name' to be pruned")
	}
	if out["env"] != "prod" {
		t.Error("expected 'env' to be retained")
	}
}

func TestApply_RemovesZeroValues(t *testing.T) {
	p := pruner.New(pruner.PruneZero)
	cfg := map[string]any{
		"retries": 0,
		"timeout": 30,
		"debug":   false,
	}
	out, err := p.Apply(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["retries"]; ok {
		t.Error("expected 'retries' to be pruned")
	}
	if _, ok := out["debug"]; ok {
		t.Error("expected 'debug' to be pruned")
	}
	if out["timeout"] != 30 {
		t.Error("expected 'timeout' to be retained")
	}
}

func TestApply_NestedMap(t *testing.T) {
	p := pruner.New(pruner.PruneNil | pruner.PruneEmpty)
	cfg := map[string]any{
		"db": map[string]any{
			"host":     "localhost",
			"password": nil,
			"prefix":   "",
		},
	}
	out, err := p.Apply(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	db, ok := out["db"].(map[string]any)
	if !ok {
		t.Fatal("expected 'db' to be a map")
	}
	if _, ok := db["password"]; ok {
		t.Error("expected 'password' to be pruned")
	}
	if _, ok := db["prefix"]; ok {
		t.Error("expected 'prefix' to be pruned")
	}
	if db["host"] != "localhost" {
		t.Error("expected 'host' to be retained")
	}
}

func TestApply_EmptyNestedMapPruned(t *testing.T) {
	p := pruner.New(pruner.PruneNil | pruner.PruneEmpty)
	cfg := map[string]any{
		"meta": map[string]any{
			"tag": nil,
		},
	}
	out, err := p.Apply(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["meta"]; ok {
		t.Error("expected empty nested map 'meta' to be pruned")
	}
}

func TestApply_NilConfig_ReturnsError(t *testing.T) {
	p := pruner.New(0)
	_, err := p.Apply(nil)
	if err == nil {
		t.Error("expected error for nil config")
	}
}

func TestMustApply_Panics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for nil config")
		}
	}()
	pruner.New(0).MustApply(nil)
}
