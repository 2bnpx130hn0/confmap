package merger

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntersectMerge_CommonKeysRetained(t *testing.T) {
	m := NewIntersect()
	m.AddLayer(map[string]any{"a": 1, "b": 2, "c": 3})
	m.AddLayer(map[string]any{"a": 10, "b": 20, "d": 40})

	result := m.Merge()

	require.Contains(t, result, "a")
	require.Contains(t, result, "b")
	assert.NotContains(t, result, "c", "key only in first layer should be excluded")
	assert.NotContains(t, result, "d", "key only in second layer should be excluded")
}

func TestIntersectMerge_LastLayerValueWins(t *testing.T) {
	m := NewIntersect()
	m.AddLayer(map[string]any{"x": "first"})
	m.AddLayer(map[string]any{"x": "second"})
	m.AddLayer(map[string]any{"x": "third"})

	result := m.Merge()

	assert.Equal(t, "third", result["x"])
}

func TestIntersectMerge_NoCommonKeys(t *testing.T) {
	m := NewIntersect()
	m.AddLayer(map[string]any{"a": 1})
	m.AddLayer(map[string]any{"b": 2})

	result := m.Merge()

	assert.Empty(t, result)
}

func TestIntersectMerge_NilLayerSkipped(t *testing.T) {
	m := NewIntersect()
	m.AddLayer(map[string]any{"a": 1, "b": 2})
	m.AddLayer(nil)
	m.AddLayer(map[string]any{"a": 99, "b": 88})

	result := m.Merge()

	// nil layer is skipped; both non-nil layers share "a" and "b"
	assert.Equal(t, 99, result["a"])
	assert.Equal(t, 88, result["b"])
}

func TestIntersectMerge_EmptyLayers(t *testing.T) {
	m := NewIntersect()

	result := m.Merge()

	assert.NotNil(t, result)
	assert.Empty(t, result)
}

func TestIntersectMerge_SingleLayer(t *testing.T) {
	m := NewIntersect()
	m.AddLayer(map[string]any{"only": true})

	result := m.Merge()

	assert.Equal(t, map[string]any{"only": true}, result)
}

func TestIntersectMerge_ThreeLayersPartialOverlap(t *testing.T) {
	m := NewIntersect()
	m.AddLayer(map[string]any{"a": 1, "b": 2, "c": 3})
	m.AddLayer(map[string]any{"a": 10, "b": 20, "e": 50})
	m.AddLayer(map[string]any{"a": 100, "c": 300, "e": 500})

	result := m.Merge()

	// Only "a" appears in all three layers.
	require.Len(t, result, 1)
	assert.Equal(t, 100, result["a"])
}
