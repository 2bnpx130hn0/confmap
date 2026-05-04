package loader_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/example/confmap/loader"
)

func writeTempFile(t *testing.T, name, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("writeTempFile: %v", err)
	}
	return path
}

func TestFileLoader_YAML(t *testing.T) {
	path := writeTempFile(t, "config.yaml", "host: localhost\nport: 8080\n")
	fl := &loader.FileLoader{Path: path}
	cfg, err := fl.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg["host"] != "localhost" {
		t.Errorf("expected host=localhost, got %v", cfg["host"])
	}
	if cfg["port"] != 8080 {
		t.Errorf("expected port=8080, got %v", cfg["port"])
	}
}

func TestFileLoader_TOML(t *testing.T) {
	path := writeTempFile(t, "config.toml", "host = \"127.0.0.1\"\nport = 9090\n")
	fl := &loader.FileLoader{Path: path}
	cfg, err := fl.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg["host"] != "127.0.0.1" {
		t.Errorf("expected host=127.0.0.1, got %v", cfg["host"])
	}
}

func TestFileLoader_UnsupportedExtension(t *testing.T) {
	path := writeTempFile(t, "config.json", `{}`)
	fl := &loader.FileLoader{Path: path}
	_, err := fl.Load()
	if err == nil {
		t.Fatal("expected error for unsupported extension")
	}
}

func TestEnvLoader(t *testing.T) {
	t.Setenv("APP_HOST", "envhost")
	t.Setenv("APP_PORT", "3000")
	t.Setenv("OTHER_KEY", "ignored")

	el := &loader.EnvLoader{Prefix: "APP"}
	cfg, err := el.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg["host"] != "envhost" {
		t.Errorf("expected host=envhost, got %v", cfg["host"])
	}
	if cfg["port"] != "3000" {
		t.Errorf("expected port=3000, got %v", cfg["port"])
	}
	if _, ok := cfg["other_key"]; ok {
		t.Error("OTHER_KEY should not be present")
	}
}

func TestFileLoader_MissingFile(t *testing.T) {
	fl := &loader.FileLoader{Path: "/nonexistent/path/config.yaml"}
	_, err := fl.Load()
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
