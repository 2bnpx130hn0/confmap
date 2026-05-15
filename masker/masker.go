// Package masker provides utilities for partially masking config values,
// replacing sensitive string content with obfuscated representations while
// preserving a configurable number of visible characters.
package masker

import "strings"

// Masker partially obscures string values in a config map for keys that
// match a set of sensitive key names.
type Masker struct {
	keys    map[string]struct{}
	visible int    // number of trailing characters to keep visible
	maskCh  string // character used for masking
}

// New creates a Masker that will mask values for the given key names.
// visible controls how many trailing characters remain unmasked.
// maskCh is the character repeated for the hidden portion (e.g. "*").
func New(sensitiveKeys []string, visible int, maskCh string) *Masker {
	if maskCh == "" {
		maskCh = "*"
	}
	km := make(map[string]struct{}, len(sensitiveKeys))
	for _, k := range sensitiveKeys {
		km[strings.ToLower(k)] = struct{}{}
	}
	return &Masker{keys: km, visible: visible, maskCh: maskCh}
}

// Apply returns a deep copy of cfg with sensitive string values masked.
func (m *Masker) Apply(cfg map[string]any) map[string]any {
	if cfg == nil {
		return nil
	}
	out := make(map[string]any, len(cfg))
	for k, v := range cfg {
		switch val := v.(type) {
		case map[string]any:
			out[k] = m.Apply(val)
		case string:
			if _, sensitive := m.keys[strings.ToLower(k)]; sensitive {
				out[k] = m.maskString(val)
			} else {
				out[k] = val
			}
		default:
			out[k] = v
		}
	}
	return out
}

// IsSensitive reports whether key is in the sensitive key set.
func (m *Masker) IsSensitive(key string) bool {
	_, ok := m.keys[strings.ToLower(key)]
	return ok
}

func (m *Masker) maskString(s string) string {
	if len(s) == 0 {
		return s
	}
	keep := m.visible
	if keep > len(s) {
		keep = len(s)
	}
	hidden := len(s) - keep
	if hidden <= 0 {
		return strings.Repeat(m.maskCh, len(s))
	}
	return strings.Repeat(m.maskCh, hidden) + s[hidden:]
}
