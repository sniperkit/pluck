package filter

// Config stores the filtering strategy for the colly collector.
type Config struct {
	Enabled     bool   `default:"true" json:"enabled" yaml:"enabled" toml:"enabled" xml:"enabled" ini:"enabled"`
	Rule        string `default:"" json:"rule" yaml:"rule" toml:"rule" xml:"rule" ini:"rule"`
	ScannerType string `default:"regex" json:"scanner" yaml:"scanner" toml:"scanner" xml:"scanner" ini:"scanner"`
	isValid     bool   `json:"-" yaml:"-" toml:"-" xml:"-" ini:"-"`
}

// String
func (f *Config) String() string {
	return ""
}

// IsValid
func (f *Config) IsValid(pattern string) bool {
	return f.isValid
}

// AddRule
func (f *Config) AddRule(pattern string) bool {
	return false
}

// checkRuleByName
func checkRuleByName(pattern string) bool {
	return false
}
