// Package coercer provides type coercion for config map values,
// converting between compatible types (e.g. string→int, float64→int).
package coercer

import (
	"fmt"
	"strconv"
)

// Coercer applies type coercion rules to a config map.
type Coercer struct {
	rules []rule
}

type rule struct {
	key      string
	targetFn func(v interface{}) (interface{}, error)
}

// New returns a new Coercer.
func New() *Coercer {
	return &Coercer{}
}

// ToString registers a rule to coerce the value at key to a string.
func (c *Coercer) ToString(key string) *Coercer {
	c.rules = append(c.rules, rule{key: key, targetFn: coerceToString})
	return c
}

// ToInt registers a rule to coerce the value at key to an int.
func (c *Coercer) ToInt(key string) *Coercer {
	c.rules = append(c.rules, rule{key: key, targetFn: coerceToInt})
	return c
}

// ToBool registers a rule to coerce the value at key to a bool.
func (c *Coercer) ToBool(key string) *Coercer {
	c.rules = append(c.rules, rule{key: key, targetFn: coerceToBool})
	return c
}

// Apply runs all registered coercion rules against cfg, modifying it in place.
// Returns an error if any coercion fails.
func (c *Coercer) Apply(cfg map[string]interface{}) error {
	for _, r := range c.rules {
		v, ok := cfg[r.key]
		if !ok {
			continue
		}
		coerced, err := r.targetFn(v)
		if err != nil {
			return fmt.Errorf("coercer: key %q: %w", r.key, err)
		}
		cfg[r.key] = coerced
	}
	return nil
}

func coerceToString(v interface{}) (interface{}, error) {
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
		return nil, fmt.Errorf("cannot coerce %T to string", v)
	}
}

func coerceToInt(v interface{}) (interface{}, error) {
	switch val := v.(type) {
	case int:
		return val, nil
	case float64:
		return int(val), nil
	case string:
		n, err := strconv.Atoi(val)
		if err != nil {
			return nil, fmt.Errorf("cannot coerce %q to int", val)
		}
		return n, nil
	default:
		return nil, fmt.Errorf("cannot coerce %T to int", v)
	}
}

func coerceToBool(v interface{}) (interface{}, error) {
	switch val := v.(type) {
	case bool:
		return val, nil
	case string:
		b, err := strconv.ParseBool(val)
		if err != nil {
			return nil, fmt.Errorf("cannot coerce %q to bool", val)
		}
		return b, nil
	case int:
		return val != 0, nil
	default:
		return nil, fmt.Errorf("cannot coerce %T to bool", v)
	}
}
