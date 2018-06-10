package main

import (
	"log"

	// colly core
	"github.com/sniperkit/colly/pkg"
	"github.com/sniperkit/colly/pkg/debug"
	"github.com/sniperkit/colly/pkg/helper"
	"github.com/sniperkit/colly/pkg/queue"

	// colly plugins
	sitemap "github.com/sniperkit/colly/plugins/data/format/sitemap"

	// datastructure helpers
	cmmap "github.com/sniperkit/colly/plugins/data/structure/map/multi"
	tablib "github.com/sniperkit/colly/plugins/data/transform/tabular"
)

const (
	SITEMAP_URL          string = "https://www.shopify.com/sitemap.xml"
	SITEMAP_URL_GZ       string = "http://www.nytimes.com/sitemaps/sitemap_news/sitemap.xml.gz"
	SITEMAP_URL_TXT      string = "https://golanglibs.com/sitemap.txt"
	SITEMAP_INDEX_URL    string = "https://www.coindesk.com/sitemap_index.xml"
	SITEMAP_INDEX_URL_GZ string = "http://www.lidl.de/sitemap_index.xml.gz"
)

var (
	version           string                   = "0.0.1-alpha"
	cacheCollectorDir string                   = "./shared/cache/collector"
	sheets            map[string][]interface{} = make(map[string][]interface{}, 0)
	dsExport          *tablib.Dataset
	dsURLs            *tablib.Dataset
	dataBook          *tablib.Databook
	mapKnownURLs      = cmmap.NewConcurrentMultiMap()
	logger            *log.Logger
	cq                *queue.Queue
)

func init() {
	initCollections()
}

func initCollections() {
	dsURLs = tablib.NewDataset([]string{"loc", "changefreq", "priority"})
	dsExport = tablib.NewDataset([]string{"url", "price", "size", "colors"})
	dataBook = tablib.NewDatabook()
	dataBook.AddSheet("patterns", dsExport) //.Sort("price"))         // add the patterns sheets to the collector databook
	dataBook.AddSheet("known_urls", dsURLs) //.Sort("priority")) // add the kown_urls sheets to the collector databook
}

func main() {

	// Create a Collector
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36"),
		colly.AllowedDomains("www.shopify.com"),
		colly.Debugger(&debug.LogDebugger{}), // Attach a debugger to the collector
		colly.Async(true),
		colly.CacheDir(cacheCollectorDir), // Cache responses to prevent multiple download of pages even if the collector is restarted
	)

	// Limit the number of threads started by colly to two
	// when visiting links which domains' matches "*httpbin.*" glob
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 2,
	})

	helper.RandomUserAgent(c)
	helper.Referrer(c)

	cs, err := sitemap.NewWithCollector(SITEMAP_INDEX_URL_GZ, c)
	if err != nil {
		log.Println("invalid sitemap.")
	}

	if cq != nil {
		cs.EnqueueAll()
	} else {
		cs.VisitAll()
	}
	cs.Count()

	c.Wait()

}
