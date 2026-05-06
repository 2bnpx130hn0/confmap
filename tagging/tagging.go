// Package tagging provides key-level tag annotation for config maps.
// Tags are arbitrary string labels attached to config keys, enabling
// downstream filtering, auditing, or documentation generation.
package tagging

import "fmt"

// Tagger holds tag annotations for config keys.
type Tagger struct {
	tags map[string][]string
}

// New returns a new Tagger instance.
func New() *Tagger {
	return &Tagger{tags: make(map[string][]string)}
}

// Tag adds one or more tags to a config key.
func (t *Tagger) Tag(key string, tags ...string) {
	t.tags[key] = append(t.tags[key], tags...)
}

// Tags returns all tags associated with a key.
func (t *Tagger) Tags(key string) []string {
	return t.tags[key]
}

// HasTag reports whether a key carries a specific tag.
func (t *Tagger) HasTag(key, tag string) bool {
	for _, v := range t.tags[key] {
		if v == tag {
			return true
		}
	}
	return false
}

// FilterByTag returns config entries whose keys carry the given tag.
func (t *Tagger) FilterByTag(cfg map[string]any, tag string) map[string]any {
	out := make(map[string]any)
	for k, v := range cfg {
		if t.HasTag(k, tag) {
			out[k] = v
		}
	}
	return out
}

// Annotate returns a copy of cfg with an extra "__tags__" key containing
// a map of key → tags for all tagged keys present in cfg.
func (t *Tagger) Annotate(cfg map[string]any) map[string]any {
	out := make(map[string]any, len(cfg)+1)
	for k, v := range cfg {
		out[k] = v
	}
	annotations := make(map[string][]string)
	for k, tags := range t.tags {
		if _, ok := cfg[k]; ok {
			annotations[k] = tags
		}
	}
	if len(annotations) > 0 {
		out["__tags__"] = annotations
	}
	return out
}

// Remove deletes all tags for a key.
func (t *Tagger) Remove(key string) error {
	if _, ok := t.tags[key]; !ok {
		return fmt.Errorf("tagging: key %q not found", key)
	}
	delete(t.tags, key)
	return nil
}
