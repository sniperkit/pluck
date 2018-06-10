package main

import (
	"github.com/sniperkit/colly/plugins/data/aggregate/web/plugin/scenario/ann"
	"github.com/sniperkit/colly/plugins/data/aggregate/web/service/plugin"
)

func main() {
	server, err := plugin.Server()
	if err != nil {
		panic(err)
	}
	server.RegisterPipeline(new(ann.Scenario))
	server.Serve()
}
