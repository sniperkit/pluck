// Copyright 2018 Adam Tauber
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package colly

import (
	"strings"

	// internal
	jsonquery "github.com/sniperkit/colly/plugins/data/extract/query/json"
)

type JsonParser string

const (
	MXJ        JsonParser = "mxj"        // https://github.com/clbanning/mxj
	GABS       JsonParser = "gabs"       // https://github.com/Jeffail/gabs
	GJSON      JsonParser = "gjson"      // https://github.com/tidwall/gjson
	LAZYJSON   JsonParser = "lazyjson"   // https://github.com/qw4990/lazyjson
	FASTJSON   JsonParser = "fastjson"   // https://github.com/valyala/fastjson
	FFJSON     JsonParser = "ffjson"     // https://github.com/pquerna/ffjson
	EASYJSON   JsonParser = "easyjson"   // https://github.com/mailru/easyjson
	JSONPARSER JsonParser = "jsonparser" // https://github.com/buger/jsonparser
	DJSON      JsonParser = "djson"      // https://github.com/a8m/djson
	JSNM       JsonParser = "jsnm"       // https://github.com/toukii/jsnm
	JSONSTREAM JsonParser = "jsonstream" // https://github.com/pb-/jsonstream
	JSONEZ     JsonParser = "jsonez"     // https://github.com/srikanth2212/jsonez
	JSON       JsonParser = "json"       // encoding/json
)

// JSONElement is the representation of a JSON tag.
type JSONElement struct {
	// Name is the name of the tag
	Name string
	// Text is the content node
	Text string
	// Request is the request object of the element's HTML document
	Request *Request
	// Response is the Response object of the element's HTML document
	Response *Response
	// DOM is the DOM object of the page. DOM is relative
	// to the current JSONElement and is either a html.Node or jsonquery.Node
	// based on how the JSONElement was created.
	DOM interface{}
}

// NewJSONElementFromJSONNode creates a JSONElement from a jsonquery.Node.
func NewJSONElementFromJSONNode(resp *Response, s *jsonquery.Node) *JSONElement {
	return &JSONElement{
		Name:     s.Data,
		Request:  resp.Request,
		Response: resp,
		Text:     s.InnerText(),
		DOM:      s,
	}
}

// Extract
func (h *JSONElement) Extract(pluckerConfig map[string]interface{}) string {
	return ""
}

// Header
func (h *JSONElement) Header(key string) (value string) {
	value = strings.TrimSpace(h.Response.Headers.Get(key))
	return
}

// Headers
func (h *JSONElement) Headers() map[string]string {
	res := make(map[string]string, 0)
	/*
		for key, val := range h.Response.Headers {
			res[key] = val
		}
	*/
	return res
}

// FindOne
func (h *JSONElement) FindOne(xpathQuery string) string {

	n := jsonquery.FindOne(h.DOM.(*jsonquery.Node), xpathQuery)
	if n == nil {
		return ""
	}
	return strings.TrimSpace(n.InnerText())
}

// Find
func (h *JSONElement) Find(xpathQuery, attrName string) []string {
	var res []string
	child := jsonquery.Find(h.DOM.(*jsonquery.Node), xpathQuery)
	for _, node := range child {
		res = append(res, node.InnerText())
	}
	return res
}
