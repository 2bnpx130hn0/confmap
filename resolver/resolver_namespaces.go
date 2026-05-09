package resolver

import (
	"github.com/your-org/confmap/namespaces"
)

// NamespacedResolver wraps a Resolver and exposes a scoped view of the
// resolved config under a fixed prefix.
type NamespacedResolver struct {
	*Resolver
	prefix string
}

// NewNamespaced returns a NamespacedResolver that restricts the resolved
// config to the subtree identified by prefix (e.g. "database").
func NewNamespaced(r *Resolver, prefix string) *NamespacedResolver {
	return &NamespacedResolver{Resolver: r, prefix: prefix}
}

// Namespace returns a live Namespace view of the current resolved config
// scoped to the resolver's prefix. The returned Namespace references the
// underlying data directly; call Resolver.Resolve again to refresh.
func (nr *NamespacedResolver) Namespace() (*namespaces.Namespace, error) {
	cfg, err := nr.Resolve()
	if err != nil {
		return nil, err
	}
	return namespaces.New(nr.prefix, cfg), nil
}

// Get is a convenience wrapper that resolves the config and retrieves a
// single namespaced key in one call.
func (nr *NamespacedResolver) Get(local string) (any, error) {
	ns, err := nr.Namespace()
	if err != nil {
		return nil, err
	}
	v, ok := ns.Get(local)
	if !ok {
		return nil, nil
	}
	return v, nil
}
