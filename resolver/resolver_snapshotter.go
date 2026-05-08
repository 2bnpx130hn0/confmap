package resolver

import (
	"fmt"

	"github.com/your-org/confmap/snapshotter"
)

// SnapshotResolver wraps a Resolver and automatically captures a named
// snapshot every time the configuration is (re)loaded.
type SnapshotResolver struct {
	*Resolver
	snap    *snapshotter.Snapshotter
	label   string
	counter int
}

// NewSnapshotting returns a SnapshotResolver that delegates to r and stores
// each resolved config in snap under sequentially numbered names prefixed by
// label (e.g. "deploy-0", "deploy-1", …).
func NewSnapshotting(r *Resolver, snap *snapshotter.Snapshotter, label string) *SnapshotResolver {
	return &SnapshotResolver{Resolver: r, snap: snap, label: label}
}

// Resolve loads and validates the configuration, then captures a snapshot.
func (sr *SnapshotResolver) Resolve() (map[string]any, error) {
	cfg, err := sr.Resolver.Resolve()
	if err != nil {
		return nil, err
	}
	name := fmt.Sprintf("%s-%d", sr.label, sr.counter)
	sr.snap.Capture(name, cfg)
	sr.counter++
	return cfg, nil
}

// Snapshotter returns the underlying snapshotter so callers can inspect,
// restore, or delete individual snapshots.
func (sr *SnapshotResolver) Snapshotter() *snapshotter.Snapshotter {
	return sr.snap
}
