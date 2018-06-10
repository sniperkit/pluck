package sitemap

// Config
type Config struct {

	// Enabled
	Enabled bool `default:"false" flag:"with-sitemap-parser" yaml:"enabled" toml:"enabled" xml:"enabled" ini:"enabled" csv:"Enabled" json:"enabled" yaml:"enabled" toml:"enabled" xml:"enabled" ini:"enabled" csv:"Enabled"`

	// URL
	URL string `flag:"sitemap-url" json:"url" yaml:"url" toml:"url" xml:"url" ini:"URL" csv:"URL"`

	// AutoDetect
	AutoDetect bool `default:"false" flag:"sitemap-auto-detect" yaml:"auto_detect" toml:"auto_detect" xml:"autoDetect" ini:"autoDetect" csv:"AutoDetect" json:"auto_detect" yaml:"auto_detect" toml:"auto_detect" xml:"autoDetect" ini:"autoDetect" csv:"AutoDetect"`

	// LimitURLs
	LimitURLs int `default:"0" flag:"sitemap-limit" yaml:"limit_urls" toml:"limit_urls" xml:"limitURLs" ini:"limitURLs" csv:"limitURLs" json:"limit_urls" yaml:"limit_urls" toml:"limit_urls" xml:"limitURLs" ini:"limitURLs" csv:"limitURLs"`

	//-- END
}
