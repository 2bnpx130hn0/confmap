// Package aliaser provides key aliasing for config maps,
// allowing alternate key names to be resolved transparently.
package aliaser

import "fmt"

// Aliaser maps alias keys to canonical keys in a config map.
type Aliaser struct {
	aliases map[string]string // alias -> canonical
}

// New creates a new Aliaser with no aliases registered.
func New() *Aliaser {
	return &Aliaser{aliases: make(map[string]string)}
}

// Register adds an alias for a canonical key.
// Reading from alias will return the value of canonical.
func (a *Aliaser) Register(alias, canonical string) error {
	if alias == "" || canonical == "" {
		return fmt.Errorf("aliaser: alias and canonical must be non-empty")
	}
	if alias == canonical {
		return fmt.Errorf("aliaser: alias %q must differ from canonical", alias)
	}
	a.aliases[alias] = canonical
	return nil
}

// Resolve returns the canonical key for a given key.
// If the key is not an alias, it is returned unchanged.
func (a *Aliaser) Resolve(key string) string {
	if canon, ok := a.aliases[key]; ok {
		return canon
	}
	return key
}

// Apply returns a new config map where all alias keys are expanded
// to their canonical form. Canonical keys already present take precedence.
func (a *Aliaser) Apply(cfg map[string]any) (map[string]any, error) {
	if cfg == nil {
		return nil, nil
	}
	out := make(map[string]any, len(cfg))
	for k, v := range cfg {
		out[k] = v
	}
	for alias, canonical := range a.aliases {
		val, hasAlias := out[alias]
		if !hasAlias {
			continue
		}
		if _, hasCanon := out[canonical]; !hasCanon {
			out[canonical] = val
		}
		delete(out, alias)
	}
	return out, nil
}

// Aliases returns a copy of the registered alias map.
func (a *Aliaser) Aliases() map[string]string {
	copy := make(map[string]string, len(a.aliases))
	for k, v := range a.aliases {
		copy[k] = v
	}
	return copy
}
