package pluck

import (
	plucker "github.com/sniperkit/colly/plugins/data/extract/text-plucker/config"
)

// Config
type Config struct {

	// Enabled
	Enabled bool `default:"false" flag:"enable-plucker" yaml:"enabled" toml:"enabled" xml:"enabled" ini:"enabled" csv:"Enabled"`

	// Rules
	Units []plucker.Config `yaml:"rules" toml:"rules" xml:"rules" ini:"rules" csv:"rules"`

	//-- END
}
