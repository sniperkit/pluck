package sitemap

import (
	"fmt"
	"net/url"
	"path"
	// helpers
	// pp "github.com/sniperkit/colly/plugins/app/debug/pp"
)

func getURLs(sitemapURL url.URL) ([]url.URL, error) {
	var urls []url.URL

	pathExtension := path.Ext(sitemapURL.String())

	switch pathExtension {
	case ".txt":
		urlsFromSitemap, sitemapError := getURLsFromSitemapTXT(sitemapURL)
		if sitemapError == nil {
			urls = append(urls, urlsFromSitemap...)
		}

		if isInvalidSitemapContent(sitemapError) {
			return nil, fmt.Errorf("%q is neither a sitemap index nor a XML sitemap", sitemapURL.String())
		}

	default:
		urlsFromIndex, indexError := getURLsFromSitemapIndex(sitemapURL)
		if indexError == nil {
			urls = urlsFromIndex
		}

		urlsFromSitemap, sitemapError := getURLsFromSitemap(sitemapURL)
		if sitemapError == nil {
			urls = append(urls, urlsFromSitemap...)
		}

		if isInvalidSitemapIndexContent(indexError) && isInvalidSitemapContent(sitemapError) {
			return nil, fmt.Errorf("%q is neither a sitemap index nor a XML sitemap", sitemapURL.String())
		}

	}

	return urls, nil
}

func getURLsFromSitemapTXT(txtSitemapURL url.URL) ([]url.URL, error) {
	var urls []url.URL

	sitemap, txtSitemapError := getTXTSitemap(txtSitemapURL)
	if txtSitemapError != nil {
		return nil, txtSitemapError
	}

	for _, urlEntry := range sitemap.URLs {
		/*
			parsedURL, parseError := url.Parse(urlEntry.Loc)
			if parseError != nil {
				return nil, parseError
			}
		*/
		urls = append(urls, urlEntry.href)
		// urls = append(urls, *parsedURL)
	}

	// pp.Println("urls=", urls)

	return urls, nil
}

func getURLsFromSitemap(xmlSitemapURL url.URL) ([]url.URL, error) {
	var urls []url.URL

	sitemap, xmlSitemapError := getXMLSitemap(xmlSitemapURL)
	if xmlSitemapError != nil {
		return nil, xmlSitemapError
	}

	for _, urlEntry := range sitemap.URLs {
		parsedURL, parseError := url.Parse(urlEntry.Loc)
		if parseError != nil {
			return nil, parseError
		}
		urls = append(urls, *parsedURL)
	}

	return urls, nil
}

func getURLsFromSitemapIndex(xmlSitemapURL url.URL) ([]url.URL, error) {
	var urls []url.URL

	sitemapIndex, sitemapIndexError := getSitemapIndex(xmlSitemapURL)
	if sitemapIndexError != nil {
		return nil, sitemapIndexError
	}

	for _, sitemap := range sitemapIndex.Sitemaps {
		locationURL, err := url.Parse(sitemap.Loc)
		if err != nil {
			return nil, err
		}
		sitemapUrls, err := getURLsFromSitemap(*locationURL)
		if err != nil {
			return nil, err
		}
		urls = append(urls, sitemapUrls...)
	}
	return urls, nil

}
