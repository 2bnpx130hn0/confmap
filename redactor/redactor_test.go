package redactor_test

import (
	"sort"
	"testing"

	"github.com/user/confmap/redactor"
)

func TestApply_RedactsMatchingKeys(t *testing.T) {
	r := redactor.New("password", "secret")
	cfg := map[string]any{
		"username": "alice",
		"password": "s3cr3t",
		"api_secret": "tok123",
	}
	out := r.Apply(cfg)
	if out["username"] != "alice" {
		t.Errorf("expected username to be unchanged, got %v", out["username"])
	}
	if out["password"] != "[REDACTED]" {
		t.Errorf("expected password to be redacted, got %v", out["password"])
	}
	if out["api_secret"] != "[REDACTED]" {
		t.Errorf("expected api_secret to be redacted, got %v", out["api_secret"])
	}
}

func TestApply_CaseInsensitiveMatch(t *testing.T) {
	r := redactor.New("token")
	cfg := map[string]any{"AUTH_TOKEN": "abc", "host": "localhost"}
	out := r.Apply(cfg)
	if out["AUTH_TOKEN"] != "[REDACTED]" {
		t.Errorf("expected AUTH_TOKEN redacted, got %v", out["AUTH_TOKEN"])
	}
	if out["host"] != "localhost" {
		t.Errorf("expected host unchanged, got %v", out["host"])
	}
}

func TestApply_NestedMap(t *testing.T) {
	r := redactor.New("secret")
	cfg := map[string]any{
		"db": map[string]any{
			"host":   "localhost",
			"secret": "dbpass",
		},
	}
	out := r.Apply(cfg)
	db, ok := out["db"].(map[string]any)
	if !ok {
		t.Fatal("expected nested map under 'db'")
	}
	if db["host"] != "localhost" {
		t.Errorf("expected host unchanged, got %v", db["host"])
	}
	if db["secret"] != "[REDACTED]" {
		t.Errorf("expected secret redacted, got %v", db["secret"])
	}
}

func TestApply_DoesNotMutateOriginal(t *testing.T) {
	r := redactor.New("password")
	cfg := map[string]any{"password": "original"}
	r.Apply(cfg)
	if cfg["password"] != "original" {
		t.Error("Apply must not mutate the original config")
	}
}

func TestKeys_ReturnsSensitivePaths(t *testing.T) {
	r := redactor.New("secret", "token")
	cfg := map[string]any{
		"app_token": "t1",
		"name":      "svc",
		"db": map[string]any{
			"secret": "pass",
			"port":   5432,
		},
	}
	keys := r.Keys(cfg)
	sort.Strings(keys)
	expected := []string{"app_token", "db.secret"}
	if len(keys) != len(expected) {
		t.Fatalf("expected %v, got %v", expected, keys)
	}
	for i, k := range keys {
		if k != expected[i] {
			t.Errorf("expected %s, got %s", expected[i], k)
		}
	}
}

func TestApply_NoPatterns(t *testing.T) {
	r := redactor.New()
	cfg := map[string]any{"password": "s3cr3t", "host": "localhost"}
	out := r.Apply(cfg)
	if out["password"] != "s3cr3t" {
		t.Error("expected no redaction with empty pattern list")
	}
}
