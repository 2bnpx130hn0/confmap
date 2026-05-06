// Package auditor provides change-tracking for confmap configuration maps.
//
// It records every key-level mutation (set, update, delete) that occurs
// between two snapshots of a config map, together with the originating
// source label and a UTC timestamp.
//
// Basic usage:
//
//	aud := auditor.New()
//	old := map[string]interface{}{"port": 8080}
//	new_ := map[string]interface{}{"port": 9090, "host": "localhost"}
//	aud.Record(old, new_, "file:app.yaml")
//
//	for _, e := range aud.Events() {
//		fmt.Printf("%s %s %v -> %v (%s)\n", e.Timestamp, e.Kind, e.OldValue, e.NewValue, e.Source)
//	}
package auditor
