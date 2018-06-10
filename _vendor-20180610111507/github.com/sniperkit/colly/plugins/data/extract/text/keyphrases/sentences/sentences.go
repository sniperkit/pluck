package sentences

import (
	"regexp"
	"strings"

	"github.com/gelembjuk/keyphrases/helper"
	"github.com/gelembjuk/keyphrases/languages"
	"gopkg.in/neurosnap/sentences.v1"
	"gopkg.in/neurosnap/sentences.v1/data"
)

var langobj languages.LangClass

func SetLangObject(lang languages.LangClass) {
	langobj = lang
}

func SetLanguage(lang string) {
	langobj, _ = languages.GetLangObject(lang)
}

func SplitTextForSentencesFromNews(text string) ([]string, error) {
	// this text is from news sources. It can have specific news format
	// clean a text from standard news message formatting , and specific language

	return SplitText(text, true)
}

func SplitTextForSentences(text string) ([]string, error) {
	return SplitText(text, false)
}

func SplitText(text string, news bool) ([]string, error) {
	// prepare tokenizer
	sentenceslist := []string{}

	langfile := "data/" + langobj.GetName() + ".json"

	b, err := data.Asset(langfile)

	if err != nil {
		return sentenceslist, err
	}

	// load the training data
	training, err := sentences.LoadTraining(b)

	if err != nil {
		return sentenceslist, err
	}

	// create the default sentence tokenizer
	tokenizer := sentences.NewSentenceTokenizer(training)

	text, _ = helper.CleanTextAfterHTML(text)

	if news {
		// this text is from news sources. It can have specific news format
		// clean a text from standard news message formatting , and specific language
		text, _, _ = langobj.CleanNewsMessage(text)
	}

	sentencesobjs := tokenizer.Tokenize(text)

	for _, s := range sentencesobjs {
		sentence := s.Text

		// remove last symbol of a sentence if it is a dot or so
		if len(sentence) < 3 {
			continue
		}

		sentence, _ = cleanAndNormaliseSentence(sentence)

		sentenceslist = append(sentenceslist, sentence)
	}

	return sentenceslist, nil
}

func NormaliseSentencesList(sentenceslist []string) ([]string, error) {
	normalsentences := []string{}

	for _, sentence := range sentenceslist {
		sentence, _ = langobj.StrongCleanAndNormaliseSentence(sentence)
		sentence = strings.ToLower(sentence)

		replace := [][]string{
			{",", ""},
		}

		for _, template := range replace {
			r := regexp.MustCompile(template[0])

			sentence = r.ReplaceAllString(sentence, template[1])
		}

		normalsentences = append(normalsentences, sentence)
	}

	return normalsentences, nil
}

func cleanAndNormaliseSentence(sentence string) (string, error) {

	sentence, _ = langobj.CleanAndNormaliseSentence(sentence)

	replace := [][]string{
		{"\\(\\s*?https?://[^ )]+\\s*?\\)", ""},
		{"[\\[\\]}{]", ""},
		{"[:;-]", " "},
		{"[.?!):]\\s*?$", " "},
		{"\\s\\s+", " "},
		{"^\\s+", ""},
		{"\\s+$", ""},
	}

	for _, template := range replace {
		r := regexp.MustCompile(template[0])

		sentence = r.ReplaceAllString(sentence, template[1])
	}

	return sentence, nil
}
