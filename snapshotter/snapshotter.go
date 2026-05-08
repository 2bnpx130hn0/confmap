// Package snapshotter provides point-in-time config snapshot management,
// allowing callers to capture, list, restore, and compare named snapshots
// of a configuration map.
package snapshotter

import (
	"fmt"
	"sync"
	"time"
)

// Snapshot holds a named, timestamped copy of a config map.
type Snapshot struct {
	Name      string
	CapturedAt time.Time
	Data      map[string]any
}

// Snapshotter manages a collection of named config snapshots.
type Snapshotter struct {
	mu        sync.RWMutex
	snapshots map[string]Snapshot
}

// New returns an initialised Snapshotter.
func New() *Snapshotter {
	return &Snapshotter{snapshots: make(map[string]Snapshot)}
}

// Capture deep-copies cfg and stores it under name.
func (s *Snapshotter) Capture(name string, cfg map[string]any) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.snapshots[name] = Snapshot{
		Name:       name,
		CapturedAt: time.Now(),
		Data:       deepCopy(cfg),
	}
}

// Get returns the snapshot stored under name, or an error if absent.
func (s *Snapshotter) Get(name string) (Snapshot, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	snap, ok := s.snapshots[name]
	if !ok {
		return Snapshot{}, fmt.Errorf("snapshotter: snapshot %q not found", name)
	}
	return snap, nil
}

// List returns the names of all stored snapshots.
func (s *Snapshotter) List() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	names := make([]string, 0, len(s.snapshots))
	for k := range s.snapshots {
		names = append(names, k)
	}
	return names
}

// Delete removes the named snapshot. It is a no-op if the name does not exist.
func (s *Snapshotter) Delete(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.snapshots, name)
}

// Restore returns a deep copy of the data stored in the named snapshot so the
// caller can use it as the current configuration.
func (s *Snapshotter) Restore(name string) (map[string]any, error) {
	snap, err := s.Get(name)
	if err != nil {
		return nil, err
	}
	return deepCopy(snap.Data), nil
}

// deepCopy performs a shallow-recursive copy of a map[string]any.
func deepCopy(src map[string]any) map[string]any {
	dst := make(map[string]any, len(src))
	for k, v := range src {
		if nested, ok := v.(map[string]any); ok {
			dst[k] = deepCopy(nested)
		} else {
			dst[k] = v
		}
	}
	return dst
}
