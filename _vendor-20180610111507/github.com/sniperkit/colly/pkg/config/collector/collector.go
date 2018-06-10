package collector

import (
	// collectors conponents config
	collector_cache "github.com/sniperkit/colly/pkg/config/collector/cache"
	collector_mode "github.com/sniperkit/colly/pkg/config/collector/mode"

	// content filtering configs
	filter_content "github.com/sniperkit/colly/pkg/config/filter"

	// http decorators config
	http_proxy "github.com/sniperkit/colly/pkg/config/http/proxy"
	http_transport "github.com/sniperkit/colly/pkg/config/http/transport"

	// preload entries to visit
	loader_sitemap "github.com/sniperkit/colly/pkg/config/loader/sitemap"
)

type Config struct {

	// RootURL
	RootURL string `required:"true" flag:"start-url" yaml:"root_url" toml:"root_url" xml:"rootURL" ini:"rootURL" csv:"RootURL" json:"root_url" yaml:"root_url" toml:"root_url" xml:"rootURL" ini:"rootURL" csv:"RootURL"`

	// UserAgent is the User-Agent string used by HTTP requests
	UserAgent string `default:"colly - https://github.com/sniperkit/colly" flag:"user-agent" yaml:"user_agent" toml:"user_agent" xml:"userAgent" ini:"userAgent" csv:"userAgent" json:"user_agent" yaml:"user_agent" toml:"user_agent" xml:"userAgent" ini:"userAgent" csv:"userAgent"`

	// RandomUserAgent specifies to generate a random User-Agent string for all HTTP requests
	RandomUserAgent bool `default:"false" flag:"with-random-user-agent" yaml:"random_user_agent" toml:"random_user_agent" xml:"randomUserAgent" ini:"randomUserAgent" csv:"randomUserAgent" json:"random_user_agent" yaml:"random_user_agent" toml:"random_user_agent" xml:"randomUserAgent" ini:"randomUserAgent" csv:"randomUserAgent"`

	// MaxDepth limits the recursion depth of visited URLs.
	// Set it to 0 for infinite recursion (default).
	MaxDepth int `default:"0" flag:"max-depth" yaml:"max_depth" toml:"max_depth" xml:"maxDepth" ini:"maxDepth" csv:"maxDepth" json:"max_depth" yaml:"max_depth" toml:"max_depth" xml:"maxDepth" ini:"maxDepth" csv:"maxDepth"`

	// AllowURLRevisit allows multiple downloads of the same URL
	AllowURLRevisit bool `default:"false" flag:"allow-url-revisit" yaml:"allow_url_revisit" toml:"allow_url_revisit" xml:"allowURLRevisit" ini:"allowURLRevisit" csv:"allowURLRevisit" json:"allow_url_revisit" yaml:"allow_url_revisit" toml:"allow_url_revisit" xml:"allowURLRevisit" ini:"allowURLRevisit" csv:"allowURLRevisit"`

	// IgnoreRobotsTxt allows the Collector to ignore any restrictions set by
	// the target host"s robots.txt file.  See http://www.robotstxt.org/ for more information.
	IgnoreRobotsTxt bool `default:"true" flag:"ignore-robots-txt" yaml:"ignore_robots_txt" toml:"ignore_robots_txt" xml:"ignoreRobotsTxt" ini:"ignoreRobotsTxt" csv:"ignoreRobotsTxt" json:"ignore_robots_txt" yaml:"ignore_robots_txt" toml:"ignore_robots_txt" xml:"ignoreRobotsTxt" ini:"ignoreRobotsTxt" csv:"ignoreRobotsTxt"`

	// CurrentMode
	CurrentMode string `default:"async" flag:"collector-mode" yaml:"current_mode" toml:"current_mode" xml:"CurrentMode" ini:"CurrentMode" csv:"CurrentMode" json:"current_mode" yaml:"current_mode" toml:"current_mode" xml:"CurrentMode" ini:"CurrentMode" csv:"CurrentMode"`

	// Modes
	Modes *collector_mode.Mode `json:"modes" yaml:"modes" toml:"modes" xml:"modes" ini:"modes" csv:"modes"`

	// Cache
	Cache *collector_cache.Config `json:"cache" yaml:"cache" toml:"cache" xml:"cache" ini:"cache" csv:"cache"`

	// Transport
	Transport *http_transport.Config `json:"transport" yaml:"transport" toml:"transport" xml:"transport" ini:"transport" csv:"transport"`

	// Proxy
	Proxy *http_proxy.Config `json:"proxy" yaml:"proxy" toml:"proxy" xml:"proxy" ini:"proxy" csv:"proxy"`

	// Filters
	Filters *filter_content.Filter `json:"filters" yaml:"filters" toml:"filters" xml:"filters" ini:"filters" csv:"filters"`

	// Sitemap
	Sitemap *loader_sitemap.Config `json:"collector" yaml:"collector" toml:"collector" xml:"collector" ini:"collector" csv:"collector"`

	//-- END
}
