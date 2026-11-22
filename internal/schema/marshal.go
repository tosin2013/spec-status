package schema

import (
	"bytes"
	"encoding/json"
)

// fallback is printed if a report cannot be serialized at all.
const fallback = "{\"schema_version\":\"1.0\",\"outcome\":\"fatal\",\"complete\":false,\"source\":{\"roadmap\":null,\"input_format\":null},\"configuration\":{\"fail_on\":\"error\",\"dependency_policy\":\"verified\",\"task_scope\":\"unique_linked_task_files\"},\"summary\":null,\"items\":[],\"diagnostics\":[{\"code\":\"SS_INVOCATION_INVALID\",\"severity\":\"error\",\"message\":\"Report could not be serialized.\",\"source\":null,\"item_id\":null,\"related\":[]}],\"error\":{\"code\":\"SS_INVOCATION_INVALID\",\"message\":\"Report could not be serialized.\",\"source\":null}}\n"

// Marshal returns the exact bytes `spec status --json` prints: one UTF-8
// JSON document without HTML escaping, and a trailing newline.
func Marshal(r Report) []byte {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(r); err != nil {
		return []byte(fallback)
	}
	return buf.Bytes()
}
