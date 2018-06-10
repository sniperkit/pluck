package sitemap

import (
	"fmt"
	"net/url"

	"github.com/sniperkit/colly/pkg"
	"github.com/sniperkit/colly/pkg/queue"
)

func NewWithCollector(inputURL string, c *colly.Collector) (*SitemapCollector, error) {
	cs, err := New(inputURL)
	if err != nil {
		return nil, err
	}
	cs.collector = c
	return cs, nil
}

func AttachCollector(inputURL string, c *colly.Collector) (*SitemapCollector, error) {
	cs, err := New(inputURL)
	if err != nil {
		return nil, err
	}
	cs.collector = c
	return cs, nil
}

func AttachQueue(inputURL string, cqueue *queue.Queue) (*SitemapCollector, error) {
	cs, err := New(inputURL)
	if err != nil {
		return nil, err
	}
	cs.cqueue = cqueue
	return cs, nil
}

func AttachCollectorWithQueue(inputURL string, collector *colly.Collector, cqueue *queue.Queue) (*SitemapCollector, error) {
	cs, err := New(inputURL)
	if err != nil {
		return nil, err
	}
	cs.collector = collector
	cs.cqueue = cqueue
	return cs, nil
}

func (cs *SitemapCollector) Count() int {
	return len(cs.URLs)
}

func (cs *SitemapCollector) List() ([]url.URL, []string) {
	cs.URLs, _ = getURLs(cs.href)
	var urls []string
	for _, u := range cs.URLs {
		urls = append(urls, u.String())
	}
	return cs.URLs, urls
}

func (cs *SitemapCollector) Index() ([]url.URL, []string) {
	var sitemaps []string
	for _, sitemap := range cs.Indices {
		sitemaps = append(sitemaps, sitemap.String())
	}
	return cs.Indices, sitemaps
}

func (cs *SitemapCollector) Sitemaps() []url.URL {
	return cs.Indices
}

func (cs *SitemapCollector) getURLs() {
	if !cs.IsValid() {
		return
	}
	// var urls []url.URL
	urlsFromIndex, indexError := getURLsFromSitemapIndex(cs.href)
	if indexError == nil {
		cs.Indices = urlsFromIndex
	}

	urlsFromSitemap, sitemapError := getURLsFromSitemap(cs.href)
	if sitemapError == nil {
		cs.URLs = append(cs.URLs, urlsFromSitemap...)
	}

	// if isInvalidSitemapIndexContent(indexError) && isInvalidXMLSitemapContent(sitemapError) {
	// 	return nil, fmt.Errorf("%q is neither a sitemap index nor a XML sitemap", xmlSitemapURL.String())
	// }
}

func (cs *SitemapCollector) All() ([]url.URL, error) {
	if !cs.IsValid() {
		return nil, fmt.Errorf("sitemap at url='%q' is not reachable", cs.href.String())
	}
	var urls []url.URL
	urlsFromIndex, indexError := getURLsFromSitemapIndex(cs.href)
	if indexError == nil {
		urls = urlsFromIndex
		cs.Indices = urlsFromIndex
	}

	urlsFromSitemap, sitemapError := getURLsFromSitemap(cs.href)
	if sitemapError == nil {
		urls = append(urls, urlsFromSitemap...)
		cs.URLs = append(cs.URLs, urlsFromSitemap...)
	}

	if isInvalidSitemapIndexContent(indexError) && isInvalidSitemapContent(sitemapError) {
		return nil, fmt.Errorf("%q is neither a sitemap index nor a XML sitemap", cs.href.String())
	}
	return urls, nil
}

func (cs *SitemapCollector) VisitAll() *colly.Collector {
	if !cs.IsValid() {
		return cs.collector
	}
	if log != nil {
		log.Println("sitemapURL=", cs.href.String())
	}
	links, err := getURLs(cs.href)
	if err != nil {
		if log != nil {
			log.Fatalln("error: ", err)
		}
		return cs.collector
	}
	if log != nil {
		log.Println("links found:", len(links))
	}
	for _, link := range links {
		if log != nil {
			log.Println("add -", link.String())
		}
		cs.collector.Visit(fmt.Sprintf("%s", link.String())) // Request visit URL by Collector
	}
	return cs.collector
}

func (cs *SitemapCollector) EnqueueAll() {
	if !cs.IsValid() {
		return
	}
	links, err := getURLs(cs.href)
	if err != nil {
		if log != nil {
			log.Fatalln("error: ", err)
		}
		return
	}
	if log != nil {
		log.Println("links found:", len(links))
	}
	for _, link := range links {
		if log != nil {
			log.Println("enqueue -", link.String())
		}
		cs.cqueue.AddURL(fmt.Sprintf("%s", link.String())) // Enqueue new URL
	}
	return
}

func VisitAll(inputURL string, c *colly.Collector) *colly.Collector {
	sitemapURL, err := url.Parse(inputURL)
	if err != nil {
		if log != nil {
			log.Fatalln("error: ", err)
		}
		return c
	}

	links, err := getURLs(*sitemapURL)
	if err != nil {
		if log != nil {
			log.Fatalln("error: ", err)
		}
		return c
	}
	if log != nil {
		log.Println("links found:", len(links))
	}

	for _, link := range links {
		if log != nil {
			log.Println("add -", link.String())
		}
		c.Visit(fmt.Sprintf("%s", link.String())) // Request visit URL by Collector
	}
	return c
}

func EnqueueAll(inputURL string, q *queue.Queue) *queue.Queue {
	sitemapURL, err := url.Parse(inputURL)
	if err != nil {
		if log != nil {
			log.Fatalln("error: ", err)
		}
		return q
	}
	if log != nil {
		log.Println("sitemapURL=", sitemapURL)
	}
	links, err := getURLs(*sitemapURL)
	if err != nil {
		if log != nil {
			log.Fatalln("error: ", err)
		}
		return q
	}
	if log != nil {
		log.Println("links found:", len(links))
	}
	for _, link := range links {
		if log != nil {
			log.Println("enqueue -", link.String())
		}
		q.AddURL(fmt.Sprintf("%s", link.String())) // Enqueue new URL
	}
	return q
}
