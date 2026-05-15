package aliaser_test

import (
	"testing"

	"github.com/iamBijoyKar/confmap/aliaser"
)

func TestRegister_ValidAlias(t *testing.T) {
	a := aliaser.New()
	if err := a.Register("db_host", "database.host"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	aliases := a.Aliases()
	if aliases["db_host"] != "database.host" {
		t.Errorf("expected alias to be registered")
	}
}

func TestRegister_EmptyAlias(t *testing.T) {
	a := aliaser.New()
	if err := a.Register("", "canonical"); err == nil {
		t.Error("expected error for empty alias")
	}
}

func TestRegister_SameAliasAndCanonical(t *testing.T) {
	a := aliaser.New()
	if err := a.Register("key", "key"); err == nil {
		t.Error("expected error when alias equals canonical")
	}
}

func TestResolve_KnownAlias(t *testing.T) {
	a := aliaser.New()
	_ = a.Register("host", "server.host")
	if got := a.Resolve("host"); got != "server.host" {
		t.Errorf("expected server.host, got %s", got)
	}
}

func TestResolve_UnknownKey(t *testing.T) {
	a := aliaser.New()
	if got := a.Resolve("unknown"); got != "unknown" {
		t.Errorf("expected key unchanged, got %s", got)
	}
}

func TestApply_ExpandsAlias(t *testing.T) {
	a := aliaser.New()
	_ = a.Register("db_host", "database.host")
	cfg := map[string]any{"db_host": "localhost", "port": 5432}
	out, err := a.Apply(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["database.host"] != "localhost" {
		t.Errorf("expected database.host=localhost")
	}
	if _, ok := out["db_host"]; ok {
		t.Error("alias key should have been removed")
	}
	if out["port"] != 5432 {
		t.Error("non-alias key should be preserved")
	}
}

func TestApply_CanonicalTakesPrecedence(t *testing.T) {
	a := aliaser.New()
	_ = a.Register("host", "server.host")
	cfg := map[string]any{"host": "alias-value", "server.host": "canon-value"}
	out, err := a.Apply(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["server.host"] != "canon-value" {
		t.Errorf("canonical should take precedence, got %v", out["server.host"])
	}
}

func TestApply_NilConfig(t *testing.T) {
	a := aliaser.New()
	out, err := a.Apply(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != nil {
		t.Error("expected nil output for nil input")
	}
}

func TestApply_NoAliasRegistered(t *testing.T) {
	a := aliaser.New()
	cfg := map[string]any{"key": "value"}
	out, err := a.Apply(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["key"] != "value" {
		t.Error("config should be unchanged when no aliases registered")
	}
}
