package tachymeter

import (
	export "github.com/sniperkit/colly/pkg/config/export"
)

// Config
type Config struct {

	// Enabled
	Enabled bool `default:"false" flag:"with-tachymeter" yaml:"enabled" toml:"enabled" xml:"enabled" ini:"enabled" csv:"Enabled" json:"enabled" yaml:"enabled" toml:"enabled" xml:"enabled" ini:"enabled" csv:"Enabled"`

	// Async
	Async bool `default:"false" flag:"with-tachymeter-async" yaml:"async" toml:"async" xml:"async" ini:"async" csv:"Async" json:"async" yaml:"async" toml:"async" xml:"async" ini:"async" csv:"Async"`

	// SampleSize
	SampleSize int `default:"50" flag:"with-tachymter-sample-size" yaml:"sample_size" toml:"sample_size" xml:"sampleSize" ini:"sampleSize" csv:"SampleSize" json:"sample_size" yaml:"sample_size" toml:"sample_size" xml:"sampleSize" ini:"sampleSize" csv:"SampleSize"`

	// HistogramBins
	HistogramBins int `default:"10" flag:"with-tachymter-histogram-bins" yaml:"histogram_bins" toml:"histogram_bins" xml:"histogramBins" ini:"histogramBins" csv:"HistogramBins" json:"histogram_bins" yaml:"histogram_bins" toml:"histogram_bins" xml:"histogramBins" ini:"histogramBins" csv:"HistogramBins"`

	// Export
	Export export.Config `json:"export" yaml:"export" toml:"export" xml:"export" ini:"export" csv:"Export"`

	//-- END
}
