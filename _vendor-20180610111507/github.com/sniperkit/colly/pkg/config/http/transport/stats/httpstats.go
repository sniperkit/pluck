package stats

import (
	http_client "github.com/sniperkit/colly/pkg/config/http/client"
	storage "github.com/sniperkit/colly/pkg/config/store"
)

// Config
type Config struct {

	// Enabled
	Enabled bool `default:"false" flag:"with-http-stats" yaml:"enabled" toml:"enabled" xml:"enabled" ini:"enabled" csv:"Enabled"`

	// Client
	Client http_client.Config `default:"" yaml:"client" toml:"client" xml:"client" ini:"client" csv:"client"`

	// Store
	Store storage.Config `json:"store" yaml:"store" toml:"store" xml:"store" ini:"store" csv:"store"`

	//-- END
}
