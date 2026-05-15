package matcher_test

import (
	"sort"
	"testing"

	"github.com/your-org/confmap/matcher"
)

func baseConfig() map[string]any {
	return map[string]any{
		"db": map[string]any{
			"host": "localhost",
			"port": 5432,
		},
		"app": map[string]any{
			"port": 8080,
			"name": "myapp",
		},
		"debug": true,
	}
}

func TestMatch_ExactKey(t *testing.T) {
	m := matcher.New("debug")
	if !m.Match("debug") {
		t.Fatal("expected debug to match")
	}
	if m.Match("db.host") {
		t.Fatal("expected db.host not to match")
	}
}

func TestMatch_WildcardSuffix(t *testing.T) {
	m := matcher.New("db.*")
	if !m.Match("db.host") {
		t.Fatal("expected db.host to match db.*")
	}
	if !m.Match("db.port") {
		t.Fatal("expected db.port to match db.*")
	}
	if m.Match("app.port") {
		t.Fatal("expected app.port not to match db.*")
	}
}

func TestFilter_ReturnsMatchingKeys(t *testing.T) {
	m := matcher.New("db.*")
	result := m.Filter(baseConfig())

	db, ok := result["db"].(map[string]any)
	if !ok {
		t.Fatal("expected db map in result")
	}
	if _, ok := db["host"]; !ok {
		t.Error("expected db.host in result")
	}
	if _, ok := db["port"]; !ok {
		t.Error("expected db.port in result")
	}
	if _, ok := result["app"]; ok {
		t.Error("expected app to be excluded")
	}
	if _, ok := result["debug"]; ok {
		t.Error("expected debug to be excluded")
	}
}

func TestFilter_MultiplePatterns(t *testing.T) {
	m := matcher.New("db.*", "app.port")
	result := m.Filter(baseConfig())

	app, ok := result["app"].(map[string]any)
	if !ok {
		t.Fatal("expected app map in result")
	}
	if _, ok := app["port"]; !ok {
		t.Error("expected app.port in result")
	}
	if _, ok := app["name"]; ok {
		t.Error("expected app.name to be excluded")
	}
}

func TestKeys_ReturnsSortedMatchingKeys(t *testing.T) {
	m := matcher.New("db.*")
	keys := m.Keys(baseConfig())
	sort.Strings(keys)

	if len(keys) != 2 {
		t.Fatalf("expected 2 keys, got %d: %v", len(keys), keys)
	}
	if keys[0] != "db.host" || keys[1] != "db.port" {
		t.Errorf("unexpected keys: %v", keys)
	}
}

func TestFilter_NoPatterns_ReturnsEmpty(t *testing.T) {
	m := matcher.New()
	result := m.Filter(baseConfig())
	if len(result) != 0 {
		t.Errorf("expected empty result, got %v", result)
	}
}
