package transport

import (
	// http_client "github.com/sniperkit/colly/pkg/config/http/client"
	http_cache "github.com/sniperkit/colly/pkg/config/http/transport/cache"
	http_stats "github.com/sniperkit/colly/pkg/config/http/transport/stats"
)

// Config defines components to load to provide custom http transport decorator.
type Config struct {

	// Enabled
	Enabled bool `default:"false" flag:"with--custom-http-transport" yaml:"enabled" toml:"enabled" xml:"enabled" ini:"enabled" csv:"Enabled"`

	// Cache stores httpcache configuration of the cache store and the http transport decorator
	Cache http_cache.Config `json:"cache" yaml:"cache" toml:"cache" xml:"cache" ini:"cache" csv:"cache"`

	// Stats stores httpstats configuration of the http transport decorator
	Stats http_stats.Config `json:"stats" yaml:"stats" toml:"stats" xml:"stats" ini:"stats" csv:"stats"`
}
