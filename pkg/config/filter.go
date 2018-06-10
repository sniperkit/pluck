package config

// FilterMode specifies the filtering list behaviour type.
type FilterMode string

// Enum list of filter modes.
const (
	WHITELIST FilterMode = "whitelist"
	BLACKLIST FilterMode = "blacklist"
	STOPLIST  FilterMode = "stoplist"
	SKIPLIST  FilterMode = "skiplist"
	WATCHLIST FilterMode = "watchlist"
)

// Filters specifies a whitelist, backlist, skiplist and stoplist of patterns
type Filters struct {

	// Enabled
	Enabled bool `default:"true" json:"enabled" yaml:"enabled" toml:"enabled" xml:"enabled" ini:"enabled"`

	// Stats enables to save more stats about the matched filters
	Stats bool `default:"false" json:"stats" yaml:"stats" toml:"stats" xml:"stats" ini:"stats"`

	// Mode
	Mode FilterMode `required:"true" default:"blacklist" json:"mode" yaml:"mode" toml:"mode" xml:"mode" ini:"mode"`

	// ReplaceBy
	ReplaceBy string `default:"*" json:"replace_by" yaml:"replace_by" toml:"replace_by" xml:"replaceBy" ini:"replaceBy"`

	// Patterns
	Patterns []*Filter `json:"patterns" yaml:"patterns" toml:"patterns" xml:"patterns" ini:"patterns"`

	//-- End
}

// Filter specifies a behaviour if the pattern is matched
type Filter struct {

	// Enabled
	Enabled bool `default:"true" json:"enabled,omitempty" yaml:"enabled,omitempty" toml:"enabled,omitempty" xml:"enabled,omitempty" ini:"enabled,omitempty"`

	// Match
	Match string `required:"true" json:"match" yaml:"match" toml:"match" xml:"match" ini:"match"`

	// Mode
	Mode MatchMode `required:"true" default:"mode" json:"mode" yaml:"mode" toml:"mode" xml:"mode" ini:"mode"`

	// ReplaceBy
	ReplaceBy string `default:"*" json:"replace_by" yaml:"replace_by" toml:"replace_by" xml:"replaceBy" ini:"replaceBy"`

	//-- End
}
