package main

import (
	imagga "github.com/sniperkit/colly/plugins/data/aggregate/web/plugin/indexer/photo-imagga"
	"github.com/sniperkit/colly/plugins/data/aggregate/web/service/plugin"
)

func main() {
	server, err := plugin.Server()
	if err != nil {
		panic(err)
	}
	server.RegisterPipeline(new(imagga.Indexer))
	server.Serve()
}
