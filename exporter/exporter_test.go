package exporter_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/yourorg/confmap/exporter"
)

var sampleCfg = map[string]any{
	"host": "localhost",
	"port": 8080,
	"debug": true,
}

func TestExport_YAML(t *testing.T) {
	out, err := exporter.Export(sampleCfg, exporter.FormatYAML)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	s := string(out)
	if !strings.Contains(s, "host") || !strings.Contains(s, "localhost") {
		t.Errorf("expected YAML to contain host key, got:\n%s", s)
	}
	if !strings.Contains(s, "port") {
		t.Errorf("expected YAML to contain port key, got:\n%s", s)
	}
}

func TestExport_TOML(t *testing.T) {
	out, err := exporter.Export(sampleCfg, exporter.FormatTOML)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	s := string(out)
	if !strings.Contains(s, "host") || !strings.Contains(s, "localhost") {
		t.Errorf("expected TOML to contain host key, got:\n%s", s)
	}
}

func TestExport_JSON(t *testing.T) {
	out, err := exporter.Export(sampleCfg, exporter.FormatJSON)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var parsed map[string]any
	if err := json.Unmarshal(out, &parsed); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if parsed["host"] != "localhost" {
		t.Errorf("expected host=localhost, got %v", parsed["host"])
	}
}

func TestExport_UnsupportedFormat(t *testing.T) {
	_, err := exporter.Export(sampleCfg, exporter.Format("xml"))
	if err == nil {
		t.Fatal("expected error for unsupported format, got nil")
	}
	if !strings.Contains(err.Error(), "unsupported format") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestExport_EmptyConfig(t *testing.T) {
	cfg := map[string]any{}
	for _, fmt := range []exporter.Format{exporter.FormatYAML, exporter.FormatTOML, exporter.FormatJSON} {
		out, err := exporter.Export(cfg, fmt)
		if err != nil {
			t.Errorf("format %s: unexpected error: %v", fmt, err)
		}
		if out == nil {
			t.Errorf("format %s: expected non-nil output for empty config", fmt)
		}
	}
}
