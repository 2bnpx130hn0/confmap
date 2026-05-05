package resolver

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/yourorg/confmap/watcher"
)

// WatchedResolver wraps a Resolver with hot-reload support.
type WatchedResolver struct {
	mu       sync.RWMutex
	resolver *Resolver
	paths    []string
	watcher  *watcher.Watcher
}

// NewWatched creates a WatchedResolver that reloads when any source file changes.
func NewWatched(r *Resolver, paths []string, interval time.Duration) (*WatchedResolver, error) {
	wr := &WatchedResolver{
		resolver: r,
		paths:    paths,
	}
	w := watcher.New(interval, func(path string) error {
		return wr.reload()
	})
	for _, p := range paths {
		if err := w.Add(p); err != nil {
			return nil, fmt.Errorf("watched resolver: add %s: %w", p, err)
		}
	}
	wr.watcher = w
	return wr, nil
}

// Start begins watching in the background until ctx is cancelled.
func (wr *WatchedResolver) Start(ctx context.Context) {
	wr.watcher.Start(ctx)
}

// Resolve returns the current merged and validated config map.
func (wr *WatchedResolver) Resolve() (map[string]any, error) {
	wr.mu.RLock()
	defer wr.mu.RUnlock()
	return wr.resolver.Resolve()
}

func (wr *WatchedResolver) reload() error {
	_, err := wr.resolver.Resolve()
	return err
}
