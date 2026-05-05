// Package watcher implements file-system polling for config hot-reload.
//
// Usage:
//
//	w := watcher.New(500*time.Millisecond, func(path string) error {
//	    // re-load and apply the changed config file
//	    return resolver.Reload(path)
//	})
//	w.Add("/etc/myapp/config.yaml")
//	w.Start(ctx)
//
// The watcher polls registered files at the configured interval. When a
// file's modification time advances, the ReloadFunc is invoked with the
// changed path. Errors from ReloadFunc are logged but do not stop the
// watcher.
//
// The watcher is safe for concurrent use.
package watcher
