// Package snapshotter provides point-in-time snapshot management for
// configuration maps.
//
// Usage:
//
//	s := snapshotter.New()
//	s.Capture("before-deploy", currentConfig)
//
//	// … apply changes …
//
//	// Roll back if needed:
//	prev, err := s.Restore("before-deploy")
//	if err != nil {
//		log.Fatal(err)
//	}
//	useConfig(prev)
//
// All operations are safe for concurrent use.
package snapshotter
