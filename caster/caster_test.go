package caster_test

import (
	"testing"

	"github.com/yourorg/confmap/caster"
)

func TestString_DirectString(t *testing.T) {
	c := caster.New(map[string]interface{}{"name": "alice"})
	v, err := c.String("name")
	if err != nil || v != "alice" {
		t.Fatalf("expected alice, got %q, err: %v", v, err)
	}
}

func TestString_FromInt(t *testing.T) {
	c := caster.New(map[string]interface{}{"port": 8080})
	v, err := c.String("port")
	if err != nil || v != "8080" {
		t.Fatalf("expected 8080, got %q, err: %v", v, err)
	}
}

func TestString_MissingKey(t *testing.T) {
	c := caster.New(map[string]interface{}{})
	_, err := c.String("missing")
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestInt_DirectInt(t *testing.T) {
	c := caster.New(map[string]interface{}{"workers": 4})
	v, err := c.Int("workers")
	if err != nil || v != 4 {
		t.Fatalf("expected 4, got %d, err: %v", v, err)
	}
}

func TestInt_FromFloat64(t *testing.T) {
	c := caster.New(map[string]interface{}{"timeout": float64(30)})
	v, err := c.Int("timeout")
	if err != nil || v != 30 {
		t.Fatalf("expected 30, got %d, err: %v", v, err)
	}
}

func TestInt_FromString(t *testing.T) {
	c := caster.New(map[string]interface{}{"port": "9090"})
	v, err := c.Int("port")
	if err != nil || v != 9090 {
		t.Fatalf("expected 9090, got %d, err: %v", v, err)
	}
}

func TestInt_InvalidString(t *testing.T) {
	c := caster.New(map[string]interface{}{"port": "abc"})
	_, err := c.Int("port")
	if err == nil {
		t.Fatal("expected error for non-numeric string")
	}
}

func TestBool_DirectBool(t *testing.T) {
	c := caster.New(map[string]interface{}{"debug": true})
	v, err := c.Bool("debug")
	if err != nil || !v {
		t.Fatalf("expected true, got %v, err: %v", v, err)
	}
}

func TestBool_FromString(t *testing.T) {
	c := caster.New(map[string]interface{}{"verbose": "false"})
	v, err := c.Bool("verbose")
	if err != nil || v {
		t.Fatalf("expected false, got %v, err: %v", v, err)
	}
}

func TestBool_MissingKey(t *testing.T) {
	c := caster.New(map[string]interface{}{})
	_, err := c.Bool("enabled")
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestBool_UnsupportedType(t *testing.T) {
	c := caster.New(map[string]interface{}{"flag": []string{"a"}})
	_, err := c.Bool("flag")
	if err == nil {
		t.Fatal("expected error for unsupported type")
	}
}
