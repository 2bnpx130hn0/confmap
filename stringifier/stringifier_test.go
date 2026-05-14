package stringifier_test

import (
	"strings"
	"testing"

	"github.com/yourusername/confmap/stringifier"
)

func TestRender_FlatMap(t *testing.T) {
	s := stringifier.New(stringifier.Options{SortKeys: true})
	cfg := map[string]any{"host": "localhost", "port": 8080}
	lines := s.Render(cfg)
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
	if lines[0] != "host=localhost" {
		t.Errorf("unexpected line: %s", lines[0])
	}
	if lines[1] != "port=8080" {
		t.Errorf("unexpected line: %s", lines[1])
	}
}

func TestRender_NestedMap(t *testing.T) {
	s := stringifier.New(stringifier.Options{SortKeys: true})
	cfg := map[string]any{
		"db": map[string]any{"host": "pg", "port": 5432},
	}
	lines := s.Render(cfg)
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
	if lines[0] != "db.host=pg" {
		t.Errorf("unexpected: %s", lines[0])
	}
	if lines[1] != "db.port=5432" {
		t.Errorf("unexpected: %s", lines[1])
	}
}

func TestRender_QuoteValues(t *testing.T) {
	s := stringifier.New(stringifier.Options{QuoteValues: true, SortKeys: true})
	cfg := map[string]any{"key": "val"}
	lines := s.Render(cfg)
	if lines[0] != `key="val"` {
		t.Errorf("expected quoted value, got %s", lines[0])
	}
}

func TestRender_CustomDelimiter(t *testing.T) {
	s := stringifier.New(stringifier.Options{Delimiter: ": ", SortKeys: true})
	cfg := map[string]any{"foo": "bar"}
	lines := s.Render(cfg)
	if lines[0] != "foo: bar" {
		t.Errorf("unexpected: %s", lines[0])
	}
}

func TestRender_Prefix(t *testing.T) {
	s := stringifier.New(stringifier.Options{Prefix: "export ", SortKeys: true})
	cfg := map[string]any{"env": "prod"}
	lines := s.Render(cfg)
	if lines[0] != "export env=prod" {
		t.Errorf("unexpected: %s", lines[0])
	}
}

func TestString_JoinsWithNewline(t *testing.T) {
	s := stringifier.New(stringifier.Options{SortKeys: true})
	cfg := map[string]any{"a": 1, "b": 2}
	out := s.String(cfg)
	if !strings.Contains(out, "\n") {
		t.Errorf("expected newline separator, got: %s", out)
	}
}

func TestRender_EmptyConfig(t *testing.T) {
	s := stringifier.New(stringifier.Options{})
	lines := s.Render(map[string]any{})
	if len(lines) != 0 {
		t.Errorf("expected empty slice, got %v", lines)
	}
}

func TestRender_DefaultDelimiter(t *testing.T) {
	s := stringifier.New(stringifier.Options{})
	cfg := map[string]any{"x": "y"}
	lines := s.Render(cfg)
	if !strings.Contains(lines[0], "=") {
		t.Errorf("expected default '=' delimiter, got %s", lines[0])
	}
}
