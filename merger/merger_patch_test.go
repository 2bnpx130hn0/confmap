package merger

import (
	"testing"
)

func TestPatch_SetNewKey(t *testing.T) {
	base := map[string]interface{}{"host": "localhost"}
	p := NewPatch([]PatchOp{{Op: "set", Path: "port", Value: 8080}})
	out, err := p.Merge(base)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["port"] != 8080 {
		t.Errorf("expected port=8080, got %v", out["port"])
	}
	if out["host"] != "localhost" {
		t.Errorf("expected host unchanged, got %v", out["host"])
	}
}

func TestPatch_OverrideExistingKey(t *testing.T) {
	base := map[string]interface{}{"host": "localhost", "port": 80}
	p := NewPatch([]PatchOp{{Op: "set", Path: "port", Value: 443}})
	out, err := p.Merge(base)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["port"] != 443 {
		t.Errorf("expected port=443, got %v", out["port"])
	}
}

func TestPatch_DeleteKey(t *testing.T) {
	base := map[string]interface{}{"host": "localhost", "debug": true}
	p := NewPatch([]PatchOp{{Op: "delete", Path: "debug"}})
	out, err := p.Merge(base)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["debug"]; ok {
		t.Error("expected 'debug' key to be deleted")
	}
}

func TestPatch_RenameKey(t *testing.T) {
	base := map[string]interface{}{"db_host": "127.0.0.1"}
	p := NewPatch([]PatchOp{{Op: "rename", Path: "db_host", To: "database_host"}})
	out, err := p.Merge(base)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["database_host"] != "127.0.0.1" {
		t.Errorf("expected database_host=127.0.0.1, got %v", out["database_host"])
	}
	if _, ok := out["db_host"]; ok {
		t.Error("expected old key 'db_host' to be removed after rename")
	}
}

func TestPatch_RenameKeyMissingSource(t *testing.T) {
	base := map[string]interface{}{"host": "localhost"}
	p := NewPatch([]PatchOp{{Op: "rename", Path: "missing", To: "new_key"}})
	_, err := p.Merge(base)
	if err == nil {
		t.Fatal("expected error for missing rename source key")
	}
}

func TestPatch_UnknownOp(t *testing.T) {
	base := map[string]interface{}{"host": "localhost"}
	p := NewPatch([]PatchOp{{Op: "upsert", Path: "host", Value: "remote"}})
	_, err := p.Merge(base)
	if err == nil {
		t.Fatal("expected error for unknown op")
	}
}

func TestPatch_DoesNotMutateBase(t *testing.T) {
	base := map[string]interface{}{"host": "localhost", "port": 80}
	p := NewPatch([]PatchOp{
		{Op: "set", Path: "port", Value: 443},
		{Op: "delete", Path: "host"},
	})
	_, err := p.Merge(base)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if base["port"] != 80 {
		t.Error("base config was mutated: port changed")
	}
	if _, ok := base["host"]; !ok {
		t.Error("base config was mutated: host deleted")
	}
}
