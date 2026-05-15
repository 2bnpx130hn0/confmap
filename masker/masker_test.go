package masker_test

import (
	"testing"

	"github.com/your-org/confmap/masker"
)

func TestApply_MasksSensitiveKey(t *testing.T) {
	m := masker.New([]string{"password"}, 2, "*")
	cfg := map[string]any{"password": "secret99", "user": "alice"}
	out := m.Apply(cfg)

	if out["password"] != "******99" {
		t.Errorf("expected ******99, got %v", out["password"])
	}
	if out["user"] != "alice" {
		t.Errorf("expected alice, got %v", out["user"])
	}
}

func TestApply_CaseInsensitiveKeyMatch(t *testing.T) {
	m := masker.New([]string{"apikey"}, 3, "#")
	cfg := map[string]any{"ApiKey": "abcdef"}
	out := m.Apply(cfg)
	if out["ApiKey"] != "###def" {
		t.Errorf("expected ###def, got %v", out["ApiKey"])
	}
}

func TestApply_NestedMap(t *testing.T) {
	m := masker.New([]string{"token"}, 0, "*")
	cfg := map[string]any{
		"db": map[string]any{
			"token": "mysecret",
			"host":  "localhost",
		},
	}
	out := m.Apply(cfg)
	db, _ := out["db"].(map[string]any)
	if db["token"] != "********" {
		t.Errorf("expected ********, got %v", db["token"])
	}
	if db["host"] != "localhost" {
		t.Errorf("expected localhost, got %v", db["host"])
	}
}

func TestApply_DoesNotMutateOriginal(t *testing.T) {
	m := masker.New([]string{"secret"}, 1, "*")
	cfg := map[string]any{"secret": "abc"}
	m.Apply(cfg)
	if cfg["secret"] != "abc" {
		t.Error("original config was mutated")
	}
}

func TestApply_NilConfig(t *testing.T) {
	m := masker.New([]string{"x"}, 2, "*")
	if got := m.Apply(nil); got != nil {
		t.Errorf("expected nil, got %v", got)
	}
}

func TestApply_ShortValueFullyMasked(t *testing.T) {
	m := masker.New([]string{"pin"}, 4, "*")
	cfg := map[string]any{"pin": "12"}
	out := m.Apply(cfg)
	// visible > len, so all chars masked
	if out["pin"] != "**" {
		t.Errorf("expected **, got %v", out["pin"])
	}
}

func TestIsSensitive(t *testing.T) {
	m := masker.New([]string{"Password", "TOKEN"}, 2, "*")
	if !m.IsSensitive("password") {
		t.Error("expected password to be sensitive")
	}
	if !m.IsSensitive("TOKEN") {
		t.Error("expected TOKEN to be sensitive")
	}
	if m.IsSensitive("username") {
		t.Error("expected username not to be sensitive")
	}
}

func TestApply_NonStringValuesUnchanged(t *testing.T) {
	m := masker.New([]string{"count"}, 2, "*")
	cfg := map[string]any{"count": 42}
	out := m.Apply(cfg)
	if out["count"] != 42 {
		t.Errorf("expected 42, got %v", out["count"])
	}
}
