package mxj

import (
	"strings"
)

// KeyStyle represents the specific style of the key.
type PathStyle uint

// Header style
const (
	// "/foo/bar/0/baz"
	JSONPointerStyle PathStyle = iota

	// "foo/bar/0/baz"
	SlashStyle

	// "foo.bar.0.baz"
	DotNotationStyle

	// "foo.bar[0].baz"
	DotBracketStyle
)

func (mv Map) FieldsToLower() Map {
	var data = make(map[string]interface{})
	for k, v := range mv {
		data[strings.ToLower(k)] = v
	}
	return data
}
