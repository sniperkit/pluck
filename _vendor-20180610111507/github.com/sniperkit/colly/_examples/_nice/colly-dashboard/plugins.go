package main

import (
	"plugin"
)

/*
Refs:
	- https://github.com/altsab/goplug/blob/master/main.go
*/

var (
	defaultCollectorPluginFilepath string = "../../../../shared/libs/bitcq/bitcq.so"
)

func loadPlugin(filePath string) {
	p, err := plugin.Open(filePath)
	if err != nil {
		panic(err)
	}

	f, err := p.Lookup("Search")
	if err != nil {
		panic(err)
	}

	search := f.(func(string, string) []byte)
	results := search("Unfriended", "qwerty")
	log.Println(string(results))
}
