// Package stringifier converts config maps to human-readable string representations
// with configurable formatting options such as delimiter, prefix, and quoting.
package stringifier

import (
	"fmt"
	"sort"
	"strings"
)

// Options controls how the config map is rendered.
type Options struct {
	Delimiter string // key-value separator, default "="
	Prefix    string // prefix prepended to every line
	QuoteValues bool   // wrap values in double quotes
	SortKeys    bool   // output keys in sorted order
}

// Stringifier renders a config map to a slice of formatted strings.
type Stringifier struct {
	opts Options
}

// New creates a Stringifier with the provided Options.
func New(opts Options) *Stringifier {
	if opts.Delimiter == "" {
		opts.Delimiter = "="
	}
	return &Stringifier{opts: opts}
}

// Render converts cfg into a slice of "key=value" lines, flattening nested maps
// using dot notation.
func (s *Stringifier) Render(cfg map[string]any) []string {
	flat := flatten("", cfg)
	keys := make([]string, 0, len(flat))
	for k := range flat {
		keys = append(keys, k)
	}
	if s.opts.SortKeys {
		sort.Strings(keys)
	}
	lines := make([]string, 0, len(keys))
	for _, k := range keys {
		v := fmt.Sprintf("%v", flat[k])
		if s.opts.QuoteValues {
			v = `"` + v + `"`
		}
		lines = append(lines, s.opts.Prefix+k+s.opts.Delimiter+v)
	}
	return lines
}

// String returns the rendered config as a single newline-joined string.
func (s *Stringifier) String(cfg map[string]any) string {
	return strings.Join(s.Render(cfg), "\n")
}

// flatten recursively flattens nested maps using dot-separated keys.
func flatten(prefix string, m map[string]any) map[string]any {
	out := make(map[string]any)
	for k, v := range m {
		full := k
		if prefix != "" {
			full = prefix + "." + k
		}
		if nested, ok := v.(map[string]any); ok {
			for nk, nv := range flatten(full, nested) {
				out[nk] = nv
			}
		} else {
			out[full] = v
		}
	}
	return out
}
