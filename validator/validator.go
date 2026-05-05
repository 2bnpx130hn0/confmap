// Package validator provides schema validation for merged config maps.
package validator

import (
	"fmt"
	"strings"
)

// FieldType represents the expected type of a config field.
type FieldType string

const (
	TypeString  FieldType = "string"
	TypeInt     FieldType = "int"
	TypeFloat   FieldType = "float"
	TypeBool    FieldType = "bool"
	TypeMap     FieldType = "map"
)

// FieldSchema defines the schema for a single config field.
type FieldSchema struct {
	Type     FieldType
	Required bool
}

// Schema maps dot-separated field paths to their expected schema.
// Example: "database.host" -> FieldSchema{Type: TypeString, Required: true}
type Schema map[string]FieldSchema

// ValidationError collects all validation errors found.
type ValidationError struct {
	Errors []string
}

func (e *ValidationError) Error() string {
	return "validation errors:\n  " + strings.Join(e.Errors, "\n  ")
}

func (e *ValidationError) HasErrors() bool {
	return len(e.Errors) > 0
}

// Validate checks that the provided config map conforms to the given schema.
// It verifies required fields are present and that values match expected types.
func Validate(cfg map[string]any, schema Schema) error {
	ve := &ValidationError{}

	for path, fieldSchema := range schema {
		val, found := getNestedValue(cfg, strings.Split(path, "."))
		if !found {
			if fieldSchema.Required {
				ve.Errors = append(ve.Errors, fmt.Sprintf("required field %q is missing", path))
			}
			continue
		}
		if err := checkType(path, val, fieldSchema.Type); err != nil {
			ve.Errors = append(ve.Errors, err.Error())
		}
	}

	if ve.HasErrors() {
		return ve
	}
	return nil
}

func getNestedValue(cfg map[string]any, keys []string) (any, bool) {
	if len(keys) == 0 {
		return nil, false
	}
	val, ok := cfg[keys[0]]
	if !ok {
		return nil, false
	}
	if len(keys) == 1 {
		return val, true
	}
	nested, ok := val.(map[string]any)
	if !ok {
		return nil, false
	}
	return getNestedValue(nested, keys[1:])
}

func checkType(path string, val any, expected FieldType) error {
	switch expected {
	case TypeString:
		if _, ok := val.(string); !ok {
			return fmt.Errorf("field %q: expected string, got %T", path, val)
		}
	case TypeInt:
		switch val.(type) {
		case int, int64, int32:
		default:
			return fmt.Errorf("field %q: expected int, got %T", path, val)
		}
	case TypeFloat:
		switch val.(type) {
		case float32, float64:
		default:
			return fmt.Errorf("field %q: expected float, got %T", path, val)
		}
	case TypeBool:
		if _, ok := val.(bool); !ok {
			return fmt.Errorf("field %q: expected bool, got %T", path, val)
		}
	case TypeMap:
		if _, ok := val.(map[string]any); !ok {
			return fmt.Errorf("field %q: expected map, got %T", path, val)
		}
	}
	return nil
}
