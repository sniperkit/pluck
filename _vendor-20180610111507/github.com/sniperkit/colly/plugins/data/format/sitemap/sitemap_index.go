package sitemap

import (
	"encoding/xml"
	"net/url"
	"strings"
)

func getSitemapIndex(xmlSitemapURL url.URL) (Indices, error) {
	response, readErr := readURL(xmlSitemapURL)
	if readErr != nil {
		return Indices{}, readErr
	}
	if !strings.Contains(string(response.GetBody()), "</sitemapindex>") {
		return Indices{}, IndexError{"Invalid content"}
	}
	var sitemapIndex Indices
	unmarshalError := xml.Unmarshal(response.GetBody(), &sitemapIndex)
	if unmarshalError != nil {
		return Indices{}, unmarshalError
	}
	return sitemapIndex, nil
}

func (sitemapIndexError IndexError) Error() string {
	return sitemapIndexError.message
}

func isInvalidSitemapIndexContent(err error) bool {
	if err == nil {
		return false
	}

	return err.Error() == "Invalid content"
}
