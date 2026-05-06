package tagging_test

import (
	"testing"

	"github.com/example/confmap/tagging"
)

func TestTag_AndHasTag(t *testing.T) {
	tr := tagging.New()
	tr.Tag("database.host", "sensitive", "required")

	if !tr.HasTag("database.host", "sensitive") {
		t.Error("expected sensitive tag")
	}
	if !tr.HasTag("database.host", "required") {
		t.Error("expected required tag")
	}
	if tr.HasTag("database.host", "optional") {
		t.Error("unexpected optional tag")
	}
}

func TestTags_ReturnsAll(t *testing.T) {
	tr := tagging.New()
	tr.Tag("port", "numeric", "optional")
	tags := tr.Tags("port")
	if len(tags) != 2 {
		t.Fatalf("expected 2 tags, got %d", len(tags))
	}
}

func TestTags_MissingKey(t *testing.T) {
	tr := tagging.New()
	if tags := tr.Tags("nonexistent"); len(tags) != 0 {
		t.Errorf("expected empty slice, got %v", tags)
	}
}

func TestFilterByTag(t *testing.T) {
	tr := tagging.New()
	tr.Tag("password", "secret")
	tr.Tag("host", "required")
	tr.Tag("debug", "optional")

	cfg := map[string]any{
		"password": "hunter2",
		"host":     "localhost",
		"debug":    true,
	}

	secrets := tr.FilterByTag(cfg, "secret")
	if len(secrets) != 1 {
		t.Fatalf("expected 1 secret key, got %d", len(secrets))
	}
	if _, ok := secrets["password"]; !ok {
		t.Error("expected password in secrets")
	}
}

func TestAnnotate_AddsTagsMeta(t *testing.T) {
	tr := tagging.New()
	tr.Tag("api_key", "secret")

	cfg := map[string]any{"api_key": "abc123", "timeout": 30}
	annotated := tr.Annotate(cfg)

	meta, ok := annotated["__tags__"]
	if !ok {
		t.Fatal("expected __tags__ annotation")
	}
	tags := meta.(map[string][]string)
	if len(tags["api_key"]) == 0 {
		t.Error("expected api_key to have tags in annotation")
	}
}

func TestAnnotate_NoTaggedKeys(t *testing.T) {
	tr := tagging.New()
	cfg := map[string]any{"foo": "bar"}
	annotated := tr.Annotate(cfg)
	if _, ok := annotated["__tags__"]; ok {
		t.Error("expected no __tags__ key when nothing is tagged")
	}
}

func TestRemove_ExistingKey(t *testing.T) {
	tr := tagging.New()
	tr.Tag("x", "t1")
	if err := tr.Remove("x"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tr.HasTag("x", "t1") {
		t.Error("tag should have been removed")
	}
}

func TestRemove_MissingKey(t *testing.T) {
	tr := tagging.New()
	if err := tr.Remove("ghost"); err == nil {
		t.Error("expected error for missing key")
	}
}
