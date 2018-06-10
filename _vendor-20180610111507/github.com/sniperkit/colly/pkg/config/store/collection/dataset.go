package collection

import (
	export "github.com/sniperkit/colly/pkg/config/export"
)

// Dataset
type Dataset struct {

	// Enabled
	Enabled bool `default:"false" json:"enabled" yaml:"enabled" toml:"enabled" xml:"enabled" ini:"enabled"`

	// IsExportable
	IsExportable bool `default:"true" json:"is_exportable" yaml:"is_exportable" toml:"is_exportable" xml:"is_exportable" ini:"is_exportable"`

	// Datasets
	Datasets []Dataset `json:"datasets" yaml:"datasets" toml:"datasets" xml:"datasets" ini:"datasets"`

	// MaxRows
	MaxRows int `default:"100000" json:"max_rows" yaml:"max_rows" toml:"max_rows" xml:"max_rows" ini:"max_rows"`

	// MaxCols
	MaxCols int `default:"100" json:"max_cols" yaml:"max_cols" toml:"max_cols" xml:"max_cols" ini:"max_cols"`

	// Charset
	Charset string `default:"UTF-8" json:"charset" yaml:"charset" toml:"charset" xml:"charset" ini:"charset"`

	// Exports
	Exports []export.Config `json:"exports" yaml:"exports" toml:"exports" xml:"exports" ini:"exports"`

	// - End
}
