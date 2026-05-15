package merger

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSchemaMerge_CompatibleVersions(t *testing.T) {
	sm := NewSchema()
	sm.AddLayer("1.0.0", map[string]any{"host": "localhost"})
	sm.AddLayer("1.2.0", map[string]any{"port": 8080})

	out, err := sm.Merge()
	require.NoError(t, err)
	assert.Equal(t, "localhost", out["host"])
	assert.Equal(t, 8080, out["port"])
}

func TestSchemaMerge_IncompatibleVersions(t *testing.T) {
	sm := NewSchema()
	sm.AddLayer("1.0.0", map[string]any{"host": "localhost"})
	sm.AddLayer("2.0.0", map[string]any{"port": 9090})

	_, err := sm.Merge()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "incompatible major versions")
}

func TestSchemaMerge_NilLayerSkipped(t *testing.T) {
	sm := NewSchema()
	sm.AddLayer("3.1.0", map[string]any{"a": 1})
	sm.AddLayer("3.2.0", nil)
	sm.AddLayer("3.3.0", map[string]any{"b": 2})

	out, err := sm.Merge()
	require.NoError(t, err)
	assert.Equal(t, 1, out["a"])
	assert.Equal(t, 2, out["b"])
}

func TestSchemaMerge_EmptyLayers(t *testing.T) {
	sm := NewSchema()
	out, err := sm.Merge()
	require.NoError(t, err)
	assert.Empty(t, out)
}

func TestSchemaMerge_OverrideOrder(t *testing.T) {
	sm := NewSchema()
	sm.AddLayer("2.0", map[string]any{"key": "first"})
	sm.AddLayer("2.1", map[string]any{"key": "second"})

	out, err := sm.Merge()
	require.NoError(t, err)
	assert.Equal(t, "second", out["key"])
}

func TestSchemaMerge_InvalidBaseVersion(t *testing.T) {
	sm := NewSchema()
	sm.AddLayer("", map[string]any{"x": 1})

	_, err := sm.Merge()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid")
}

func TestSchemaMerge_NoDotVersion(t *testing.T) {
	sm := NewSchema()
	sm.AddLayer("v1", map[string]any{"a": true})
	sm.AddLayer("v1", map[string]any{"b": false})

	out, err := sm.Merge()
	require.NoError(t, err)
	assert.Equal(t, true, out["a"])
	assert.Equal(t, false, out["b"])
}
