package config

// Config specifies parameters for plucking content
type Config struct {

	// Name sets the key in the returned map, after completion
	Name string `json:"name" yaml:"name" toml:"name" xml:"name" ini:"name"`

	// Activators must be found in order, before capturing commences
	Activators []string `json:"activators" yaml:"activators" toml:"activators" xml:"activators" ini:"activators"`

	// Deactivator restarts capturing
	Deactivator string `json:"deactivator" yaml:"deactivator" toml:"deactivator" xml:"deactivator" ini:"deactivator"`

	// Permanent set the number of activators that stay permanently (counted from left to right)
	Permanent int `json:"permanent" yaml:"permanent" toml:"permanent" xml:"permanent" ini:"permanent"`

	// Finisher trigger the end of capturing this pluck
	Finisher string `json:"finisher" yaml:"finisher" toml:"finisher" xml:"finisher" ini:"finisher"`

	// Limit specifies the number of times capturing can occur
	Limit int `default:"-1" json:"limit" yaml:"limit" toml:"limit" xml:"limit" ini:"limit"`

	// Sanitize enables the html stripping
	Sanitize bool `default:"false" json:"sanitize" yaml:"sanitize" toml:"sanitize" xml:"sanitize" ini:"sanitize"`

	// Maximum set the number of characters for a capture
	Maximum int `json:"maximum" yaml:"maximum" toml:"maximum" xml:"maximum" ini:"maximum"`

	//-- End
}
