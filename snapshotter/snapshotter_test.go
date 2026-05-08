package snapshotter_test

import (
	"testing"

	"github.com/your-org/confmap/snapshotter"
)

func baseConfig() map[string]any {
	return map[string]any{
		"host": "localhost",
		"port": 8080,
		"db": map[string]any{
			"name": "mydb",
		},
	}
}

func TestCapture_And_Get(t *testing.T) {
	s := snapshotter.New()
	s.Capture("v1", baseConfig())

	snap, err := s.Get("v1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if snap.Name != "v1" {
		t.Errorf("expected name v1, got %s", snap.Name)
	}
	if snap.Data["host"] != "localhost" {
		t.Errorf("expected host localhost, got %v", snap.Data["host"])
	}
}

func TestGet_MissingSnapshot(t *testing.T) {
	s := snapshotter.New()
	_, err := s.Get("nonexistent")
	if err == nil {
		t.Fatal("expected error for missing snapshot")
	}
}

func TestCapture_IsDeepCopy(t *testing.T) {
	s := snapshotter.New()
	cfg := baseConfig()
	s.Capture("v1", cfg)

	// Mutate original after capture.
	cfg["host"] = "changed"

	snap, _ := s.Get("v1")
	if snap.Data["host"] != "localhost" {
		t.Errorf("snapshot should be isolated from original mutation, got %v", snap.Data["host"])
	}
}

func TestList_ReturnsAllNames(t *testing.T) {
	s := snapshotter.New()
	s.Capture("v1", baseConfig())
	s.Capture("v2", baseConfig())

	names := s.List()
	if len(names) != 2 {
		t.Errorf("expected 2 snapshots, got %d", len(names))
	}
}

func TestDelete_RemovesSnapshot(t *testing.T) {
	s := snapshotter.New()
	s.Capture("v1", baseConfig())
	s.Delete("v1")

	if len(s.List()) != 0 {
		t.Error("expected empty list after delete")
	}
}

func TestRestore_ReturnsMutableCopy(t *testing.T) {
	s := snapshotter.New()
	s.Capture("v1", baseConfig())

	restored, err := s.Restore("v1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Mutate restored copy and verify snapshot is unchanged.
	restored["host"] = "mutated"
	snap, _ := s.Get("v1")
	if snap.Data["host"] != "localhost" {
		t.Errorf("original snapshot mutated unexpectedly")
	}
}

func TestRestore_MissingSnapshot(t *testing.T) {
	s := snapshotter.New()
	_, err := s.Restore("ghost")
	if err == nil {
		t.Fatal("expected error restoring missing snapshot")
	}
}
