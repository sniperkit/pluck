package main

import (
	tablib "github.com/sniperkit/colly/plugins/data/transform/tabular"
)

// data structure - databooks and datasets
var (
	sheets   map[string][]interface{} = make(map[string][]interface{}, 0)
	dsExport *tablib.Dataset
	dsURLs   *tablib.Dataset
	dataBook *tablib.Databook
)

func initDataCollections() {
	dsURLs = tablib.NewDataset([]string{"loc", "changefreq", "priority"})
	dsExport = tablib.NewDataset([]string{"url", "price", "size", "colors"})
	dataBook = tablib.NewDatabook()
	dataBook.AddSheet("patterns", dsExport) //.Sort("price"))         // add the patterns sheets to the collector databook
	dataBook.AddSheet("known_urls", dsURLs) //.Sort("priority")) // add the kown_urls sheets to the collector databook
}
