package filter

// Domain
type Domain struct {

	// Enabled
	Enabled bool `default:"false" json:"enabled" yaml:"enabled" toml:"enabled" xml:"enabled" ini:"enabled"`

	// Address
	Address string `json:"address" yaml:"address" toml:"address" xml:"address" ini:"address"`

	// Protocol
	Protocol string `json:"protocol" yaml:"protocol" toml:"protocol" xml:"protocol" ini:"protocol"`

	// Host
	Host string `json:"host" yaml:"host" toml:"host" xml:"host" ini:"host"`

	// Port
	Port string `json:"port" yaml:"port" toml:"port" xml:"port" ini:"port"`

	// ForceSSL
	ForceSSL bool `default:"true" json:"force_ssl" yaml:"force_ssl" toml:"force_ssl" xml:"forceSSL" ini:"forceSSL"`

	// VerifySSL
	VerifySSL bool `default:"false" json:"ssl_verify" yaml:"ssl_verify" toml:"ssl_verify" xml:"verifySSL" ini:"verifySSL"`

	//-- END
}
