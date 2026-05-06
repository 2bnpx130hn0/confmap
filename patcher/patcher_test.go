package patcher_test

import (
	"testing"

	"github.com/example/confmap/patcher"
)

func baseConfig() map[string]any {
	return map[string]any{
		"app": map[string]any{
			"name": "myapp",
			"port": 8080,
		},
		"debug": false,
	}
}

func TestApply_TopLevelKey(t *testing.T) {
	p := patcher.New(baseConfig())
	if err := p.Apply("debug", true); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Config()["debug"] != true {
		t.Errorf("expected debug=true, got %v", p.Config()["debug"])
	}
}

func TestApply_NestedKey(t *testing.T) {
	p := patcher.New(baseConfig())
	if err := p.Apply("app.port", 9090); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	app := p.Config()["app"].(map[string]any)
	if app["port"] != 9090 {
		t.Errorf("expected port=9090, got %v", app["port"])
	}
}

func TestApply_CreatesIntermediateMaps(t *testing.T) {
	p := patcher.New(baseConfig())
	if err := p.Apply("feature.flags.dark_mode", true); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	feature := p.Config()["feature"].(map[string]any)
	flags := feature["flags"].(map[string]any)
	if flags["dark_mode"] != true {
		t.Errorf("expected dark_mode=true, got %v", flags["dark_mode"])
	}
}

func TestApply_ErrorOnNonMapIntermediate(t *testing.T) {
	cfg := map[string]any{"app": "not-a-map"}
	p := patcher.New(cfg)
	if err := p.Apply("app.port", 8080); err == nil {
		t.Error("expected error when intermediate is not a map")
	}
}

func TestDelete_ExistingKey(t *testing.T) {
	p := patcher.New(baseConfig())
	if err := p.Delete("app.name"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	app := p.Config()["app"].(map[string]any)
	if _, exists := app["name"]; exists {
		t.Error("expected 'name' to be deleted")
	}
}

func TestDelete_MissingIntermediateKey(t *testing.T) {
	p := patcher.New(baseConfig())
	if err := p.Delete("nonexistent.key"); err == nil {
		t.Error("expected error when intermediate key is missing")
	}
}

func TestDelete_TopLevelKey(t *testing.T) {
	p := patcher.New(baseConfig())
	if err := p.Delete("debug"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, exists := p.Config()["debug"]; exists {
		t.Error("expected 'debug' to be deleted")
	}
}

func TestConfig_ReturnsMutatedMap(t *testing.T) {
	cfg := baseConfig()
	p := patcher.New(cfg)
	_ = p.Apply("new_key", "new_value")
	if p.Config()["new_key"] != "new_value" {
		t.Errorf("expected new_key=new_value, got %v", p.Config()["new_key"])
	}
}
