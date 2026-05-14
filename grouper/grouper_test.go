package grouper_test

import (
	"testing"

	"github.com/iamNilotpal/confmap/grouper"
)

var base = map[string]any{
	"db_host":    "localhost",
	"db_port":    5432,
	"cache_host": "redis",
	"cache_ttl":  300,
	"debug":      true,
}

func TestByPrefix_SplitsCorrectly(t *testing.T) {
	g := grouper.New(base)
	groups := g.ByPrefix([]string{"db", "cache"}, "_")

	if groups["db"]["host"] != "localhost" {
		t.Errorf("expected db.host=localhost, got %v", groups["db"]["host"])
	}
	if groups["db"]["port"] != 5432 {
		t.Errorf("expected db.port=5432, got %v", groups["db"]["port"])
	}
	if groups["cache"]["host"] != "redis" {
		t.Errorf("expected cache.host=redis, got %v", groups["cache"]["host"])
	}
	if groups["cache"]["ttl"] != 300 {
		t.Errorf("expected cache.ttl=300, got %v", groups["cache"]["ttl"])
	}
}

func TestByPrefix_UnmatchedGoesToEmpty(t *testing.T) {
	g := grouper.New(base)
	groups := g.ByPrefix([]string{"db", "cache"}, "_")

	if _, ok := groups[""]["debug"]; !ok {
		t.Error("expected unmatched key 'debug' in empty group")
	}
}

func TestByPrefix_NoMatchingPrefixes(t *testing.T) {
	g := grouper.New(base)
	groups := g.ByPrefix([]string{"svc"}, "_")

	if len(groups[""]) != len(base) {
		t.Errorf("expected all keys in empty group, got %d", len(groups[""]))
	}
	if len(groups["svc"]) != 0 {
		t.Errorf("expected empty svc group, got %d keys", len(groups["svc"]))
	}
}

func TestByFunc_CustomGrouping(t *testing.T) {
	g := grouper.New(base)
	groups := g.ByFunc(func(key string) string {
		if key == "debug" {
			return "meta"
		}
		return "other"
	})

	if _, ok := groups["meta"]["debug"]; !ok {
		t.Error("expected 'debug' in meta group")
	}
	if len(groups["other"]) != 4 {
		t.Errorf("expected 4 keys in other group, got %d", len(groups["other"]))
	}
}

func TestGroups_ReturnsMatchedPrefixes(t *testing.T) {
	g := grouper.New(base)
	names := g.Groups([]string{"db", "cache", "svc"}, "_")

	found := map[string]bool{}
	for _, n := range names {
		found[n] = true
	}
	if !found["db"] {
		t.Error("expected 'db' in groups")
	}
	if !found["cache"] {
		t.Error("expected 'cache' in groups")
	}
	if found["svc"] {
		t.Error("did not expect 'svc' in groups")
	}
}

func TestByPrefix_EmptyConfig(t *testing.T) {
	g := grouper.New(map[string]any{})
	groups := g.ByPrefix([]string{"db"}, "_")
	if len(groups) != 0 {
		t.Errorf("expected no groups for empty config, got %d", len(groups))
	}
}
