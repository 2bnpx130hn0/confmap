// Package namespaces provides scoped views into a config map,
// allowing callers to read and write keys under a fixed prefix.
package namespaces

import (
	"fmt"
	"strings"
)

// Namespace wraps a flat config map and exposes a scoped view
// under a given prefix, e.g. "database" or "server".
type Namespace struct {
	prefix string
	data   map[string]any
}

// New returns a Namespace that scopes all operations under prefix.
// The supplied data map is referenced directly (not copied).
func New(prefix string, data map[string]any) *Namespace {
	return &Namespace{prefix: prefix, data: data}
}

// key returns the fully-qualified key for a given local name.
func (n *Namespace) key(local string) string {
	if n.prefix == "" {
		return local
	}
	return n.prefix + "." + local
}

// Get retrieves the value at the namespaced key.
// It traverses nested maps using dot-separated segments.
func (n *Namespace) Get(local string) (any, bool) {
	parts := strings.Split(n.key(local), ".")
	var current any = n.data
	for _, p := range parts {
		m, ok := current.(map[string]any)
		if !ok {
			return nil, false
		}
		current, ok = m[p]
		if !ok {
			return nil, false
		}
	}
	return current, true
}

// Set writes value at the namespaced key, creating intermediate
// maps as needed. Returns an error if an intermediate segment is
// occupied by a non-map value.
func (n *Namespace) Set(local string, value any) error {
	parts := strings.Split(n.key(local), ".")
	current := n.data
	for _, p := range parts[:len(parts)-1] {
		v, exists := current[p]
		if !exists {
			next := map[string]any{}
			current[p] = next
			current = next
			continue
		}
		next, ok := v.(map[string]any)
		if !ok {
			return fmt.Errorf("namespaces: segment %q is not a map", p)
		}
		current = next
	}
	current[parts[len(parts)-1]] = value
	return nil
}

// Keys returns all leaf-level local keys visible under this namespace.
func (n *Namespace) Keys() []string {
	root, ok := n.Get("")
	if n.prefix == "" {
		root = n.data
		ok = true
	}
	if !ok {
		return nil
	}
	m, ok := root.(map[string]any)
	if !ok {
		return nil
	}
	var keys []string
	collect(m, "", &keys)
	return keys
}

func collect(m map[string]any, prefix string, out *[]string) {
	for k, v := range m {
		full := k
		if prefix != "" {
			full = prefix + "." + k
		}
		if nested, ok := v.(map[string]any); ok {
			collect(nested, full, out)
		} else {
			*out = append(*out, full)
		}
	}
}
