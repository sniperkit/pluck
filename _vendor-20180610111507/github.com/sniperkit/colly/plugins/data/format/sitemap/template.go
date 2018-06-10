package sitemap

const (
	header = `<?xml version="1.0" encoding="UTF-8"?>
	<urlset xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="http://www.sitemaps.org/schemas/sitemap/0.9 http://www.sitemaps.org/schemas/sitemap/0.9/sitemap.xsd" xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`
	footer = `
	</urlset>`
	template = `
	<url>
	  <loc>%s</loc>
	  <lastmod>%s</lastmod>
	  <changefreq>%s</changefreq>
	  <priority>%.1f</priority>
	</url> 	`

	indexHeader = `<?xml version="1.0" encoding="UTF-8"?>
  <sitemapindex xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`
	indexFooter = `
	</sitemapindex>`
	indexTemplate = `
	<sitemap>
		<loc>%s%s</loc>
		<lastmod>%s</lastmod>
	</sitemap>`
)

// https://github.com/pengux/sitemap
// https://github.com/henrybear327/CCU-Search-Engine/blob/master/assignment2/version2/preprocess.go
const (
	// MaxSitemapItems is the maximum number of items for a single sitemap
	MaxSitemapItems = 50000

	// SitemapXML is the XML structure for urlset in sitemaps
	SitemapXML = `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
	xsi:schemaLocation="http://www.sitemaps.org/schemas/sitemap/0.9 http://www.sitemaps.org/schemas/sitemap/0.9/sitemap.xsd"
	xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">%s
</urlset>`

	// SitemapItemXML is the XML format for the URL item in sitemap
	SitemapItemXML = `
	<url>
		<loc>%s</loc>
		<lastmod>%s</lastmod>
		<changefreq>%s</changefreq>
		<priority>%.1f</priority>
	</url>`

	// SitemapIndexXML is the XML structure of a sitemap index
	SitemapIndexXML = `<?xml version="1.0" encoding="UTF-8"?>
<sitemapindex xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">%s
</sitemapindex>
`

	// SitemapIndexItemXML is the XML structure of a sitemap index item
	SitemapIndexItemXML = `
	<sitemap>
		<loc>%s</loc>
		<lastmod>%s</lastmod>
	</sitemap>`
)
