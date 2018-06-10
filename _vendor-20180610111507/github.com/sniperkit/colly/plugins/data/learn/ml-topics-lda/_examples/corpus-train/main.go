package main

import (
	"log"

	topics "github.com/sniperkit/colly/plugins/mlearn/lda-topics"
)

func main() {
	processor := topics.NewProcessor(
		topics.Transformations{
			topics.ToLower,
			topics.RemoveTwitterUsernames,
			topics.Sanitize,
			topics.MinLen,
			topics.GetStopwordFilter("../stopwords/en"),
			topics.GetStopwordFilter("../stopwords/se"),
		},
	)
	corpus, err := processor.ImportSingleFileCorpus(topics.NewCorpus(), "./corpus")
	if err != nil {
		log.Fatalln("error while importing the corpus from a local file: ", err)
	}

	lda := topics.NewLDA(
		&topics.Configuration{
			Verbose:       true,
			PrintInterval: 500,
			PrintNumWords: 8},
	)
	err = lda.Init(corpus, 8, 0, 0)
	if err != nil {
		log.Fatalln("error while creating the lda classifier: ", err)
	}

	_, err = lda.Train(10000)
	if err != nil {
		log.Fatalln("error while training the lda classifier: ", err)
	}

	lda.PrintTopWords(8)
}
