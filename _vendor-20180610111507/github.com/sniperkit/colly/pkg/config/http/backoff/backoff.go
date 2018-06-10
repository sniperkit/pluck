package backoff

// Config
type Config struct {

	// Enabled
	Enabled bool `default:"false" flag:"with-backoff" yaml:"enabled" toml:"enabled" xml:"enabled" ini:"enabled" csv:"Enabled"`

	// Retry
	Retry int `default:"3" flag:"with-backoff-retry" yaml:"retry" toml:"retry" xml:"retry" ini:"retry" csv:"retry"`

	//-- End
}
