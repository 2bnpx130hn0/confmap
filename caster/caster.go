// Package caster provides type casting utilities for config map values.
// It allows safely converting raw interface{} values from config maps
// into strongly typed Go primitives.
package caster

import (
	"fmt"
	"strconv"
)

// Caster wraps a config map and provides typed accessor methods.
type Caster struct {
	data map[string]interface{}
}

// New creates a new Caster from the given config map.
func New(data map[string]interface{}) *Caster {
	return &Caster{data: data}
}

// String returns the value at key as a string.
// Returns an error if the key is missing or the value cannot be cast.
func (c *Caster) String(key string) (string, error) {
	v, ok := c.data[key]
	if !ok {
		return "", fmt.Errorf("caster: key %q not found", key)
	}
	switch val := v.(type) {
	case string:
		return val, nil
	case int:
		return strconv.Itoa(val), nil
	case float64:
		return strconv.FormatFloat(val, 'f', -1, 64), nil
	case bool:
		return strconv.FormatBool(val), nil
	default:
		return "", fmt.Errorf("caster: cannot cast %T to string for key %q", v, key)
	}
}

// Int returns the value at key as an int.
// Returns an error if the key is missing or the value cannot be cast.
func (c *Caster) Int(key string) (int, error) {
	v, ok := c.data[key]
	if !ok {
		return 0, fmt.Errorf("caster: key %q not found", key)
	}
	switch val := v.(type) {
	case int:
		return val, nil
	case float64:
		return int(val), nil
	case string:
		n, err := strconv.Atoi(val)
		if err != nil {
			return 0, fmt.Errorf("caster: cannot parse %q as int for key %q", val, key)
		}
		return n, nil
	default:
		return 0, fmt.Errorf("caster: cannot cast %T to int for key %q", v, key)
	}
}

// Bool returns the value at key as a bool.
// Returns an error if the key is missing or the value cannot be cast.
func (c *Caster) Bool(key string) (bool, error) {
	v, ok := c.data[key]
	if !ok {
		return false, fmt.Errorf("caster: key %q not found", key)
	}
	switch val := v.(type) {
	case bool:
		return val, nil
	case string:
		b, err := strconv.ParseBool(val)
		if err != nil {
			return false, fmt.Errorf("caster: cannot parse %q as bool for key %q", val, key)
		}
		return b, nil
	default:
		return false, fmt.Errorf("caster: cannot cast %T to bool for key %q", v, key)
	}
}
