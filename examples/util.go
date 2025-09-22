package examples

import (
	"bytes"
	"encoding/json"
)

// join returns a JSON object joining multiple JSON byte slices into one object.
// currently only supports joining JSON objects by merging their keys, without
// any special handling of key conflicts.
func join(parts ...[]byte) []byte {
	b := &bytes.Buffer{}
	b.WriteByte('{')
	for i, part := range parts {
		part = canonicalize(part)
		if i > 0 {
			b.WriteByte(',')
		}
		part = part[1 : len(part)-1] // strip '{' and '}'
		b.Write(part)
	}
	b.WriteByte('}')
	return b.Bytes()
}

// canonicalize returns a canonical JSON encoding of b.
func canonicalize(b []byte) []byte {
	var v any
	json.Unmarshal(b, &v)
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return b
}
