package coercer_test

import (
	"testing"

	"github.com/user/confmap/coercer"
)

func TestToString_FromInt(t *testing.T) {
	cfg := map[string]interface{}{"port": 8080}
	c := coercer.New().ToString("port")
	if err := c.Apply(cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg["port"] != "8080" {
		t.Errorf("expected \"8080\", got %v", cfg["port"])
	}
}

func TestToString_FromBool(t *testing.T) {
	cfg := map[string]interface{}{"debug": true}
	c := coercer.New().ToString("debug")
	if err := c.Apply(cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg["debug"] != "true" {
		t.Errorf("expected \"true\", got %v", cfg["debug"])
	}
}

func TestToInt_FromString(t *testing.T) {
	cfg := map[string]interface{}{"timeout": "30"}
	c := coercer.New().ToInt("timeout")
	if err := c.Apply(cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg["timeout"] != 30 {
		t.Errorf("expected 30, got %v", cfg["timeout"])
	}
}

func TestToInt_FromFloat64(t *testing.T) {
	cfg := map[string]interface{}{"workers": float64(4)}
	c := coercer.New().ToInt("workers")
	if err := c.Apply(cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg["workers"] != 4 {
		t.Errorf("expected 4, got %v", cfg["workers"])
	}
}

func TestToBool_FromString(t *testing.T) {
	cfg := map[string]interface{}{"enabled": "true"}
	c := coercer.New().ToBool("enabled")
	if err := c.Apply(cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg["enabled"] != true {
		t.Errorf("expected true, got %v", cfg["enabled"])
	}
}

func TestToBool_FromInt(t *testing.T) {
	cfg := map[string]interface{}{"verbose": 1}
	c := coercer.New().ToBool("verbose")
	if err := c.Apply(cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg["verbose"] != true {
		t.Errorf("expected true, got %v", cfg["verbose"])
	}
}

func TestApply_MissingKeyIsSkipped(t *testing.T) {
	cfg := map[string]interface{}{"host": "localhost"}
	c := coercer.New().ToInt("port")
	if err := c.Apply(cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := cfg["port"]; ok {
		t.Error("expected missing key to remain absent")
	}
}

func TestApply_InvalidCoercion_ReturnsError(t *testing.T) {
	cfg := map[string]interface{}{"port": "not-a-number"}
	c := coercer.New().ToInt("port")
	if err := c.Apply(cfg); err == nil {
		t.Error("expected error for invalid int coercion, got nil")
	}
}

func TestApply_MultipleRules(t *testing.T) {
	cfg := map[string]interface{}{
		"port":    "9090",
		"debug":   "false",
		"version": 2,
	}
	c := coercer.New().ToInt("port").ToBool("debug").ToString("version")
	if err := c.Apply(cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg["port"] != 9090 {
		t.Errorf("port: expected 9090, got %v", cfg["port"])
	}
	if cfg["debug"] != false {
		t.Errorf("debug: expected false, got %v", cfg["debug"])
	}
	if cfg["version"] != "2" {
		t.Errorf("version: expected \"2\", got %v", cfg["version"])
	}
}
