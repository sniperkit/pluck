package collection

// Config
type Config struct {

	// Enabled
	Enabled bool `default:"false" flag:"with-collections" yaml:"enabled" toml:"enabled" xml:"enabled" ini:"enabled" csv:"Enabled"`

	// Databooks
	Databooks []Databook `json:"databooks" yaml:"databooks" toml:"databooks" xml:"databooks" ini:"databooks" csv:"databooks"`

	// Datasets
	Datasets []Dataset `json:"datasets" yaml:"datasets" toml:"datasets" xml:"datasets" ini:"datasets" csv:"datasets"`

	//-- END
}
