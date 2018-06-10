package main

import (
	"log"

	"github.com/sniperkit/colly/pkg"
	"github.com/sniperkit/colly/pkg/helper"
	"github.com/sniperkit/colly/pkg/queue"
	"github.com/sniperkit/colly/plugins/storage/external/sqlite3"
)

var version = "0.0.1-alpha"

func main() {
	// Instantiate collector
	c := colly.NewCollector(
		// Allow requests only to www.example.com
		colly.AllowedDomains("www.example.com"),
		//colly.Async(true),
		//colly.Debugger(&debug.LogDebugger{}),
	)

	storage := &sqlite3.Storage{
		Filename: "./results.db",
	}

	defer storage.Close()

	err := c.SetStorage(storage)
	if err != nil {
		panic(err)
	}

	helper.RandomUserAgent(c)
	helper.Referrer(c)

	q, _ := queue.New(8, storage)
	q.AddURL("http://www.example.com")

	// c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 4})

	// Find and visit all links
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		q.AddURL(e.Request.AbsoluteURL(e.Attr("href")))
	})

	c.OnResponse(func(r *colly.Response) {
		log.Println(r.Request.URL, "\t", r.StatusCode)
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Println(r.Request.URL, "\t", r.StatusCode, "\nError:", err)
	})

	q.Run(c)
	log.Println(c)
}
