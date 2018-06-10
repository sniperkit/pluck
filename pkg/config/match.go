package config

/*
	Boolean queries allow the following special operators to be used:

	- Operator MAYBE:
	   - `hello MAYBE world`
	- Operator OR:
	   - `hello | world`
	- Operator NOT:
	   - `hello -world`
	   - `hello !world`
	- Grouping:
	   - `( hello world )`
*/

// MatchMode specifies the matching mode to use to scan the plucked occurence
type MatchMode string

// Enum list of all available matching modes
const (
	MATCH_ALL     MatchMode = "all"     // matches plucked occurence with all the pattern words
	MATCH_ANY     MatchMode = "any"     // matches plucked occurence with any of the pattern words
	MATCH_PHRASE  MatchMode = "phrase"  // matches plucked occurence as a phrase, requiring perfect match of the phrase attribute
	MATCH_BOOLEAN MatchMode = "boolean" // matches plucked occurence as a boolean expression
)

// Match
type Match struct {

	// Patterns matching modes inside plucked occurences
	Mode string `default:"any" json:"mode" yaml:"mode" toml:"mode" xml:"mode" ini:"mode"`

	// separator inside the match to use if we want to join all the occurences into a slice of strings
	Separator string `json:"separator" yaml:"separator" toml:"separator" xml:"separator" ini:"separator"`

	// Split plucked occurences with a user-defined separator
	Split bool `default:"true" json:"split" yaml:"split" toml:"split" xml:"split" ini:"split"`

	// separator inside the match to use if we want to join all the occurences into a slice of strings
	Phrase string `json:"phrase" yaml:"phrase" toml:"phrase" xml:"phrase" ini:"phrase"`

	//-- End
}
