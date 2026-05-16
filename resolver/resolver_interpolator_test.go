package resolver_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/confmap/interpolator"
	"github.com/confmap/loader"
	"github.com/confmap/resolver"
)

func writeInterpYAML(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "cfg.yaml")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestInterpolatedResolver_ExpandsVars(t *testing.T) {
	path := writeInterpYAML(t, "host: ${MY_HOST}\nport: 8080\n")
	fl, _ := loader.NewFileLoader(path)
	interp := interpolator.New(map[string]string{"MY_HOST": "example.com"}, true)

	r := resolver.NewInterpolated([]loader.Loader{fl}, nil, interp)
	out, err := r.Resolve()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["host"] != "example.com" {
		t.Errorf("expected example.com, got %v", out["host"])
	}
}

func TestInterpolatedResolver_StrictMissingVar(t *testing.T) {
	path := writeInterpYAML(t, "host: ${UNDEFINED_VAR}\n")
	fl, _ := loader.NewFileLoader(path)
	interp := interpolator.New(map[string]string{}, true)

	r := resolver.NewInterpolated([]loader.Loader{fl}, nil, interp)
	_, err := r.Resolve()
	if err == nil {
		t.Fatal("expected error for undefined variable in strict mode")
	}
}

func TestInterpolatedResolver_NilInterpolator_Passthrough(t *testing.T) {
	path := writeInterpYAML(t, "key: ${NOT_EXPANDED}\n")
	fl, _ := loader.NewFileLoader(path)

	r := resolver.NewInterpolated([]loader.Loader{fl}, nil, nil)
	out, err := r.Resolve()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["key"] != "${NOT_EXPANDED}" {
		t.Errorf("expected raw string, got %v", out["key"])
	}
}

func TestInterpolatedResolver_ValidationFailure(t *testing.T) {
	path := writeInterpYAML(t, "name: hello\n")
	fl, _ := loader.NewFileLoader(path)
	interp := interpolator.New(nil, false)
	schema := map[string]any{
		"required": []any{"name", "missing_required"},
	}

	r := resolver.NewInterpolated([]loader.Loader{fl}, schema, interp)
	_, err := r.Resolve()
	if err == nil {
		t.Fatal("expected validation error")
	}
}
