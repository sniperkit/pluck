package main

import (
	// "os"

	pluck "github.com/sniperkit/colly/plugins/data/extract/text/pluck"
)

func pluck_selector_query_from_config_file(filePath string) {
	p, err := pluck.New()
	panic(err)
	p.Verbose(true)
	p.Load(filePath)
	p.ResultJSON(true)
	p.PluckFile("output/test_plucker.json")
}

func pluck_with_config(queries ...string) {
	p, err := pluck.New()
	panic(err)
	p.Verbose(true)
	for _, pc := range PLUCKER_CONFIG_UNITS {
		p.Add(*pc)
	}
	for _, query := range queries {
		p.PluckString(query)
	}
	p.ResultJSON(true)
	p.PluckFile("output/test_plucker.json")
}
