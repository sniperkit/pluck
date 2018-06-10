package collection

import (
	export "github.com/sniperkit/colly/pkg/config/export"
)

// Databook
type Databook struct {

	// Enabled
	Enabled bool `default:"false" json:"enabled" yaml:"enabled" toml:"enabled" xml:"enabled" ini:"enabled"`

	// IsExportable
	IsExportable bool `default:"true" json:"is_exportable" yaml:"is_exportable" toml:"is_exportable" xml:"is_exportable" ini:"is_exportable"`

	// Datasets
	Datasets []string `json:"datasets" yaml:"datasets" toml:"datasets" xml:"datasets" ini:"datasets"`

	// MaxDatasets
	MaxDatasets int `default:"5" json:"max_datasets" yaml:"max_datasets" toml:"max_datasets" xml:"maxDatasets" ini:"maxDatasets"`

	// Charset
	Charset string `default:"UTF-8" json:"charset" yaml:"charset" toml:"charset" xml:"charset" ini:"charset"`

	// Exports
	Exports []export.Config `json:"exports" yaml:"exports" toml:"exports" xml:"exports" ini:"exports"`

	// - End
}
