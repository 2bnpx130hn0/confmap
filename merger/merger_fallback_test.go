package merger

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFallback_PrimaryWins(t *testing.T) {
	primary := map[string]any{"host": "prod.example.com", "port": 9090}
	fallback := map[string]any{"host": "localhost", "port": 8080, "debug": false}

	result := NewFallback(primary, fallback).Merge()

	assert.Equal(t, "prod.example.com", result["host"])
	assert.Equal(t, 9090, result["port"])
	assert.Equal(t, false, result["debug"], "fallback key should be present")
}

func TestFallback_MissingPrimaryKeyFilledFromFallback(t *testing.T) {
	primary := map[string]any{"host": "prod.example.com"}
	fallback := map[string]any{"host": "localhost", "timeout": 30}

	result := NewFallback(primary, fallback).Merge()

	assert.Equal(t, "prod.example.com", result["host"])
	assert.Equal(t, 30, result["timeout"])
}

func TestFallback_NilPrimaryValueSkipped(t *testing.T) {
	primary := map[string]any{"host": nil}
	fallback := map[string]any{"host": "localhost"}

	result := NewFallback(primary, fallback).Merge()

	assert.Equal(t, "localhost", result["host"], "nil primary value should fall back")
}

func TestFallback_NestedMapsAreMergedRecursively(t *testing.T) {
	primary := map[string]any{
		"db": map[string]any{"host": "db.prod"},
	}
	fallback := map[string]any{
		"db": map[string]any{"host": "db.local", "port": 5432},
	}

	result := NewFallback(primary, fallback).Merge()

	require.IsType(t, map[string]any{}, result["db"])
	db := result["db"].(map[string]any)
	assert.Equal(t, "db.prod", db["host"])
	assert.Equal(t, 5432, db["port"], "nested fallback key should be preserved")
}

func TestFallback_DoesNotMutatePrimary(t *testing.T) {
	primary := map[string]any{"a": 1}
	fallback := map[string]any{"b": 2}

	NewFallback(primary, fallback).Merge()

	_, exists := primary["b"]
	assert.False(t, exists, "primary should not be mutated")
}

func TestFallback_DoesNotMutateFallback(t *testing.T) {
	primary := map[string]any{"a": 1}
	fallback := map[string]any{"b": 2}

	NewFallback(primary, fallback).Merge()

	_, exists := fallback["a"]
	assert.False(t, exists, "fallback should not be mutated")
}

func TestFallback_EmptyPrimaryReturnsAllFallback(t *testing.T) {
	primary := map[string]any{}
	fallback := map[string]any{"x": 10, "y": 20}

	result := NewFallback(primary, fallback).Merge()

	assert.Equal(t, 10, result["x"])
	assert.Equal(t, 20, result["y"])
}

func TestFallback_EmptyFallbackReturnsPrimary(t *testing.T) {
	primary := map[string]any{"a": "hello"}
	fallback := map[string]any{}

	result := NewFallback(primary, fallback).Merge()

	assert.Equal(t, "hello", result["a"])
	assert.Len(t, result, 1)
}
