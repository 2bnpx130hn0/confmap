package resolver_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/yourorg/confmap/loader"
	"github.com/yourorg/confmap/resolver"
)

func writeTempYAML(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "confmap-reload-*.yaml")
	if err != nil {
		t.Fatalf("create temp: %v", err)
	}
	_, _ = f.WriteString(content)
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestWatchedResolver_ReloadsOnChange(t *testing.T) {
	path := writeTempYAML(t, "host: localhost\nport: 8080\n")

	fl := loader.NewFileLoader(path)
	r := resolver.New([]loader.Loader{fl}, nil)

	wr, err := resolver.NewWatched(r, []string{path}, 20*time.Millisecond)
	if err != nil {
		t.Fatalf("NewWatched: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	wr.Start(ctx)

	cfg, err := wr.Resolve()
	if err != nil {
		t.Fatalf("initial Resolve: %v", err)
	}
	if cfg["host"] != "localhost" {
		t.Errorf("expected host=localhost, got %v", cfg["host"])
	}

	// Update the file
	if err := os.WriteFile(path, []byte("host: remotehost\nport: 9090\n"), 0644); err != nil {
		t.Fatalf("write: %v", err)
	}
	time.Sleep(80 * time.Millisecond)

	cfg2, err := wr.Resolve()
	if err != nil {
		t.Fatalf("post-change Resolve: %v", err)
	}
	if cfg2["host"] != "remotehost" {
		t.Errorf("expected host=remotehost after reload, got %v", cfg2["host"])
	}
}

func TestWatchedResolver_MissingFileError(t *testing.T) {
	fl := loader.NewFileLoader("/no/such/file.yaml")
	r := resolver.New([]loader.Loader{fl}, nil)
	_, err := resolver.NewWatched(r, []string{"/no/such/file.yaml"}, 50*time.Millisecond)
	if err == nil {
		t.Error("expected error for missing watched file")
	}
}
