package colly

import (
	tabular "github.com/sniperkit/colly/plugins/data/transform/tabular"
)

/*
	Notes:

	- Export formats supported:
	  - JSON (Sets + Books)
	  - YAML (Sets + Books)
	  - XLSX (Sets + Books)
	  - XML (Sets + Books)
	  - TSV (Sets)
	  - CSV (Sets)
	  - ASCII + Markdown (Sets)
	  - MySQL (Sets)
	  - Postgres (Sets)

	Loading formats supported:
	  - JSON (Sets + Books)
	  - YAML (Sets + Books)
	  - XML (Sets)
	  - CSV (Sets)
	  - TSV (Sets)

*/

// TABElement is the representation of a TAB tag.
type TABElement struct {
	// Name is the name of the tag
	Name string

	// Dataset
	Dataset *tabular.Dataset

	// Text
	Text string

	// Request is the request object of the element's HTML document
	Request *Request

	// Response is the Response object of the element's HTML document
	Response *Response

	// err stores the loading error
	err error
}

// NewTABElementFromTABNode creates a TABElement from a jsonquery.Node.
func NewTABElementFromTABNode(resp *Response, query string, ds *tabular.Dataset) *TABElement {

	// new TABElement object
	t := &TABElement{
		Dataset:  ds,
		Name:     query,
		Request:  resp.Request,
		Response: resp,
	}

	return t
}

// Slice
func (h *TABElement) Slice(pluckerConfig map[string]interface{}) string {
	return ""
}

// Select
func (h *TABElement) Select(extractorConfig map[string]interface{}) string {
	return ""
}
