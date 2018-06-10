package filter

// Type represents a Filter engine
type Type string

const (
	REGEXP Type = "regexp" // Using default golang regexp package (default)
	PLUCK  Type = "pluck"  // Alternative to xpath and regexp, it allows to extract pattern with activators/desactivators patterns
	LEXER  Type = "lexer"  // to do...
	RUNE   Type = "rune"   // to do...
	XQUERY Type = "xquery" // to do...
	XPATH  Type = "xpath"  // to do...
	AST    Type = "ast"    // to do...
)

// Filters
type Filter struct {

	// Response
	Response Response `json:"response" yaml:"response" toml:"response" xml:"response" ini:"response" csv:"Response"`

	// Blacklists
	Blacklists Blacklist `json:"blacklists" yaml:"blacklists" toml:"blacklists" xml:"blackLists" ini:"blackLists" csv:"BlackLists"`

	// Whitelists
	Whitelists Whitelist `json:"whitelists" yaml:"whitelists" toml:"whitelists" xml:"whiteLists" ini:"whiteLists" csv:"Whitelists"`

	// Skiplists
	Skiplists Skiplist `json:"skiplists" yaml:"skiplists" toml:"skiplists" xml:"skiplists" ini:"skiplists" csv:"Skiplists"`

	//-- END
}
