// Package namespaces provides scoped, prefix-based access to a shared
// config map.
//
// A Namespace isolates a subtree of the configuration so that
// components can read and write their own keys without being aware of
// the full config structure.
//
// Example:
//
//	data := map[string]any{
//		"database": map[string]any{
//			"host": "localhost",
//			"port": 5432,
//		},
//	}
//
//	ns := namespaces.New("database", data)
//	host, _ := ns.Get("host")   // returns "localhost"
//	_ = ns.Set("pool", 5)       // writes data["database"]["pool"] = 5
package namespaces
