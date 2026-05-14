package differ_test

import (
	"testing"

	"github.com/your-org/confmap/differ"
)

func TestDiff_Added(t *testing.T) {
	d := differ.New()
	oldCfg := map[string]interface{}{"host": "localhost"}
	newCfg := map[string]interface{}{"host": "localhost", "port": 8080}

	deltas := d.Diff(oldCfg, newCfg)
	if len(deltas) != 1 {
		t.Fatalf("expected 1 delta, got %d", len(deltas))
	}
	if deltas[0].Type != differ.Added || deltas[0].Key != "port" {
		t.Errorf("unexpected delta: %+v", deltas[0])
	}
}

func TestDiff_Removed(t *testing.T) {
	d := differ.New()
	oldCfg := map[string]interface{}{"host": "localhost", "debug": true}
	newCfg := map[string]interface{}{"host": "localhost"}

	deltas := d.Diff(oldCfg, newCfg)
	if len(deltas) != 1 {
		t.Fatalf("expected 1 delta, got %d", len(deltas))
	}
	if deltas[0].Type != differ.Removed || deltas[0].Key != "debug" {
		t.Errorf("unexpected delta: %+v", deltas[0])
	}
}

func TestDiff_Changed(t *testing.T) {
	d := differ.New()
	oldCfg := map[string]interface{}{"host": "localhost"}
	newCfg := map[string]interface{}{"host": "prod.example.com"}

	deltas := d.Diff(oldCfg, newCfg)
	if len(deltas) != 1 {
		t.Fatalf("expected 1 delta, got %d", len(deltas))
	}
	if deltas[0].Type != differ.Changed {
		t.Errorf("expected Changed, got %s", deltas[0].Type)
	}
	if deltas[0].OldValue != "localhost" || deltas[0].NewValue != "prod.example.com" {
		t.Errorf("unexpected values: %+v", deltas[0])
	}
}

func TestDiff_Nested(t *testing.T) {
	d := differ.New()
	oldCfg := map[string]interface{}{
		"db": map[string]interface{}{"host": "localhost", "port": 5432},
	}
	newCfg := map[string]interface{}{
		"db": map[string]interface{}{"host": "db.prod", "port": 5432},
	}

	deltas := d.Diff(oldCfg, newCfg)
	if len(deltas) != 1 {
		t.Fatalf("expected 1 delta, got %d", len(deltas))
	}
	if deltas[0].Key != "db.host" {
		t.Errorf("expected key db.host, got %s", deltas[0].Key)
	}
}

func TestDiff_NoDifference(t *testing.T) {
	d := differ.New()
	cfg := map[string]interface{}{"key": "value"}
	deltas := d.Diff(cfg, cfg)
	if len(deltas) != 0 {
		t.Errorf("expected no deltas, got %d", len(deltas))
	}
}

func TestDiff_BothEmpty(t *testing.T) {
	d := differ.New()
	deltas := d.Diff(map[string]interface{}{}, map[string]interface{}{})
	if len(deltas) != 0 {
		t.Errorf("expected no deltas, got %d", len(deltas))
	}
}

func TestDiff_MultipleChanges(t *testing.T) {
	d := differ.New()
	oldCfg := map[string]interface{}{"host": "localhost", "port": 8080, "debug": true}
	newCfg := map[string]interface{}{"host": "prod.example.com", "port": 9090}

	deltas := d.Diff(oldCfg, newCfg)
	if len(deltas) != 3 {
		t.Fatalf("expected 3 deltas, got %d", len(deltas))
	}

	deltaByKey := make(map[string]differ.Delta)
	for _, delta := range deltas {
		deltaByKey[delta.Key] = delta
	}

	if deltaByKey["host"].Type != differ.Changed {
		t.Errorf("expected host to be Changed, got %s", deltaByKey["host"].Type)
	}
	if deltaByKey["port"].Type != differ.Changed {
		t.Errorf("expected port to be Changed, got %s", deltaByKey["port"].Type)
	}
	if deltaByKey["debug"].Type != differ.Removed {
		t.Errorf("expected debug to be Removed, got %s", deltaByKey["debug"].Type)
	}
}
