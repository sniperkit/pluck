package main

import (
	"fmt"
	"time"

	"github.com/sniperkit/colly/plugins/data/extract/text/nlp"
)

func main() {

	type Song struct {
		Name       string
		Artist     string
		Message    string
		ReleasedAt time.Time
	}

	songSamples := []string{
		"hello {Message} ",
		"{Name} {Artist}",
		"play {Name} by {Artist}",
		"play {Name} from {Artist}",
		"play {Name}",
		"from {Artist} play {Name}",
		"play something from {ReleasedAt}",
	}

	nl := nlp.New()
	err := nl.RegisterModel(Song{}, songSamples, nlp.WithTimeFormat("2006"))
	if err != nil {
		panic(err)
	}

	err = nl.Learn() // you must call Learn after all models are registered and before calling P
	if err != nil {
		panic(err)
	}

	// after learning you can call P the times you want
	s := nl.P("hello sir can you pleeeeeease King play by Lauren Aquilina")
	if song, ok := s.(*Song); ok {
		fmt.Println("Success")
		fmt.Printf("%#v\n", song)
	} else {
		fmt.Println("Failed")
	}

	// Prints
	//
	// Success
	// &main.Song{Name: "King", Artist: "Lauren Aquilina"}
}
