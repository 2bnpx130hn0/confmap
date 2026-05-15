package squasher_test

import (
	"testing"

	"github.com/iamando/confmap/squasher"
)

func TestSquash_SingleLayer(t *testing.T) {
	s := squasher.New(".")
	layer := map[string]any{"host": "localhost", "port": 5432}
	out, err := s.Squash(layer)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["host"] != "localhost" || out["port"] != 5432 {
		t.Fatalf("unexpected output: %v", out)
	}
}

func TestSquash_LaterLayerWins(t *testing.T) {
	s := squasher.New(".")
	base := map[string]any{"debug": false, "timeout": 30}
	override := map[string]any{"debug": true}
	out, err := s.Squash(base, override)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["debug"] != true {
		t.Errorf("expected debug=true, got %v", out["debug"])
	}
	if out["timeout"] != 30 {
		t.Errorf("expected timeout=30, got %v", out["timeout"])
	}
}

func TestSquash_NestedFlattened(t *testing.T) {
	s := squasher.New(".")
	layer := map[string]any{
		"db": map[string]any{
			"host": "127.0.0.1",
			"port": 5432,
		},
	}
	out, err := s.Squash(layer)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["db.host"] != "127.0.0.1" {
		t.Errorf("expected db.host=127.0.0.1, got %v", out["db.host"])
	}
	if out["db.port"] != 5432 {
		t.Errorf("expected db.port=5432, got %v", out["db.port"])
	}
}

func TestSquash_NilLayerSkipped(t *testing.T) {
	s := squasher.New(".")
	out, err := s.Squash(nil, map[string]any{"key": "val"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["key"] != "val" {
		t.Errorf("expected key=val, got %v", out["key"])
	}
}

func TestSquash_EmptyLayers(t *testing.T) {
	s := squasher.New(".")
	out, err := s.Squash()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 0 {
		t.Errorf("expected empty map, got %v", out)
	}
}

func TestExpand_RoundTrip(t *testing.T) {
	s := squasher.New(".")
	original := map[string]any{
		"server": map[string]any{
			"host": "example.com",
			"tls":  true,
		},
		"timeout": 10,
	}
	flat, err := s.Squash(original)
	if err != nil {
		t.Fatalf("squash error: %v", err)
	}
	nested, err := s.Expand(flat)
	if err != nil {
		t.Fatalf("expand error: %v", err)
	}
	server, ok := nested["server"].(map[string]any)
	if !ok {
		t.Fatalf("expected server to be a map")
	}
	if server["host"] != "example.com" {
		t.Errorf("expected server.host=example.com, got %v", server["host"])
	}
	if nested["timeout"] != 10 {
		t.Errorf("expected timeout=10, got %v", nested["timeout"])
	}
}

func TestExpand_ConflictError(t *testing.T) {
	s := squasher.New(".")
	flat := map[string]any{
		"a":   "scalar",
		"a.b": "nested",
	}
	_, err := s.Expand(flat)
	if err == nil {
		t.Fatal("expected error on key segment collision, got nil")
	}
}

func TestNew_DefaultSeparator(t *testing.T) {
	s := squasher.New("")
	layer := map[string]any{"x": map[string]any{"y": 1}}
	out, err := s.Squash(layer)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["x.y"]; !ok {
		t.Errorf("expected default separator '.', keys: %v", out)
	}
}
