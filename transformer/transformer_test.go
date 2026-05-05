package transformer_test

import (
	"errors"
	"testing"

	"github.com/yourorg/confmap/transformer"
)

func TestSetDefault_KeyAbsent(t *testing.T) {
	data := map[string]any{"host": "localhost"}
	fn := transformer.SetDefault("port", 8080)
	if err := fn(data); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if data["port"] != 8080 {
		t.Errorf("expected port=8080, got %v", data["port"])
	}
}

func TestSetDefault_KeyPresent(t *testing.T) {
	data := map[string]any{"port": 9090}
	fn := transformer.SetDefault("port", 8080)
	_ = fn(data)
	if data["port"] != 9090 {
		t.Errorf("expected port unchanged at 9090, got %v", data["port"])
	}
}

func TestRename(t *testing.T) {
	data := map[string]any{"db_host": "127.0.0.1"}
	fn := transformer.Rename("db_host", "database.host")
	if err := fn(data); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := data["db_host"]; ok {
		t.Error("old key should have been removed")
	}
	if data["database.host"] != "127.0.0.1" {
		t.Errorf("expected database.host=127.0.0.1, got %v", data["database.host"])
	}
}

func TestRename_MissingKey(t *testing.T) {
	data := map[string]any{"other": "value"}
	fn := transformer.Rename("db_host", "database.host")
	if err := fn(data); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(data) != 1 {
		t.Errorf("data should be unchanged, got %v", data)
	}
}

func TestCoerceString(t *testing.T) {
	data := map[string]any{"port": 3000}
	fn := transformer.CoerceString("port")
	if err := fn(data); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if data["port"] != "3000" {
		t.Errorf("expected port=\"3000\", got %v", data["port"])
	}
}

func TestApply_Chain(t *testing.T) {
	data := map[string]any{"level": 2}
	tr := transformer.New(
		transformer.SetDefault("debug", false),
		transformer.CoerceString("level"),
	)
	if err := tr.Apply(data); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if data["debug"] != false {
		t.Errorf("expected debug=false, got %v", data["debug"])
	}
	if data["level"] != "2" {
		t.Errorf("expected level=\"2\", got %v", data["level"])
	}
}

func TestApply_ErrorPropagates(t *testing.T) {
	sentinel := errors.New("boom")
	tr := transformer.New(func(_ map[string]any) error { return sentinel })
	err := tr.Apply(map[string]any{})
	if !errors.Is(err, sentinel) {
		t.Errorf("expected sentinel error, got %v", err)
	}
}
