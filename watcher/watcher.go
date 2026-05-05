// Package watcher provides hot-reload functionality for config files.
// It watches registered files for changes and triggers a reload callback.
package watcher

import (
	"context"
	"log"
	"os"
	"sync"
	"time"
)

// ReloadFunc is called when a watched file changes.
type ReloadFunc func(path string) error

// Watcher monitors files for changes and invokes a callback on modification.
type Watcher struct {
	mu       sync.Mutex
	files    map[string]time.Time
	interval time.Duration
	onReload ReloadFunc
}

// New creates a Watcher that polls at the given interval.
func New(interval time.Duration, fn ReloadFunc) *Watcher {
	return &Watcher{
		files:    make(map[string]time.Time),
		interval: interval,
		onReload: fn,
	}
}

// Add registers a file path to be watched.
func (w *Watcher) Add(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	w.mu.Lock()
	w.files[path] = info.ModTime()
	w.mu.Unlock()
	return nil
}

// Start begins polling in a background goroutine until ctx is cancelled.
func (w *Watcher) Start(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(w.interval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				w.poll()
			}
		}
	}()
}

func (w *Watcher) poll() {
	w.mu.Lock()
	defer w.mu.Unlock()
	for path, lastMod := range w.files {
		info, err := os.Stat(path)
		if err != nil {
			log.Printf("watcher: stat error for %s: %v", path, err)
			continue
		}
		if info.ModTime().After(lastMod) {
			w.files[path] = info.ModTime()
			if err := w.onReload(path); err != nil {
				log.Printf("watcher: reload error for %s: %v", path, err)
			}
		}
	}
}
