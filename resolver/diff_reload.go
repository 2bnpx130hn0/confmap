package resolver

import (
	"log"

	"github.com/your-org/confmap/differ"
)

// DiffCallback is called when a config reload produces changes.
// It receives the list of deltas between the old and new config.
type DiffCallback func(deltas []differ.Delta)

// WatchedWithDiff extends WatchedResolver with diff-aware reload callbacks.
type WatchedWithDiff struct {
	*WatchedResolver
	differ   *differ.Differ
	onChange DiffCallback
}

// NewWatchedWithDiff creates a WatchedResolver that also computes and reports
// config diffs on every successful reload.
func NewWatchedWithDiff(w *WatchedResolver, cb DiffCallback) *WatchedWithDiff {
	return &WatchedWithDiff{
		WatchedResolver: w,
		differ:          differ.New(),
		onChange:        cb,
	}
}

// StartWithDiff begins watching for file changes and invokes the DiffCallback
// with computed deltas whenever the config reloads successfully.
func (wd *WatchedWithDiff) StartWithDiff() error {
	return wd.WatchedResolver.Start(func(newCfg map[string]interface{}) {
		oldCfg := wd.WatchedResolver.Config()
		deltas := wd.differ.Diff(oldCfg, newCfg)
		if len(deltas) == 0 {
			log.Println("[confmap] reload: no config changes detected")
			return
		}
		log.Printf("[confmap] reload: %d change(s) detected", len(deltas))
		if wd.onChange != nil {
			wd.onChange(deltas)
		}
	})
}
