package main

import (
	"fmt"
	"plugin"
)

func main() {
	p, err := plugin.Open("../../shared/libs/bitcq/bitcq.so")
	if err != nil {
		panic(err)
	}

	f, err := p.Lookup("Search")
	if err != nil {
		panic(err)
	}

	search := f.(func(string, string) []byte)
	results := search("Unfriended", "qwerty")
	fmt.Println(string(results))
}
