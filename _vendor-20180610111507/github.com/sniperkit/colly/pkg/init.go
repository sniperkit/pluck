package colly

import (
	"net/http/cookiejar"
	"regexp"
	"sync"
	"sync/atomic"

	// core
	cfg "github.com/sniperkit/colly/pkg/config"
	debug "github.com/sniperkit/colly/pkg/debug"
	storage "github.com/sniperkit/colly/pkg/storage"

	// plugins
	robotstxt "github.com/sniperkit/colly/plugins/data/format/robotstxt"
)

// NewCollector creates a new Collector instance with cfg.Default configuration
func NewCollector(options ...func(*Collector)) *Collector {
	c := &Collector{}
	c.Init()
	for _, f := range options {
		f(c)
	}
	c.parseSettingsFromEnv()
	return c
}

// NewCollector creates a new Collector instance with cfg.Default configuration
func NewCollectorWithConfig(cfg *cfg.Config) (c *Collector) {
	c = &Collector{}
	if cfg != nil {

		//// ---

		// Cache Storage
		c.store = &storage.InMemoryStorage{}
		c.store.Init()
		c.MaxBodySize = cfg.Filters.Response.MaxBodySize
		c.backend = &httpBackend{}

		// Requests
		jar, _ := cookiejar.New(nil)
		c.backend.Init(jar)
		c.backend.Client.CheckRedirect = c.checkRedirectFunc()
		c.wg = &sync.WaitGroup{}
		c.lock = &sync.RWMutex{}
		c.robotsMap = make(map[string]*robotstxt.RobotsData)
		c.IgnoreRobotsTxt = cfg.Collector.IgnoreRobotsTxt
		c.ID = atomic.AddUint32(&collectorCounter, 1)

		// Filters
		c.AllowedDomains = cfg.Filters.Whitelists.Domains
		c.DisallowedDomains = cfg.Filters.Blacklists.Domains

		c.AllowURLRevisit = cfg.Collector.AllowURLRevisit

		c.CacheDir = cfg.Collector.Cache.Store.Directory

		// cfg.Blacklists.Domains

		c.MaxDepth = cfg.Collector.MaxDepth
		c.ParseHTTPErrorResponse = cfg.Filters.Response.ParseHTTPErrorResponse

		if cfg.Collector.UserAgent != "" {
			c.UserAgent = cfg.Collector.UserAgent
		} else {
			c.UserAgent = `Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36`
		}

		if cfg.Collector.CurrentMode == "async" {
			c.Async = true
		}

		// Advanced features
		c.DetectCharset = cfg.Filters.Response.DetectCharset
		c.DetectTabular = cfg.Filters.Response.DetectTabular
		c.DetectMimeType = cfg.Filters.Response.DetectMimeType

		c.DebugMode = cfg.App.DebugMode
		c.VerboseMode = cfg.App.VerboseMode

	} else {

		c.Init()
		c.parseSettingsFromEnv()

	}

	return
}

// UserAgent sets the user agent used by the Collector.
func UserAgent(ua string) func(*Collector) {
	return func(c *Collector) {
		c.UserAgent = ua
	}
}

// MaxDepth limits the recursion depth of visited URLs.
func MaxDepth(depth int) func(*Collector) {
	return func(c *Collector) {
		c.MaxDepth = depth
	}
}

// AllowTabular enables OnTAB event callback by the Collector.
func AllowTabular(status bool) func(*Collector) {
	return func(c *Collector) {
		c.AllowTabular = status
	}
}

// AllowedDomains sets the domain whitelist used by the Collector.
func AllowedDomains(domains ...string) func(*Collector) {
	return func(c *Collector) {
		c.AllowedDomains = domains
	}
}

// ParseHTTPErrorResponse allows parsing responses with HTTP errors
func ParseHTTPErrorResponse() func(*Collector) {
	return func(c *Collector) {
		c.ParseHTTPErrorResponse = cfg.DefaultParseHTTPErrorResponse
		// c.Config.ParseHTTPErrorResponse = cfg.DefaultParseHTTPErrorResponse
	}
}

