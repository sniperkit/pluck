// +build go1.10

package jq

import (
	"bytes"
)

func decodeJSON(data []byte, v interface{}) error {
	d := json.NewDecoder(bytes.NewReader(data))
	d.DisallowUnknownFields()
	return d.Decode(v)
}
