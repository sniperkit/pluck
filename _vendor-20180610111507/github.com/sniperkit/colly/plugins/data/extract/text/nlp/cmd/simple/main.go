package main

import (
	"fmt"

	"github.com/sniperkit/colly/plugins/data/extract/text/nlp"
)

type Song struct {
	Name   string
	Artist string
}

func main() {
	songSamples := []string{
		"play {Name} by {Artist}",
		"play {Name} from {Artist}",
		"play {Name}",
		"from {Artist} play {Name}",
		"I want to hear {Name} by {Artist}",
		"I want to hear {Name}",
	}

	nl := nlp.New()
	err := nl.RegisterModel(Song{}, songSamples)
	if err != nil {
		panic(err)
	}

	err = nl.Learn() // you must call Learn after all models are registered and before calling P
	if err != nil {
		panic(err)
	}

	// after learning you can call P the times you want
	s := nl.P("hello sir can you pleeeeeease play King by Lauren Aquilina")
	if song, ok := s.(Song); ok {
		fmt.Println("Success")
		fmt.Printf("%#v\n", song)
	} else {
		fmt.Println("Failed")
	}
}

// Prints
//
// Success
// main.Song{Name: "King", Artist: "Lauren Aquilina"}