// ForceDir specifies that the program will try to create missing storage directories.
func ForceDir() func(*Collector) {
	return func(c *Collector) {
		c.ForceDir = cfg.DefaultForceDir
	}
}

// ForceDirRecursive specifies that the program will try to create missing storage directories recursively.
func ForceDirRecursive() func(*Collector) {
	return func(c *Collector) {
		c.ForceDirRecursive = cfg.DefaultForceDirRecursive
	}
}

// DisallowedDomains sets the domain blacklist used by the Collector.
func DisallowedDomains(domains ...string) func(*Collector) {
	return func(c *Collector) {
		c.DisallowedDomains = domains
	}
}

// DisallowedURLFilters sets the list of regular expressions which restricts
// visiting URLs. If any of the rules matches to a URL the request will be stopped.
func DisallowedURLFilters(filters ...*regexp.Regexp) func(*Collector) {
	return func(c *Collector) {
		c.DisallowedURLFilters = filters
		// c.Config.DisallowedURLFilters = filters
	}
}

// URLFilters sets the list of regular expressions which restricts
// visiting URLs. If any of the rules matches to a URL the request won't be stopped.
func URLFilters(filters ...*regexp.Regexp) func(*Collector) {
	return func(c *Collector) {
		c.URLFilters = filters
	}
}

// AllowURLRevisit instructs the Collector to allow multiple downloads of the same URL
func AllowURLRevisit() func(*Collector) {
	return func(c *Collector) {
		c.AllowURLRevisit = cfg.DefaultAllowURLRevisit
	}
}

// MaxBodySize sets the limit of the retrieved response body in bytes.
func MaxBodySize(sizeInBytes int) func(*Collector) {
	return func(c *Collector) {
		c.MaxBodySize = sizeInBytes
	}
}

// CacheDir specifies the location where GET requests are cached as files.
func CacheDir(path string) func(*Collector) {
	return func(c *Collector) {
		c.CacheDir = path
	}
}

// IgnoreRobotsTxt instructs the Collector to ignore any restrictions
// set by the target host's robots.txt file.
func IgnoreRobotsTxt() func(*Collector) {
	return func(c *Collector) {
		c.IgnoreRobotsTxt = cfg.DefaultIgnoreRobotsTxt
	}
}

// RandomUserAgent
func RandomUserAgent() func(*Collector) {
	return func(c *Collector) {
		c.RandomUserAgent = cfg.DefaultRandomUserAgent
	}
}

// ID sets the unique identifier of the Collector.
func ID(id uint32) func(*Collector) {
	return func(c *Collector) {
		c.ID = id
	}
}

// Async turns on asynchronous network requests.
func Async(a bool) func(*Collector) {
	return func(c *Collector) {
		c.Async = a
	}
}

// DetectCharset enables character encoding detection for non-utf8 response bodies
// without explicit charset declaration. This feature uses https://github.com/saintfish/chardet
func DetectCharset() func(*Collector) {
	return func(c *Collector) {
		c.DetectCharset = cfg.DefaultDetectCharset
	}
}

// DetectMimeType enables mime type detection
func DetectMimeType() func(*Collector) {
	return func(c *Collector) {
		c.DetectMimeType = cfg.DefaultDetectMimeType
	}
}

// DetectTabular
func DetectTabular() func(*Collector) {
	return func(c *Collector) {
		c.DetectTabular = cfg.DefaultDetectMimeType
	}
}

// DebugMode enables text-based content summarization.
func DebugMode() func(*Collector) {
	return func(c *Collector) {
		c.DebugMode = cfg.DefaultDebugMode
	}
}

// VerboseMode enables text-based content summarization.
func VerboseMode() func(*Collector) {
	return func(c *Collector) {
		c.VerboseMode = cfg.DefaultVerboseMode
	}
}

// Debugger sets the debugger used by the Collector.
func Debugger(d debug.Debugger) func(*Collector) {
	return func(c *Collector) {
		d.Init()
		c.debugger = d
	}
}
