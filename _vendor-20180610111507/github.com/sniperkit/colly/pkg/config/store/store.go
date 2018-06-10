package store

// Config
type Config struct {

	// Enabled
	Enabled bool `default:"false" yaml:"enabled" toml:"enabled" xml:"enabled" ini:"enabled" csv:"enabled" json:"enabled" yaml:"enabled" toml:"enabled" xml:"enabled" ini:"enabled" csv:"enabled"`

	// Domain
	Domain string `json:"domain" yaml:"domain" toml:"domain" xml:"domain" ini:"domain"`

	// Protocol
	Protocol string `json:"protocol" yaml:"protocol" toml:"protocol" xml:"protocol" ini:"protocol"`

	// Host
	Host string `json:"host" yaml:"host" toml:"host" xml:"host" ini:"host"`

	// Port
	Port string `json:"port" yaml:"port" toml:"port" xml:"port" ini:"port"`

	// Directory
	Directory string `json:"prefix_path" yaml:"prefix_path" toml:"prefix_path" xml:"prefix_path" ini:"prefix_path"`

	// ForceSSL
	ForceSSL bool `default:"true" json:"force_ssl" yaml:"force_ssl" toml:"force_ssl" xml:"force_ssl" ini:"force_ssl"`

	// VerifySSL
	VerifySSL bool `default:"false" json:"ssl_verify" yaml:"ssl_verify" toml:"ssl_verify" xml:"verifySSL" ini:"verifySSL"`

	//-- End
}
