package resolver_test

import (
	"errors"
	"testing"

	"github.com/iamBijoyKar/confmap/aliaser"
	"github.com/iamBijoyKar/confmap/resolver"
)

type aliasStaticLoader struct{ data map[string]any }

func (l *aliasStaticLoader) Load() (map[string]any, error) { return l.data, nil }

type aliasErrLoader struct{}

func (l *aliasErrLoader) Load() (map[string]any, error) {
	return nil, errors.New("load error")
}

func TestAliasedResolver_ExpandsAlias(t *testing.T) {
	a := aliaser.New()
	_ = a.Register("db_host", "database.host")

	loaders := []resolver.Loader{
		&aliasStaticLoader{data: map[string]any{"db_host": "localhost", "port": 5432}},
	}
	ar, err := resolver.NewAliased(loaders, nil, a)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	cfg, err := ar.Resolve()
	if err != nil {
		t.Fatalf("resolve error: %v", err)
	}
	if cfg["database.host"] != "localhost" {
		t.Errorf("expected database.host=localhost, got %v", cfg["database.host"])
	}
	if _, ok := cfg["db_host"]; ok {
		t.Error("alias key should have been removed")
	}
}

func TestAliasedResolver_NilAliaser(t *testing.T) {
	_, err := resolver.NewAliased(nil, nil, nil)
	if err == nil {
		t.Error("expected error for nil aliaser")
	}
}

func TestAliasedResolver_LoaderError(t *testing.T) {
	a := aliaser.New()
	loaders := []resolver.Loader{&aliasErrLoader{}}
	ar, err := resolver.NewAliased(loaders, nil, a)
	if err != nil {
		t.Fatalf("unexpected setup error: %v", err)
	}
	_, err = ar.Resolve()
	if err == nil {
		t.Error("expected loader error to propagate")
	}
}

func TestAliasedResolver_NoAliases(t *testing.T) {
	a := aliaser.New()
	loaders := []resolver.Loader{
		&aliasStaticLoader{data: map[string]any{"key": "value"}},
	}
	ar, _ := resolver.NewAliased(loaders, nil, a)
	cfg, err := ar.Resolve()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg["key"] != "value" {
		t.Error("config should be unchanged when no aliases registered")
	}
}
