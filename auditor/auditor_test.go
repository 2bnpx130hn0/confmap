package auditor_test

import (
	"testing"

	"github.com/your-org/confmap/auditor"
)

func TestRecord_SetNewKey(t *testing.T) {
	a := auditor.New()
	old := map[string]interface{}{}
	new_ := map[string]interface{}{"host": "localhost"}
	a.Record(old, new_, "test")
	evts := a.Events()
	if len(evts) != 1 {
		t.Fatalf("expected 1 event, got %d", len(evts))
	}
	if evts[0].Kind != auditor.EventSet {
		t.Errorf("expected EventSet, got %s", evts[0].Kind)
	}
	if evts[0].Key != "host" {
		t.Errorf("unexpected key %s", evts[0].Key)
	}
	if evts[0].NewValue != "localhost" {
		t.Errorf("unexpected value %v", evts[0].NewValue)
	}
}

func TestRecord_UpdateExistingKey(t *testing.T) {
	a := auditor.New()
	old := map[string]interface{}{"port": 8080}
	new_ := map[string]interface{}{"port": 9090}
	a.Record(old, new_, "env")
	evts := a.Events()
	if len(evts) != 1 {
		t.Fatalf("expected 1 event, got %d", len(evts))
	}
	if evts[0].Kind != auditor.EventUpdate {
		t.Errorf("expected EventUpdate, got %s", evts[0].Kind)
	}
	if evts[0].OldValue != 8080 || evts[0].NewValue != 9090 {
		t.Errorf("unexpected values old=%v new=%v", evts[0].OldValue, evts[0].NewValue)
	}
}

func TestRecord_DeletedKey(t *testing.T) {
	a := auditor.New()
	old := map[string]interface{}{"debug": true, "host": "localhost"}
	new_ := map[string]interface{}{"host": "localhost"}
	a.Record(old, new_, "file")
	evts := a.Events()
	if len(evts) != 1 {
		t.Fatalf("expected 1 event, got %d", len(evts))
	}
	if evts[0].Kind != auditor.EventDelete {
		t.Errorf("expected EventDelete, got %s", evts[0].Kind)
	}
	if evts[0].Key != "debug" {
		t.Errorf("unexpected key %s", evts[0].Key)
	}
}

func TestRecord_NoChange(t *testing.T) {
	a := auditor.New()
	cfg := map[string]interface{}{"x": 1}
	a.Record(cfg, cfg, "noop")
	if len(a.Events()) != 0 {
		t.Error("expected no events for identical configs")
	}
}

func TestRecord_SourceLabel(t *testing.T) {
	a := auditor.New()
	a.Record(map[string]interface{}{}, map[string]interface{}{"k": "v"}, "file:app.yaml")
	if a.Events()[0].Source != "file:app.yaml" {
		t.Errorf("unexpected source %s", a.Events()[0].Source)
	}
}

func TestClear(t *testing.T) {
	a := auditor.New()
	a.Record(map[string]interface{}{}, map[string]interface{}{"k": "v"}, "s")
	a.Clear()
	if len(a.Events()) != 0 {
		t.Error("expected empty log after Clear")
	}
}

func TestRecord_MultipleChanges(t *testing.T) {
	a := auditor.New()
	old := map[string]interface{}{"host": "localhost", "port": 8080}
	new_ := map[string]interface{}{"host": "remotehost", "timeout": 30}
	a.Record(old, new_, "merge")
	evts := a.Events()
	// Expect: host updated, port deleted, timeout set
	if len(evts) != 3 {
		t.Fatalf("expected 3 events, got %d", len(evts))
	}
	kinds := map[auditor.EventKind]int{}
	for _, e := range evts {
		kinds[e.Kind]++
	}
	if kinds[auditor.EventUpdate] != 1 || kinds[auditor.EventDelete] != 1 || kinds[auditor.EventSet] != 1 {
		t.Errorf("unexpected event kind counts: %v", kinds)
	}
}
