package collector

import (
	storage_config "github.com/sniperkit/colly/pkg/config/store"
	// storage_backend "github.com/sniperkit/colly/pkg/config/store/backend"
)

// Config
type Config struct {

	// Enabled
	Enabled bool `default:"false" flag:"with-cache" yaml:"enabled" toml:"enabled" xml:"enabled" ini:"enabled" csv:"Enabled" json:"enabled" yaml:"enabled" toml:"enabled" xml:"enabled" ini:"enabled" csv:"Enabled"`

	// Backend
	Backend string `default:"inMemory" flag:"with-cache-backend" yaml:"backend" toml:"backend" xml:"backend" ini:"backend" csv:"backend" json:"backend" yaml:"backend" toml:"backend" xml:"backend" ini:"backend" csv:"backend"`

	// Store
	Store storage_config.Config `json:"store" yaml:"store" toml:"store" xml:"store" ini:"store" csv:"store"`

	//-- END
}
