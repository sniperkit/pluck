package cache

import (
	"time"

	storage "github.com/sniperkit/colly/pkg/config/store"
)

// Config
type Config struct {

	// Enabled
	Enabled bool `default:"false" flag:"with-http-cache" yaml:"enabled" toml:"enabled" xml:"enabled" ini:"enabled" csv:"Enabled"`

	// Backend
	Backend string `default:"badger" flag:"http-cache-backend" yaml:"backend" toml:"backend" xml:"backend" ini:"backend" csv:"backend"`

	// TTL
	TTL time.Duration `default:"3600s" flag:"http-cache-ttl" yaml:"ttl" toml:"ttl" xml:"ttl" ini:"ttl" csv:"TTL"`

	// Store
	Store storage.Config `json:"store" yaml:"store" toml:"store" xml:"store" ini:"store" csv:"store"`

	//-- END
}
