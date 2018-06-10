package legacy

import (
	"net/http"
	"regexp"
)

// private variables
var (

	// fullyQualifiedPath
	fullyQualifiedPath bool = false

	// collectorBaseDir
	collectorBaseDir string

	// collectorWorkDir
	collectorWorkDir string

	// collectorAppName
	collectorAppName string = DEFAULT_APP_NAME
)

// Collector provides the scraper instance for a scraping job
type Compatible struct {

	// UserAgent is the User-Agent string used by HTTP requests
	UserAgent string

	// MaxDepth limits the recursion depth of visited URLs.
	// Set it to 0 for infinite recursion (default).
	MaxDepth int

	// AllowedDomains is a domain whitelist.
	// Leave it blank to allow any domains to be visited
	AllowedDomains []string

	// DisallowedDomains is a domain blacklist.
	DisallowedDomains []string

	// DisallowedURLFilters is a list of regular expressions which restricts
	// visiting URLs. If any of the rules matches to a URL the
	// request will be stopped. DisallowedURLFilters will
	// be evaluated before URLFilters
	// Leave it blank to allow any URLs to be visited
	DisallowedURLFilters []*regexp.Regexp

	// URLFilters is a list of regular expressions which restricts
	// visiting URLs. If any of the rules matches to a URL the
	// request won't be stopped. DisallowedURLFilters will
	// be evaluated before URLFilters
	// Leave it blank to allow any URLs to be visited
	URLFilters []*regexp.Regexp

	// AllowURLRevisit allows multiple downloads of the same URL
	AllowURLRevisit bool

	// MaxBodySize is the limit of the retrieved response body in bytes.
	// 0 means unlimited.
	// The default value for MaxBodySize is 10MB (10 * 1024 * 1024 bytes).
	MaxBodySize int

	// CacheDir specifies a location where GET requests are cached as files.
	// When it's not defined, caching is disabled.
	CacheDir string

	// IgnoreRobotsTxt allows the Collector to ignore any restrictions set by
	// the target host's robots.txt file.  See http://www.robotstxt.org/ for more
	// information.
	IgnoreRobotsTxt bool

	// Async turns on asynchronous network communication. Use Collector.Wait() to
	// be sure all requests have been finished.
	Async bool

	// ParseHTTPErrorResponse allows parsing HTTP responses with non 2xx status codes.
	// By default, Colly parses only successful HTTP responses. Set ParseHTTPErrorResponse
	// to true to enable it.
	ParseHTTPErrorResponse bool

	// ID is the unique identifier of a collector
	ID uint32

	// DetectCharset can enable character encoding detection for non-utf8 response bodies
	// without explicit charset declaration. This feature uses https://github.com/saintfish/chardet
	DetectCharset bool

	// RedirectHandler allows control on how a redirect will be managed
	RedirectHandler func(req *http.Request, via []*http.Request) error

	//-- End
}
