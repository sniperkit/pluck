package main

import (
	"fmt"
	"time"

	colly "github.com/sniperkit/colly/pkg"
)

func indexer() {
	fmt.Println("Colly based JSON-LD indexer")

	url := "https://httpbin.org/delay/2"

	// urls := []

	c := colly.NewCollector()

	c.Limit(&colly.LimitRule{
		// DomainGlob:  "*httpbin.*",
		Parallelism: 1,
		//Delay:      5 * time.Second,
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Starting", r.URL, time.Now())
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL, time.Now())
	})

	for i := 0; i < 120; i++ {
		go c.Visit(fmt.Sprintf("%s?n=%d", url, i))
	}

	// use this to get the body and route out to other functions to work on
	// c.OnResponse(func(r *colly.Response) {
	// 	fmt.Println(r.Ctx.Get("url"))
	// })

	c.Visit(url)
	c.Wait()

}
