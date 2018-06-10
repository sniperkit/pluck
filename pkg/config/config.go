package config

// Config specifies parameters for plucking
type Config struct {

	// Enabled
	Enabled bool `default:"true" json:"enabled" yaml:"enabled" toml:"enabled" xml:"enabled" ini:"enabled"`

	// Debug
	Debug bool `default:"false" json:"debug" yaml:"debug" toml:"debug" xml:"debug" ini:"debug"`

	// Verbose
	Verbose bool `default:"true" json:"verbose" yaml:"verbose" toml:"verbose" xml:"verbose" ini:"verbose"`

	// must be found in order, before capturing commences
	Activators []string `json:"activators" yaml:"activators" toml:"activators" xml:"activators" ini:"activators"`

	// number of activators that stay permanently (counted from left to right)
	Permanent int `json:"permanent" yaml:"permanent" toml:"permanent" xml:"permanent" ini:"permanent"`

	// restarts capturing
	Deactivator string `required:"true" json:"deactivator" yaml:"deactivator" toml:"deactivator" xml:"deactivator" ini:"deactivator"`

	// finishes capturing this pluck
	Finisher string `json:"finisher" yaml:"finisher" toml:"finisher" xml:"finisher" ini:"finisher"`

	// specifies the number of times capturing can occur
	Limit int `json:"limit" yaml:"limit" toml:"limit" xml:"limit" ini:"limit"`

	// the key in the returned map, after completion
	Name string `json:"name" yaml:"name" toml:"name" xml:"name" ini:"name"`

	// Sanitize html content
	Sanitize bool `json:"sanitize" yaml:"sanitize" toml:"sanitize" xml:"sanitize" ini:"sanitize"`

	// maximum number of characters for a capture
	Maximum int `json:"maximum" yaml:"maximum" toml:"maximum" xml:"maximum" ini:"maximum"`

	// Match specifies...
	Match Match `json:"match" yaml:"match" toml:"match" xml:"match" ini:"match"`

	/*
		// separator inside the match to use if we want to join all the occurences into a slice of strings
		Separator string `json:"separator" yaml:"separator" toml:"separator" xml:"separator" ini:"separator"`

		// Split plucked occurences with a user-defined separator
		Split bool `default:"true" json:"split" yaml:"split" toml:"split" xml:"split" ini:"split"`

		// Patterns matching modes inside plucked occurences
		Mode string `default:"any" json:"mode" yaml:"mode" toml:"mode" xml:"mode" ini:"mode"`

		// separator inside the match to use if we want to join all the occurences into a slice of strings
		Phrase string `json:"phrase" yaml:"phrase" toml:"phrase" xml:"phrase" ini:"phrase"`
	*/

	// identify some optional patterns to split down results
	Patterns []string `json:"patterns" yaml:"patterns" toml:"patterns" xml:"patterns" ini:"patterns"`

	// set a word list to include a plucked occurrence
	Whitelist []string `json:"whitelist" yaml:"whitelist" toml:"whitelist" xml:"whitelist" ini:"whitelist"`

	// set a word list to exclude a plucked occurrence
	Blacklist []string `json:"activatoblacklistrs" yaml:"blacklist" toml:"blacklist" xml:"blacklist" ini:"blacklist"`

	//-- End
}
