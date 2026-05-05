// Package resolver ties together the loader, merger, and validator packages
// to provide a single Resolve() call that produces a fully merged and validated
// configuration map.
//
// Usage:
//
//	schema := map[string]interface{}{
//		"fields": map[string]interface{}{
//			"host": map[string]interface{}{"type": "string", "required": true},
//			"port": map[string]interface{}{"type": "int", "required": true},
//		},
//	}
//
//	r := resolver.New(
//		schema,
//		resolver.Source{Name: "defaults", Loader: loader.NewFileLoader("defaults.yaml")},
//		resolver.Source{Name: "env",      Loader: loader.NewEnvLoader("APP_")},
//	)
//
//	cfg, err := r.Resolve()
//	if err != nil {
//		log.Fatal(err)
//	}
package resolver
