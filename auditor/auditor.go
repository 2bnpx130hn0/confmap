// Package auditor provides change-tracking for config maps,
// recording when keys were set, updated, or deleted along with
// a timestamp and optional source label.
package auditor

import (
	"fmt"
	"sync"
	"time"
)

// EventKind describes the type of audit event.
type EventKind string

const (
	EventSet    EventKind = "set"
	EventUpdate EventKind = "update"
	EventDelete EventKind = "delete"
)

// Event records a single change to a config key.
type Event struct {
	Key       string
	Kind      EventKind
	OldValue  interface{}
	NewValue  interface{}
	Source    string
	Timestamp time.Time
}

// Auditor tracks config mutations and exposes an immutable log.
type Auditor struct {
	mu     sync.RWMutex
	events []Event
}

// New returns a ready-to-use Auditor.
func New() *Auditor {
	return &Auditor{}
}

// Record compares oldCfg and newCfg (flat or nested map[string]interface{})
// and appends an Event for every key that was added, changed, or removed.
// source is a human-readable label (e.g. "file:app.yaml", "env").
func (a *Auditor) Record(oldCfg, newCfg map[string]interface{}, source string) {
	now := time.Now().UTC()
	a.mu.Lock()
	defer a.mu.Unlock()

	for k, nv := range newCfg {
		ov, exists := oldCfg[k]
		if !exists {
			a.events = append(a.events, Event{Key: k, Kind: EventSet, NewValue: nv, Source: source, Timestamp: now})
		} else if fmt.Sprintf("%v", ov) != fmt.Sprintf("%v", nv) {
			a.events = append(a.events, Event{Key: k, Kind: EventUpdate, OldValue: ov, NewValue: nv, Source: source, Timestamp: now})
		}
	}

	for k, ov := range oldCfg {
		if _, exists := newCfg[k]; !exists {
			a.events = append(a.events, Event{Key: k, Kind: EventDelete, OldValue: ov, Source: source, Timestamp: now})
		}
	}
}

// Events returns a snapshot of all recorded events.
func (a *Auditor) Events() []Event {
	a.mu.RLock()
	defer a.mu.RUnlock()
	out := make([]Event, len(a.events))
	copy(out, a.events)
	return out
}

// Clear resets the event log.
func (a *Auditor) Clear() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.events = nil
}
