package validator_test

import (
	"testing"

	"github.com/yourorg/confmap/validator"
)

func TestValidate_AllRequiredPresent(t *testing.T) {
	cfg := map[string]any{
		"host": "localhost",
		"port": 5432,
		"database": map[string]any{
			"name": "mydb",
		},
	}
	schema := validator.Schema{
		"host":          {Type: validator.TypeString, Required: true},
		"port":          {Type: validator.TypeInt, Required: true},
		"database.name": {Type: validator.TypeString, Required: true},
	}
	if err := validator.Validate(cfg, schema); err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func TestValidate_MissingRequiredField(t *testing.T) {
	cfg := map[string]any{
		"host": "localhost",
	}
	schema := validator.Schema{
		"host": {Type: validator.TypeString, Required: true},
		"port": {Type: validator.TypeInt, Required: true},
	}
	err := validator.Validate(cfg, schema)
	if err == nil {
		t.Fatal("expected validation error for missing required field")
	}
	ve, ok := err.(*validator.ValidationError)
	if !ok {
		t.Fatalf("expected *ValidationError, got %T", err)
	}
	if len(ve.Errors) != 1 {
		t.Errorf("expected 1 error, got %d: %v", len(ve.Errors), ve.Errors)
	}
}

func TestValidate_WrongType(t *testing.T) {
	cfg := map[string]any{
		"enabled": "yes", // should be bool
	}
	schema := validator.Schema{
		"enabled": {Type: validator.TypeBool, Required: true},
	}
	err := validator.Validate(cfg, schema)
	if err == nil {
		t.Fatal("expected type mismatch error")
	}
}

func TestValidate_OptionalFieldAbsent(t *testing.T) {
	cfg := map[string]any{
		"host": "localhost",
	}
	schema := validator.Schema{
		"host":    {Type: validator.TypeString, Required: true},
		"timeout": {Type: validator.TypeInt, Required: false},
	}
	if err := validator.Validate(cfg, schema); err != nil {
		t.Errorf("optional missing field should not cause error, got: %v", err)
	}
}

func TestValidate_NestedFieldTypeMismatch(t *testing.T) {
	cfg := map[string]any{
		"database": map[string]any{
			"port": "not-a-number",
		},
	}
	schema := validator.Schema{
		"database.port": {Type: validator.TypeInt, Required: true},
	}
	err := validator.Validate(cfg, schema)
	if err == nil {
		t.Fatal("expected type mismatch error for nested field")
	}
}

func TestValidate_EmptyConfig(t *testing.T) {
	cfg := map[string]any{}
	schema := validator.Schema{
		"host": {Type: validator.TypeString, Required: true},
	}
	err := validator.Validate(cfg, schema)
	if err == nil {
		t.Fatal("expected error for missing required field in empty config")
	}
}
