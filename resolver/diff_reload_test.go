package resolver_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/your-org/confmap/differ"
	"github.com/your-org/confmap/loader"
	"github.com/your-org/confmap/merger"
	"github.com/your-org/confmap/resolver"
	"github.com/your-org/confmap/watcher"
)

func writeDiffYAML(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestWatchedWithDiff_ReportsDeltas(t *testing.T) {
	dir := t.TempDir()
	path := writeDiffYAML(t, dir, "cfg.yaml", "host: localhost\nport: 3000\n")

	fl := loader.NewFileLoader(path)
	m := merger.New()
	w := watcher.New([]string{path})
	wr, err := resolver.NewWatched([]loader.Loader{fl}, m, nil, w)
	if err != nil {
		t.Fatalf("NewWatched error: %v", err)
	}

	var received []differ.Delta
	doneCh := make(chan struct{})

	wd := resolver.NewWatchedWithDiff(wr, func(deltas []differ.Delta) {
		received = deltas
		close(doneCh)
	})

	if err := wd.StartWithDiff(); err != nil {
		t.Fatalf("StartWithDiff error: %v", err)
	}
	defer wd.Stop()

	time.Sleep(50 * time.Millisecond)
	if err := os.WriteFile(path, []byte("host: prod.example.com\nport: 3000\n"), 0644); err != nil {
		t.Fatal(err)
	}

	select {
	case <-doneCh:
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for diff callback")
	}

	if len(received) != 1 {
		t.Fatalf("expected 1 delta, got %d", len(received))
	}
	if received[0].Key != "host" || received[0].Type != differ.Changed {
		t.Errorf("unexpected delta: %+v", received[0])
	}
}
