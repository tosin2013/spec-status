package schema

import (
	"bytes"
	"encoding/json"
)

// fallback is printed if a report cannot be serialized at all.
const fallback = "{\"schema_version\":\"1.0\",\"outcome\":\"fatal\",\"complete\":false,\"source\":{\"roadmap\":null,\"input_format\":null},\"configuration\":{\"fail_on\":\"error\",\"dependency_policy\":\"verified\",\"task_scope\":\"unique_linked_task_files\"},\"summary\":null,\"items\":[],\"diagnostics\":[{\"code\":\"SS_INVOCATION_INVALID\",\"severity\":\"error\",\"message\":\"Report could not be serialized.\",\"source\":null,\"item_id\":null,\"related\":[]}],\"error\":{\"code\":\"SS_INVOCATION_INVALID\",\"message\":\"Report could not be serialized.\",\"source\":null}}\n"

// Marshal returns the exact bytes `spec status --json` prints: one UTF-8
// JSON document without HTML escaping, and a trailing newline.
// Nil Items, Diagnostics, Dependencies, Related and StatusCounts are written
// as [] or {} as the schema requires; a nil Tasks.Phases stays null, which
// means the tasks file was not read. r itself is not modified.
func Marshal(r Report) []byte {
	r = normalized(r)
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(r); err != nil {
		return []byte(fallback)
	}
	return buf.Bytes()
}

func normalized(r Report) Report {
	if r.Items == nil {
		r.Items = []Item{}
	} else {
		items := make([]Item, len(r.Items))
		copy(items, r.Items)
		for i := range items {
			if items[i].Dependencies == nil {
				items[i].Dependencies = []Dependency{}
			}
		}
		r.Items = items
	}
	if r.Diagnostics == nil {
		r.Diagnostics = []Diagnostic{}
	} else {
		diags := make([]Diagnostic, len(r.Diagnostics))
		copy(diags, r.Diagnostics)
		for i := range diags {
			if diags[i].Related == nil {
				diags[i].Related = []Ref{}
			}
		}
		r.Diagnostics = diags
	}
	if r.Summary != nil && r.Summary.StatusCounts == nil {
		s := *r.Summary
		s.StatusCounts = map[string]int{}
		r.Summary = &s
	}
	return r
}
