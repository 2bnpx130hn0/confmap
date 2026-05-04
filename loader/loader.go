package loader

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v3"
)

// RawConfig holds the raw key-value pairs loaded from a source.
type RawConfig map[string]interface{}

// Loader defines the interface for loading configuration from a source.
type Loader interface {
	Load() (RawConfig, error)
}

// FileLoader loads configuration from a YAML or TOML file.
type FileLoader struct {
	Path string
}

// Load reads and parses the file based on its extension.
func (f *FileLoader) Load() (RawConfig, error) {
	data, err := os.ReadFile(f.Path)
	if err != nil {
		return nil, fmt.Errorf("loader: reading file %q: %w", f.Path, err)
	}

	ext := strings.ToLower(filepath.Ext(f.Path))
	var cfg RawConfig

	switch ext {
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(data, &cfg); err != nil {
			return nil, fmt.Errorf("loader: parsing YAML %q: %w", f.Path, err)
		}
	case ".toml":
		if err := toml.Unmarshal(data, &cfg); err != nil {
			return nil, fmt.Errorf("loader: parsing TOML %q: %w", f.Path, err)
		}
	default:
		return nil, fmt.Errorf("loader: unsupported file extension %q", ext)
	}

	return cfg, nil
}

// EnvLoader loads configuration from environment variables with a given prefix.
type EnvLoader struct {
	Prefix string
}

// Load reads environment variables that match the prefix, strips the prefix,
// lowercases the key, and returns them as a flat RawConfig.
func (e *EnvLoader) Load() (RawConfig, error) {
	cfg := make(RawConfig)
	prefix := strings.ToUpper(e.Prefix)
	if prefix != "" && !strings.HasSuffix(prefix, "_") {
		prefix += "_"
	}

	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key, val := parts[0], parts[1]
		if prefix == "" || strings.HasPrefix(key, prefix) {
			normKey := strings.ToLower(strings.TrimPrefix(key, prefix))
			cfg[normKey] = val
		}
	}
	return cfg, nil
}
