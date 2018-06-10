package config

// Debug
type Debug struct {

	// Debug
	Debug bool `default:"false" json:"debug" yaml:"debug" toml:"debug" xml:"debug" ini:"debug"`

	// Verbose
	Verbose bool `default:"true" json:"verbose" yaml:"verbose" toml:"verbose" xml:"verbose" ini:"verbose"`

	//-- End
}
