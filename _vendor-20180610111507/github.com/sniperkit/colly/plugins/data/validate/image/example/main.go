package main

import (
	"log"
	"os"

	"github.com/sniperkit/colly/plugins/data/validate/image"
)

func main() {
	trusted, err := secureimage.Check(os.Args[1])

	if err != nil {
		panic(err)
	}

	if trusted {
		log.Println("yes.")
	} else {
		log.Println("bad image file")
	}
}
