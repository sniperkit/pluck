// +build !go1.10

package jq

func decodeJSON(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
