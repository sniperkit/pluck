package config

// Configs specifies...
type Configs struct {

	// Debug
	Debug bool `default:"false" json:"debug,omitempty" yaml:"debug,omitempty" toml:"debug,omitempty" xml:"debug,omitempty" ini:"debug,omitempty"`

	// Verbose
	Verbose bool `default:"true" json:"verbose,omitempty" yaml:"verbose,omitempty" toml:"verbose,omitempty" xml:"verbose,omitempty" ini:"verbose,omitempty"`

	// XDGBaseDir specifies
	XDGBaseDir string `json:"xdg_base_dir,omitempty" yaml:"xdg_base_dir,omitempty" toml:"xdg_base_dir,omitempty" xml:"xdgBaseDir,omitempty" ini:"xdgBaseDir,omitempty"`

	// Pluck specifies the list of content plucking units
	Pluck []Config `json:"plucker" yaml:"plucker" toml:"plucker" xml:"plucker" ini:"plucker"`
}

// Config specifies parameters for plucking
type Config struct {

	// Enabled
	Enabled bool `default:"true" json:"enabled,omitempty" yaml:"enabled,omitempty" toml:"enabled,omitempty" xml:"enabled,omitempty" ini:"enabled,omitempty"`

	// Debug
	Debug bool `default:"false" json:"debug,omitempty" yaml:"debug,omitempty" toml:"debug,omitempty" xml:"debug,omitempty" ini:"debug,omitempty"`

	// Verbose
	Verbose bool `default:"true" json:"verbose,omitempty" yaml:"verbose,omitempty" toml:"verbose,omitempty" xml:"verbose,omitempty" ini:"verbose,omitempty"`

	// Sanitize html content
	Sanitize bool `default:"false" json:"sanitize,omitempty" yaml:"sanitize,omitempty" toml:"sanitize,omitempty" xml:"sanitize,omitempty" ini:"sanitize,omitempty"`

	// the key in the returned map, after completion
	Name string `required:"true" json:"name" yaml:"name" toml:"name" xml:"name" ini:"name"`

	// must be found in order, before capturing commences
	Activators []string `json:"activators" yaml:"activators" toml:"activators" xml:"activators" ini:"activators"`

	// number of activators that stay permanently (counted from left to right)
	Permanent int `json:"permanent" yaml:"permanent" toml:"permanent" xml:"permanent" ini:"permanent"`

	// restarts capturing
	Deactivator string `required:"true" json:"deactivator" yaml:"deactivator" toml:"deactivator" xml:"deactivator" ini:"deactivator"`

	// finishes capturing this pluck
	Finisher string `json:"finisher,omitempty" yaml:"finisher,omitempty" toml:"finisher,omitempty" xml:"finisher,omitempty" ini:"finisher,omitempty"`

	// specifies the number of times capturing can occur
	Limit int `default:"-1" json:"limit" yaml:"limit" toml:"limit" xml:"limit" ini:"limit"`

	// maximum number of characters for a capture
	Maximum int `json:"maximum,omitempty" yaml:"maximum,omitempty" toml:"maximum,omitempty" xml:"maximum,omitempty" ini:"maximum,omitempty"`

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
