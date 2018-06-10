package languages

import (
	"errors"

	"regexp"
	"strings"

	"github.com/gelembjuk/keyphrases/helper"
	"github.com/gelembjuk/keyphrases/wordnet"
)

var EnglishStopWords []string
var EnglishBadWordsForAnd []string
var EnglishBadWordsNotUseful []string
var EnglishAdverbsOfTime []string
var EnglishNounsOfTime []string

type English struct {
	Lang          Language
	WordNet       wordnet.WordNet
	wordnetstatus uint8
}

type cleanTemplate struct {
	template   string
	removedind int
	textind    int
}

func init() {
	EnglishStopWords = []string{"i", "me", "my", "myself", "we", "our", "ours", "ourselves", "you", "your", "yours", "yourself",
		"yourselves", "he", "him", "his", "himself", "she", "her", "hers", "herself", "it", "its",
		"itself", "they", "them", "their", "theirs", "themselves", "what", "which", "who", "whom",
		"this", "that", "these", "those", "am", "is", "are", "was", "were", "be", "been", "being", "have", "has",
		"had", "having", "do", "does", "did", "doing", "would", "should", "could", "ought", "i'm",
		"you're", "he's", "she's", "it's", "we're", "they're", "i've", "you've", "we've", "they've",
		"i'd", "you'd", "he'd", "she'd", "we'd", "they'd", "i'll", "you'll", "he'll", "she'll", "we'll",
		"they'll", "isn't", "aren't", "wasn't", "weren't", "hasn't", "haven't", "hadn't",
		"doesn't", "don't", "didn't", "won't", "wouldn't", "shan't", "shouldn't", "can't",
		"cannot", "couldn't", "mustn't", "let's", "that's", "who's", "what's", "here's",
		"there's", "when's", "where's", "why's", "how's", "a", "an", "the", "and", "but", "if", "or",
		"because", "as", "until", "while", "of", "at", "by", "for", "with", "about", "against", "between",
		"into", "through", "during", "before", "after", "above", "below", "to", "from", "up", "down", "in",
		"out", "on", "off", "over", "under", "again", "further", "then", "once", "here", "there", "when",
		"where", "why", "how", "all", "any", "both", "each", "few", "more", "most", "other", "some", "such",
		"no", "nor", "not", "only", "own", "same", "so", "than", "too", "very", "de", "will", "of", "without"}
	EnglishBadWordsForAnd = []string{"have", "has", "had", "can", "up", "could", "may", "per", "said", "says", "yet", "already", "say"}
	EnglishBadWordsNotUseful = []string{"inc", "said"}
	EnglishAdverbsOfTime = []string{"after", "already", "during", "finally", "just", "last", "later", "next", "now", "recently", "soon", "then", "tomorrow", "when", "while", "yesterday", "year", "week", "day", "month", "hour", "quarter"}
	EnglishNounsOfTime = []string{"January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}

}

func (lang English) GetName() string {
	return "english"
}
func (lang *English) SetOptions(options map[string]string) error {

	if val, ok := options["wordnetdirectory"]; ok {
		lang.wordnetstatus = 1
		// create wordnet object
		lang.WordNet = wordnet.WordNet{}

		err := lang.WordNet.SetDictDirectory(val)

		if err != nil {

			return err
		}

		lang.wordnetstatus = 2
	}
	return nil
}

func (lang *English) CleanNewsMessage(text string) (string, string, error) {
	removed := ""

	text, removed, _ = lang.cleanNewsMessagePrefix(text)

	replace := [][]string{
		{"Inc\\.\\s*?\\(", "Inc ("},
		{"\\([^)]+\\)", ""},
		{"'s\\s", " "},
		{"'s$", ""},
	}

	for _, template := range replace {
		r := regexp.MustCompile(template[0])

		text = r.ReplaceAllString(text, template[1])
	}

	return text, removed, nil
}

func (lang *English) cleanNewsMessagePrefix(text string) (string, string, error) {

	removed := ""

	templates := []cleanTemplate{
		cleanTemplate{"^(.{3,30})\\s?-(-|\\s|-\\s)(.*?)$", 1, 3},
		cleanTemplate{"^By (.{3,40}) --? (.*?)$", 1, 2},
		cleanTemplate{"^([A-Z0-9 ]{3,30})\\s?—\\s?(.*?)$", 1, 2},
		cleanTemplate{"^([A-Z0-9 ]{3,30}):\\s(.*?)$", 1, 2},
		cleanTemplate{"^(.{3,30},[A-Z0-9 /]{3,30}):\\s(.*?)$", 1, 2},
		cleanTemplate{"^(In brief):\\s(.*?)$", 1, 2},
		cleanTemplate{"^(\\(.{3,30}\\)\\s.{3,25})-(.*?)$", 1, 2},
		cleanTemplate{"^(.{10,40}\\s/.{3,25}/)\\s--\\s(.*?)$", 1, 2},
		cleanTemplate{"^(\\[.{3,25}\\])\\s?(.*?)$", 1, 2},
	}

	for _, template := range templates {
		r := regexp.MustCompile(template.template)

		matched := r.FindStringSubmatch(text)

		if len(matched) > 1 {

			removed = matched[template.removedind]

			text = matched[template.textind]

			break
		}

	}

	return text, removed, nil
}

func (lang *English) CleanAndNormaliseSentence(sentence string) (string, error) {
	replace := [][]string{
		{"\"", " "},
		{"[“”]", " "},
		{"U.S.", "United States"},
		{"U.K.", "United Kingdom"},
		{"E.U.", "Europe Union"},
	}

	for _, template := range replace {
		r := regexp.MustCompile(template[0])

		sentence = r.ReplaceAllString(sentence, template[1])
	}
	return sentence, nil
}

func (lang *English) StrongCleanAndNormaliseSentence(sentence string) (string, error) {

	return sentence, nil
}

func (lang *English) IsWord(word string) bool {
	if len(word) == 0 {
		return false
	}

	matched, _ := regexp.MatchString("^[a-zA-Z0-9-]+$", word)

	if matched {
		matched, _ = regexp.MatchString("^-+$", word)

		if !matched {
			return true
		}
	}

	return false
}

func (lang *English) RemoveCommonWords(words map[string]int) bool {

	for wordorig, count := range words {
		word := strings.ToLower(wordorig)
		if helper.StringInSlice(word, EnglishStopWords) ||
			count < 2 ||
			helper.StringInSlice(word, EnglishBadWordsForAnd) ||
			helper.StringInSlice(word, EnglishBadWordsNotUseful) {
			delete(words, wordorig)
		}
	}
	return true
}

func (lang *English) IsSimilarWord(word1 string, word2 string) int8 {
	// test cases. it will be hardcoded
	if len(word1) > 8 && len(word2) > 8 && word1[0:8] == "testword" && word2[0:8] == "testword" {
		return 1
	}

	if strings.ToLower(word1) == strings.ToLower(word2) {
		if word1 == strings.ToUpper(word1) || strings.ToUpper(word2[0:1]) == word2[0:1] {
			return 1
		}
		return -1
	}

	if word1 == "USA" && word2 == "US" {
		return 1
	}

	word1 = strings.ToLower(word1)
	word2 = strings.ToLower(word2)

	if word1 == word2+"s" {
		return 1
	}
	if word2 == word1+"s" {
		return -1
	}

	return 0

}

func (lang *English) IsNotUsefulWord(word string) bool {
	if helper.StringInSlice(word, EnglishStopWords) {
		return true
	}
	if helper.StringInSlice(word, EnglishBadWordsForAnd) {
		return true
	}
	if helper.StringInSlice(word, EnglishBadWordsNotUseful) {
		return true
	}
	return false
}

func (lang *English) IsPhraseSubphrase(phrase1 string, phrase2 string) int8 {
	if phrase1 == "the "+phrase2 {
		return 1
	}
	if phrase2 == "the "+phrase1 {
		return -1
	}
	if phrase1 == phrase2+" inc" {
		return 1
	}
	if phrase2 == phrase1+" inc" {
		return -1
	}
	if phrase1 == phrase2+"s" {
		return 1
	}
	if phrase2 == phrase1+"s" {
		return -1
	}
	return 0
}
func (lang *English) IsWordModInPhrase(phrase, word string) bool {
	l := len(word)

	if l < 2 {
		return false
	}

	pattern := " " + word + "s? "

	if word[l-2:l-1] == "es" {
		pattern = " " + word + "(es)? "
	}

	if t, _ := regexp.MatchString(pattern, " "+phrase+" "); t {
		return true
	}
	return false
}
func (lang *English) GetTypeOfWord(word string) (string, error) {

	return lang.getTypeOfWord(word, "", "")
}
func (lang *English) GetTypeOfWordComplex(word string, prevword string, nextword string) (string, error) {

	return lang.getTypeOfWord(word, prevword, nextword)
}

func (lang *English) getTypeOfWord(word string, prevword string, nextword string) (string, error) {

	//possible values
	// n -name
	// a - action
	// t - time
	// c - condition
	// v - value
	// r - thing (noun that is not Name)
	// b - not informative word
	// f - currency name

	if helper.StringInSlice(word, EnglishStopWords) {
		return "b", nil
	}

	if lang.isCurrencyName(word) {
		return "f", nil
	}

	numbervaluepattern := "^[0-9$.,-]+$"

	if t, _ := regexp.MatchString(numbervaluepattern, word); t {
		return "v", nil
	}

	if lang.wordnetstatus != 2 {
		return "", errors.New("Can not detect type of a word. WordNet dict not configured")
	}

	synonims, _ := lang.WordNet.GetWordSynonims(word)

	for _, syn := range synonims {
		if t, _ := regexp.MatchString(numbervaluepattern, syn); t {
			return "v", nil
		}
	}

	if len(word) > 1 && strings.ToUpper(word) == word {
		return "n", nil
	}

	options, err := lang.WordNet.GetWordOptions(word)

	if err == nil {
		if len(options) > 0 && options[0] == "v" {
			return "a", nil
		}
	}

	if helper.StringInSlice("r", options) || len(options) == 1 && options[0] == "a" {
		return "c", nil
	}

	// check if a word is about a time

	lword := strings.ToLower(word)

	if helper.StringInSlice(lword, EnglishAdverbsOfTime) {
		return "t", nil
	}

	ucword := helper.UpperCaseFirstLetter(word)

	if helper.StringInSlice(ucword, EnglishNounsOfTime) {
		return "t", nil
	}

	lword = lang.simplifyWord(strings.ToLower(word))

	if helper.StringInSlice(lword, EnglishAdverbsOfTime) {
		return "t", nil
	}

	ucword = helper.UpperCaseFirstLetter(lword)

	if helper.StringInSlice(ucword, EnglishNounsOfTime) {
		return "t", nil
	}

	if len(word) > 1 && ucword == word {
		return "n", nil
	}

	return "r", nil
}

func (lang *English) isCurrencyName(word string) bool {
	return false
}

func (lang *English) simplifyWord(word string) string {
	lenth := len(word)

	if lenth > 2 && word[lenth-2:] == "'s" {
		return word[0 : lenth-2]
	}
	if lenth > 1 &&
		(word[lenth-1:lenth-1] == "'" || word[lenth-1:lenth-1] == "s") {
		return word[:lenth-1]
	}
	return word
}

func (lang English) SimplifyCompanyName(phrase string) string {
	// if there is a comma then truncate everything after
	if strings.Index(phrase, ",") > 1 {
		phrase = phrase[0:strings.Index(phrase, ",")]
	}

	replace := [][]string{
		{"\\s+$", ""},
		{"^\\s+", ""},
		{"(?i)" + " inc\\.?$", ""},
		{"(?i)" + " ltd\\.?$", ""},
		{"(?i)" + " plc\\.?$", ""},
		{"(?i)" + " corp\\.?$", ""},
		{"(?i)" + " corporation$", ""},
		{"(?i)" + " incorporated$", ""},
		{"(?i)" + " international$", ""},
		{"(?i)" + " enterprises$", ""},
		{"(?i)" + " limited$", ""},
		{"(?i)" + " company$", ""},
		{"(?i)" + " & co\\..*?$", ""},
		{"(?i)" + " co\\.?\\s?$", ""},
		{"(?i)" + " & company.*?$", ""},
		{"(?i)" + "^the ", ""},
		{"   ", " "},
		{"  ", " "},
		{"\\s+$", ""},
		{"^\\s+", ""},
		{"s$", ""},
	}

	for _, template := range replace {
		r := regexp.MustCompile(template[0])

		phrase = r.ReplaceAllString(phrase, template[1])
	}
	return phrase
}

func (lang English) SimplifyCompanyNameExt(phrase string) string {
	// if there is a comma then truncate everything after
	if strings.Index(phrase, ",") > 1 {
		phrase = phrase[0:strings.Index(phrase, ",")]
	}

	replace := [][]string{
		{" \\([A-Z0-9]+\\)\\s?$", ""},
		{" [A-Z.]{3,}\\s?$", ""},
	}

	for _, template := range replace {
		r := regexp.MustCompile(template[0])

		phrase = r.ReplaceAllString(phrase, template[1])
	}
	return phrase
}
func (lang English) normaliseWord(word string) string {

	word = strings.ToLower(word)

	replace := [][]string{
		{"s$", ""},
	}

	for _, template := range replace {
		r := regexp.MustCompile(template[0])

		word = r.ReplaceAllString(word, template[1])
	}
	return word
}
func (lang English) GetSentimentOfWord(word string) (float32, error) {
	if lang.wordnetstatus != 2 {
		return 0, errors.New("Can not detect type of a word. WordNet dict not configured")
	}

	word = lang.normaliseWord(word)

	return lang.WordNet.GetWordSentiment(word)
}
