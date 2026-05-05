package watcher_test

import (
	"context"
	"os"
	"sync/atomic"
	"testing"
	"time"

	"github.com/yourorg/confmap/watcher"
)

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "confmap-watch-*.yaml")
	if err != nil {
		t.Fatalf("create temp: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp: %v", err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestWatcher_DetectsChange(t *testing.T) {
	path := writeTempConfig(t, "key: value")

	var callCount int32
	w := watcher.New(20*time.Millisecond, func(p string) error {
		if p == path {
			atomic.AddInt32(&callCount, 1)
		}
		return nil
	})

	if err := w.Add(path); err != nil {
		t.Fatalf("Add: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	w.Start(ctx)

	time.Sleep(30 * time.Millisecond)
	// Modify the file
	if err := os.WriteFile(path, []byte("key: changed"), 0644); err != nil {
		t.Fatalf("write file: %v", err)
	}
	time.Sleep(60 * time.Millisecond)

	if atomic.LoadInt32(&callCount) == 0 {
		t.Error("expected reload callback to be called after file change")
	}
}

func TestWatcher_AddMissingFile(t *testing.T) {
	w := watcher.New(50*time.Millisecond, func(string) error { return nil })
	err := w.Add("/nonexistent/path/config.yaml")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestWatcher_NoCallbackWithoutChange(t *testing.T) {
	path := writeTempConfig(t, "stable: true")

	var callCount int32
	w := watcher.New(20*time.Millisecond, func(string) error {
		atomic.AddInt32(&callCount, 1)
		return nil
	})
	if err := w.Add(path); err != nil {
		t.Fatalf("Add: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	w.Start(ctx)
	time.Sleep(80 * time.Millisecond)

	if atomic.LoadInt32(&callCount) != 0 {
		t.Errorf("expected 0 callbacks, got %d", callCount)
	}
}
