// Package exporter provides functionality to serialize a merged config map
// into various output formats such as YAML, TOML, and JSON.
package exporter

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v3"
)

// Format represents the output serialization format.
type Format string

const (
	FormatYAML Format = "yaml"
	FormatTOML Format = "toml"
	FormatJSON Format = "json"
)

// Export serializes the given config map into the specified format.
// It returns the serialized bytes or an error if serialization fails.
func Export(cfg map[string]any, format Format) ([]byte, error) {
	switch format {
	case FormatYAML:
		return exportYAML(cfg)
	case FormatTOML:
		return exportTOML(cfg)
	case FormatJSON:
		return exportJSON(cfg)
	default:
		return nil, fmt.Errorf("exporter: unsupported format %q", format)
	}
}

func exportYAML(cfg map[string]any) ([]byte, error) {
	out, err := yaml.Marshal(cfg)
	if err != nil {
		return nil, fmt.Errorf("exporter: yaml marshal failed: %w", err)
	}
	return out, nil
}

func exportTOML(cfg map[string]any) ([]byte, error) {
	var sb strings.Builder
	enc := toml.NewEncoder(&sb)
	if err := enc.Encode(cfg); err != nil {
		return nil, fmt.Errorf("exporter: toml marshal failed: %w", err)
	}
	return []byte(sb.String()), nil
}

func exportJSON(cfg map[string]any) ([]byte, error) {
	out, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("exporter: json marshal failed: %w", err)
	}
	return out, nil
}
