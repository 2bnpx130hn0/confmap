package merger

import "fmt"

// PatchOp represents a single patch operation.
type PatchOp struct {
	Op    string // "set", "delete", "rename"
	Path  string
	Value interface{}
	To    string // used for rename
}

// Patcher applies a sequence of patch operations to a base config map.
type Patcher struct {
	ops []PatchOp
}

// NewPatch creates a new Patcher with the given operations.
func NewPatch(ops []PatchOp) *Patcher {
	return &Patcher{ops: ops}
}

// Merge applies all patch operations to the base config and returns the result.
// The base config is not mutated.
func (p *Patcher) Merge(base map[string]interface{}) (map[string]interface{}, error) {
	out := DeepCopy(base)
	for _, op := range p.ops {
		switch op.Op {
		case "set":
			out[op.Path] = op.Value
		case "delete":
			delete(out, op.Path)
		case "rename":
			if op.To == "" {
				return nil, fmt.Errorf("merger/patch: rename op missing 'to' field for path %q", op.Path)
			}
			val, ok := out[op.Path]
			if !ok {
				return nil, fmt.Errorf("merger/patch: rename source key %q not found", op.Path)
			}
			out[op.To] = val
			delete(out, op.Path)
		default:
			return nil, fmt.Errorf("merger/patch: unknown op %q", op.Op)
		}
	}
	return out, nil
}
