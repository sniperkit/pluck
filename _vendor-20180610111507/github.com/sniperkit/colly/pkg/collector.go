// Copyright 2018 Adam Tauber
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package colly implements a HTTP scraping framework
package colly

import (
	"bytes"
	"crypto/rand"
	"errors"
	"fmt"
	"hash/fnv"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/net/html"
	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"

	"github.com/PuerkitoBio/goquery"

	// collector - core
	cfg "github.com/sniperkit/colly/pkg/config"
	debug "github.com/sniperkit/colly/pkg/debug"
	metric "github.com/sniperkit/colly/pkg/metric"
	storage "github.com/sniperkit/colly/pkg/storage"

	// content - transform
	tabular "github.com/sniperkit/colly/plugins/data/transform/tabular"
	sanitize "github.com/sniperkit/colly/plugins/data/transform/text/sanitize"

	// format
	robotstxt "github.com/sniperkit/colly/plugins/data/format/robotstxt"

	// encoding - iterators
	jsoniter "github.com/json-iterator/go"

	// format - operation
	htmlquery "github.com/sniperkit/colly/plugins/data/extract/query/html"
	jsonquery "github.com/sniperkit/colly/plugins/data/extract/query/json"
	xmlquery "github.com/sniperkit/colly/plugins/data/extract/query/xml"
)

var (
	collectorCounter uint32
	collectorConfig  *cfg.Config
	collectorMetrics *metric.Snapshot
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// Collector provides the scraper instance for a scraping job
type Collector struct {

	// Collector's settings
	cfg.Config

	//////////////////////////////////////////////////
	///// Collector - info
	//////////////////////////////////////////////////

	// ID is the unique identifier of a collector
	ID uint32 `default:"colly" json:"identifier" yaml:"identifier" toml:"identifier" xml:"identifier" ini:"identifier" csv:"identifier"`

	// Title/name of the current crawling campaign
	Title string `default:"Colly - Web Scraper" json:"title" yaml:"title" toml:"title" xml:"title" ini:"title" csv:"title"`

	// UserAgent is the User-Agent string used by HTTP requests
	UserAgent string `default:"colly - https://github.com/sniperkit/colly" json:"user_agent" yaml:"user_agent" toml:"user_agent" xml:"userAgent" ini:"userAgent" csv:"userAgent"`

	// RandomUserAgent specifies to generate a random User-Agent string for all HTTP requests
	RandomUserAgent bool `default:"false" json:"random_user_agent" yaml:"random_user_agent" toml:"random_user_agent" xml:"randomUserAgent" ini:"randomUserAgent" csv:"randomUserAgent"`

	//////////////////////////////////////////////////
	///// Collector - crawling parameters
	//////////////////////////////////////////////////

	// Async turns on asynchronous network communication. Use Collector.Wait() to be sure all requests have been finished.
	Async bool `default:"false" json:"async" yaml:"async" toml:"async" xml:"async" ini:"async" csv:"async"`

	// MaxDepth limits the recursion depth of visited URLs.
	// Set it to 0 for infinite recursion (default).
	MaxDepth int `default:"0" json:"max_depth" yaml:"max_depth" toml:"max_depth" xml:"maxDepth" ini:"maxDepth" csv:"maxDepth"`

	// AllowURLRevisit allows multiple downloads of the same URL
	AllowURLRevisit bool `default:"false" json:"allow_url_revisit" yaml:"allow_url_revisit" toml:"allow_url_revisit" xml:"allowURLRevisit" ini:"allowURLRevisit" csv:"allowURLRevisit"`

	// IgnoreRobotsTxt allows the Collector to ignore any restrictions set by
	// the target host"s robots.txt file.  See http://www.robotstxt.org/ for more information.
	IgnoreRobotsTxt bool `default:"true" json:"ignore_robots_txt" yaml:"ignore_robots_txt" toml:"ignore_robots_txt" xml:"ignoreRobotsTxt" ini:"ignoreRobotsTxt" csv:"ignoreRobotsTxt"`

	//////////////////////////////////////////////////
	///// Request - Filtering parameters
	//////////////////////////////////////////////////

	////// Not exportable attributes

	// AllowedDomains is a domain whitelist.
	// Leave it blank to allow any domains to be visited
	AllowedDomains []string `json:"allowed_domains" yaml:"allowed_domains" toml:"allowed_domains" xml:"allowedDomains" ini:"allowedDomains" csv:"AllowedDomains"`

	// DisallowedDomains is a domain blacklist.
	DisallowedDomains []string `json:"disallowed_domains" yaml:"disallowed_domains" toml:"disallowed_domains" xml:"disallowedDomains" ini:"disallowedDomains" csv:"DisallowedDomains"`

	// DisallowedURLFilters is a list of regular expressions which restricts
	// visiting URLs. If any of the rules matches to a URL the
	// request will be stopped. DisallowedURLFilters will
	// be evaluated before URLFilters
	// Important: Leave it blank to allow any URLs to be visited
	DisallowedURLFilters []*regexp.Regexp `json:"-" yaml:"-" toml:"-" xml:"-" ini:"-" csv:"-"`

	// URLFilters is a list of regular expressions which restricts
	// visiting URLs. If any of the rules matches to a URL the
	// request won"t be stopped. DisallowedURLFilters will
	// be evaluated before URLFilters
	// Important: Leave it blank to allow any URLs to be visited
	URLFilters []*regexp.Regexp `json:"-" yaml:"-" toml:"-" xml:"-" ini:"-" csv:"-"`

	// MaxBodySize is the limit of the retrieved response body in bytes.
	// 0 means unlimited.
	// The default value for MaxBodySize is 10MB (10 * 1024 * 1024 bytes).
	MaxBodySize int `default:"0" json:"max_body_size" yaml:"max_body_size" toml:"max_body_size" xml:"maxBodySize" ini:"maxBodySize" csv:"MaxBodySize"`

	// AllowTabular
	AllowTabular bool `default:"false" json:"allow_tabular" yaml:"allow_tabular" toml:"allowTabular" xml:"allowTabular" ini:"allowTabular" csv:"AllowTabular"`

	// UseJsonParser
	UseJsonParser string `json:"useJsonParser" yaml:"useJsonParser" toml:"useJsonParser" xml:"useJsonParser" ini:"useJsonParser" csv:"UseJsonParser"`
	// UseJsonParser JsonParser `json:"useJsonParser" yaml:"useJsonParser" toml:"useJsonParser" xml:"useJsonParser" ini:"useJsonParser" csv:"UseJsonParser"`

	//////////////////////////////////////////////////
	///// Response processing
	//////////////////////////////////////////////////

	// ParseHTTPErrorResponse allows parsing HTTP responses with non 2xx status codes.
	// By default, Colly parses only successful HTTP responses. Set ParseHTTPErrorResponse to true to enable it.
	ParseHTTPErrorResponse bool `default:"true" json:"parse_http_error_response" yaml:"parse_http_error_response" toml:"parse_http_error_response" xml:"parseHTTPErrorResponse" ini:"parseHTTPErrorResponse" csv:"parseHTTPErrorResponse"`

	// DetectCharset can enable character encoding detection for non-utf8 response bodies
	// without explicit charset declaration. This feature uses https://github.com/saintfish/chardet
	DetectCharset bool `default:"true" json:"detect_charset" yaml:"detect_charset" toml:"detect_charset" xml:"detectCharset" ini:"detectCharset" csv:"DetectCharset"`

	// DetectMimeType
	DetectMimeType bool `default:"true" json:"detect_mime_type" yaml:"detect_mime_type" toml:"detect_mime_type" xml:"detectMimeType" ini:"detectMimeType" csv:"detectMimeType"`

	// DetectTabular
	DetectTabular bool `default:"true" json:"detect_tabular_data" yaml:"detect_tabular_data" toml:"detect_tabular_data" xml:"detectTabularData" ini:"detectTabularData" csv:"DetectTabularData"`

	// XDGBaseDir
	XDGBaseDir string `json:"xdg_base_dir" yaml:"xdg_base_dir" toml:"xdg_base_dir" xml:"xdgBaseDir" ini:"xdgBaseDir" csv:"XDGBaseDir"`

	// BaseDirectory
	BaseDir string `json:"base_dir" yaml:"base_dir" toml:"base_dir" xml:"baseDir" ini:"baseDir" csv:"BaseDir"`

	// LogsDirectory
	LogsDir string `json:"logs_dir" yaml:"logs_dir" toml:"logs_dir" xml:"logsDir" ini:"logsDir" csv:"LogsDir"`

	// CacheDir specifies a location where GET requests are cached as files.
	// When it"s not defined, caching is disabled.
	CacheDir string `default:"./shared/storage/cache/http/backends/internal" json:"cache_dir" yaml:"cache_dir" toml:"cache_dir" xml:"cacheDir" ini:"cacheDir" csv:"CacheDir"`

	// ExportDir
	ExportDir string `default:"./shared/exports" json:"export_dir" yaml:"export_dir" toml:"export_dir" xml:"exportDir" ini:"exportDir" csv:"ExportDir"`

	// ForceDir specifies that the program will try to create missing storage directories.
	ForceDir bool `default:"true" json:"force_dir" yaml:"force_dir" toml:"force_dir" xml:"forceDir" ini:"forceDir" csv:"ForceDir"`

	// ForceDirRecursive specifies that the program will try to create missing storage directories recursively.
	ForceDirRecursive bool `default:"true" json:"force_dir_recursive" yaml:"force_dir_recursive" toml:"force_dir_recursive" xml:"forceDirRecursive" ini:"forceDirRecursive" csv:"ForceDirRecursive"`

	//////////////////////////////////////////////////
	///// Debug mode
	//////////////////////////////////////////////////

	// DebugMode
	DebugMode bool `default:"false" json:"debug_mode" yaml:"debug_mode" toml:"debug_mode" xml:"debugMode" ini:"debugMode" csv:"DebugMode"`

	// VerboseMode
	VerboseMode bool `default:"verbose_mode" json:"verbose_mode" yaml:"verbose_mode" toml:"verbose_mode" xml:"verboseMode" ini:"verboseMode" csv:"VerboseMode"`

	//////////////////////////////////////////////////
	///// Dashboard TUI (terminal ui only)
	//////////////////////////////////////////////////

	// IsDashboard
	DashboardMode bool `default:"true" json:"dashboard_mode" yaml:"dashboard_mode" toml:"dashboard_mode" xml:"dashboardMode" ini:"dashboardMode" csv:"dashboardMode"`

	//////////////////////////////////////////////////
	///// Export application"s config to local file
	//////////////////////////////////////////////////

	// AllowExportConfigSchema
	AllowExportConfigSchema bool `default:"true" json:"-" yaml:"-" toml:"-" xml:"-" ini:"-" csv:"-"`

	// AllowExportConfigAll
	AllowExportConfigAutoload bool `default:"false" json:"-" yaml:"-" toml:"-" xml:"-" ini:"-" csv:"-"`

	//////////////////////////////////////////////////
	///// Experimental stuff
	//////////////////////////////////////////////////

	// RedirectHandler allows control on how a redirect will be managed
	RedirectHandler func(req *http.Request, via []*http.Request) error

	// not exported attributes
	store         storage.Storage
	debugger      debug.Debugger
	robotsMap     map[string]*robotstxt.RobotsData
	requestCount  uint32
	responseCount uint32
	backend       *httpBackend
	wg            *sync.WaitGroup
	lock          *sync.RWMutex

	// content callbacks
	htmlCallbacks []*htmlCallbackContainer
	xmlCallbacks  []*xmlCallbackContainer
	jsonCallbacks []*jsonCallbackContainer
	tabCallbacks  []*tabCallbackContainer

	// collector callbacks
	requestCallbacks  []RequestCallback
	responseCallbacks []ResponseCallback
	errorCallbacks    []ErrorCallback
	scrapedCallbacks  []ScrapedCallback
}

// Init initializes the Collector's private variables and sets default configuration for the Collector
func (c *Collector) Init() {
	c.UserAgent = "colly - https://github.com/sniperkit/colly/pkg"
	c.MaxDepth = 0
	c.store = &storage.InMemoryStorage{}
	c.store.Init()
	c.MaxBodySize = 10 * 1024 * 1024
	c.backend = &httpBackend{}
	jar, _ := cookiejar.New(nil)
	c.backend.Init(jar)
	c.backend.Client.CheckRedirect = c.checkRedirectFunc()
	c.wg = &sync.WaitGroup{}
	c.lock = &sync.RWMutex{}
	c.robotsMap = make(map[string]*robotstxt.RobotsData)
	c.IgnoreRobotsTxt = true
	c.ID = atomic.AddUint32(&collectorCounter, 1)
	c.AllowURLRevisit = false
}

func (c *Collector) IsDebug() bool {
	if c.debugger != nil {
		return true
	}
	return false
}

// Appengine will replace the Collector's backend http.Client
// With an Http.Client that is provided by appengine/urlfetch
// This function should be used when the scraper is initiated
// by a http.Request to Google App Engine
func (c *Collector) Appengine(req *http.Request) {
	ctx := appengine.NewContext(req)
	client := urlfetch.Client(ctx)
	client.Jar = c.backend.Client.Jar
	client.CheckRedirect = c.backend.Client.CheckRedirect
	client.Timeout = c.backend.Client.Timeout

	c.backend.Client = client
}

// Visit starts Collector's collecting job by creating a
// request to the URL specified in parameter.
// Visit also calls the previously provided callbacks
func (c *Collector) Visit(URL string) error {
	return c.scrape(URL, "GET", 1, nil, nil, nil, true)
}

// Post starts a collector job by creating a POST request.
// Post also calls the previously provided callbacks
func (c *Collector) Post(URL string, requestData map[string]string) error {
	return c.scrape(URL, "POST", 1, createFormReader(requestData), nil, nil, true)
}

// PostMultipart starts a collector job by creating a Multipart POST request
// with raw binary data.  PostMultipart also calls the previously provided callbacks
func (c *Collector) PostMultipart(URL string, requestData map[string][]byte) error {
	boundary := randomBoundary()
	hdr := http.Header{}
	hdr.Set("Content-Type", "multipart/form-data; boundary="+boundary)
	hdr.Set("User-Agent", c.UserAgent)
	return c.scrape(URL, "POST", 1, createMultipartReader(boundary, requestData), nil, hdr, true)
}

// Request starts a collector job by creating a custom HTTP request
// where method, context, headers and request data can be specified.
// Set requestData, ctx, hdr parameters to nil if you don't want to use them.
// Valid methods:
//   - "GET"
//   - "POST"
//   - "PUT"
//   - "DELETE"
//   - "PATCH"
//   - "OPTIONS"
func (c *Collector) Request(method, URL string, requestData io.Reader, ctx *Context, hdr http.Header) error {
	return c.scrape(URL, method, 1, requestData, ctx, hdr, true)
}

// SetDebugger attaches a debugger to the collector
func (c *Collector) SetDebugger(d debug.Debugger) {
	d.Init()
	c.debugger = d
}

// UnmarshalRequest creates a Request from serialized data
func (c *Collector) UnmarshalRequest(r []byte) (*Request, error) {
	req := &serializableRequest{}
	err := json.Unmarshal(r, req)
	if err != nil {
		return nil, err
	}

	u, err := url.Parse(req.URL)
	if err != nil {
		return nil, err
	}

	ctx := NewContext()
	for k, v := range req.Ctx {
		ctx.Put(k, v)
	}

	return &Request{
		Method:    req.Method,
		URL:       u,
		Body:      bytes.NewReader(req.Body),
		Ctx:       ctx,
		ID:        atomic.AddUint32(&c.requestCount, 1),
		Headers:   &http.Header{},
		collector: c,
	}, nil
}

func (c *Collector) scrape(u, method string, depth int, requestData io.Reader, ctx *Context, hdr http.Header, checkRevisit bool) error {
	if err := c.requestCheck(u, method, depth, checkRevisit); err != nil {
		return err
	}
	parsedURL, err := url.Parse(u)
	if err != nil {
		return err
	}
	if parsedURL.Scheme == "" {
		parsedURL.Scheme = DEFAULT_HTTP_SCHEME
	}
	if !c.isDomainAllowed(parsedURL.Host) {
		return ErrForbiddenDomain
	}
	if !c.IgnoreRobotsTxt {
		if err = c.checkRobots(parsedURL); err != nil {
			return err
		}
	}
	if hdr == nil {
		hdr = http.Header{"User-Agent": []string{c.UserAgent}}
	}
	rc, ok := requestData.(io.ReadCloser)
	if !ok && requestData != nil {
		rc = ioutil.NopCloser(requestData)
	}
	req := &http.Request{
		Method:     method,
		URL:        parsedURL,
		Proto:      DEFAULT_HTTP_REQUEST_PROTO,
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     hdr,
		Body:       rc,
		Host:       parsedURL.Host,
	}
	setRequestBody(req, requestData)
	u = parsedURL.String()
	c.wg.Add(1)
	if c.Async {
		go c.fetch(u, method, depth, requestData, ctx, hdr, req)
		return nil
	}
	return c.fetch(u, method, depth, requestData, ctx, hdr, req)
}

func setRequestBody(req *http.Request, body io.Reader) {
	if body != nil {
		switch v := body.(type) {
		case *bytes.Buffer:
			req.ContentLength = int64(v.Len())
			buf := v.Bytes()
			req.GetBody = func() (io.ReadCloser, error) {
				r := bytes.NewReader(buf)
				return ioutil.NopCloser(r), nil
			}
		case *bytes.Reader:
			req.ContentLength = int64(v.Len())
			snapshot := *v
			req.GetBody = func() (io.ReadCloser, error) {
				r := snapshot
				return ioutil.NopCloser(&r), nil
			}
		case *strings.Reader:
			req.ContentLength = int64(v.Len())
			snapshot := *v
			req.GetBody = func() (io.ReadCloser, error) {
				r := snapshot
				return ioutil.NopCloser(&r), nil
			}
		}
		if req.GetBody != nil && req.ContentLength == 0 {
			req.Body = http.NoBody
			req.GetBody = func() (io.ReadCloser, error) { return http.NoBody, nil }
		}
	}
}

func (c *Collector) fetch(u, method string, depth int, requestData io.Reader, ctx *Context, hdr http.Header, req *http.Request) error {
	defer c.wg.Done()
	if ctx == nil {
		ctx = NewContext()
	}
	request := &Request{
		URL:       req.URL,
		Headers:   &req.Header,
		Ctx:       ctx,
		Depth:     depth,
		Method:    method,
		Body:      requestData,
		collector: c,
		ID:        atomic.AddUint32(&c.requestCount, 1),
	}

	c.handleOnRequest(request)

	if request.abort {
		return nil
	}

	if method == "POST" && req.Header.Get("Content-Type") == "" {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	}

	if req.Header.Get("Accept") == "" {
		req.Header.Set("Accept", "*/*")
	}

	origURL := req.URL
	response, err := c.backend.Cache(req, c.MaxBodySize, c.CacheDir)
	if err := c.handleOnError(response, err, request, ctx); err != nil {
		return err
	}
	if req.URL != origURL {
		request.URL = req.URL
		request.Headers = &req.Header
	}
	atomic.AddUint32(&c.responseCount, 1)
	response.Ctx = ctx
	response.Request = request

	err = response.fixCharset(c.DetectCharset, request.ResponseCharacterEncoding)
	if err != nil {
		return err
	}

	c.handleOnResponse(response)

	isTabularContent, format := isTabular(response.Headers.Get("Content-Type"), request.URL.String())

	if c.AllowTabular && isTabularContent {
		err = c.handleOnTAB(response)
		if err != nil {
			c.handleOnError(response, err, request, ctx)
		}

	} else {

		switch format {
		case "json":
			err = c.handleOnJSON(response)
			if err != nil {
				c.handleOnError(response, err, request, ctx)
			}

		case "xml":
			err = c.handleOnXML(response)
			if err != nil {
				c.handleOnError(response, err, request, ctx)
			}

		case "html":
			fallthrough

		default:
			err = c.handleOnHTML(response)
			if err != nil {
				c.handleOnError(response, err, request, ctx)
			}

		}

	}

	c.handleOnScraped(response)

	return err
}

func (c *Collector) requestCheck(u, method string, depth int, checkRevisit bool) error {
	if u == "" {
		return ErrMissingURL
	}
	if c.MaxDepth > 0 && c.MaxDepth < depth {
		return ErrMaxDepth
	}
	if len(c.DisallowedURLFilters) > 0 {
		if isMatchingFilter(c.DisallowedURLFilters, []byte(u)) {
			return ErrForbiddenURL
		}
	}
	if len(c.URLFilters) > 0 {
		if !isMatchingFilter(c.URLFilters, []byte(u)) {
			return ErrNoURLFiltersMatch
		}
	}
	if checkRevisit && !c.AllowURLRevisit && method == "GET" {
		h := fnv.New64a()
		h.Write([]byte(u))
		uHash := h.Sum64()
		visited, err := c.store.IsVisited(uHash)
		if err != nil {
			return err
		}
		if visited {
			return ErrAlreadyVisited
		}
		return c.store.Visited(uHash)
	}
	return nil
}

func (c *Collector) isDomainAllowed(domain string) bool {
	for _, d2 := range c.DisallowedDomains {
		if d2 == domain {
			return false
		}
	}
	if c.AllowedDomains == nil || len(c.AllowedDomains) == 0 {
		return true
	}
	for _, d2 := range c.AllowedDomains {
		if d2 == domain {
			return true
		}
	}
	return false
}

func (c *Collector) checkRobots(u *url.URL) error {
	c.lock.RLock()
	robot, ok := c.robotsMap[u.Host]
	c.lock.RUnlock()

	if !ok {
		// no robots file cached
		resp, err := c.backend.Client.Get(u.Scheme + "://" + u.Host + "/robots.txt")
		if err != nil {
			return err
		}
		robot, err = robotstxt.FromResponse(resp)
		if err != nil {
			return err
		}
		c.lock.Lock()
		c.robotsMap[u.Host] = robot
		c.lock.Unlock()
	}

	uaGroup := robot.FindGroup(c.UserAgent)
	if uaGroup == nil {
		return nil
	}

	if !uaGroup.Test(u.EscapedPath()) {
		return ErrRobotsTxtBlocked
	}
	return nil
}

// String is the text representation of the collector.
// It contains useful debug information about the collector's internals
func (c *Collector) String() string {
	return fmt.Sprintf(
		"Requests made: %d (%d responses) | Callbacks: OnRequest: %d, OnHTML: %d, OnResponse: %d, OnError: %d",
		c.requestCount,
		c.responseCount,
		len(c.requestCallbacks),
		len(c.htmlCallbacks),
		len(c.responseCallbacks),
		len(c.errorCallbacks),
	)
}

// Wait returns when the collector jobs are finished
func (c *Collector) Wait() {
	c.wg.Wait()
}

// OnRequest registers a function. Function will be executed on every
// request made by the Collector
func (c *Collector) OnRequest(f RequestCallback) {
	c.lock.Lock()
	if c.requestCallbacks == nil {
		c.requestCallbacks = make([]RequestCallback, 0, 4)
	}
	c.requestCallbacks = append(c.requestCallbacks, f)
	c.lock.Unlock()
}

// OnResponse registers a function. Function will be executed on every response
func (c *Collector) OnResponse(f ResponseCallback) {
	c.lock.Lock()
	if c.responseCallbacks == nil {
		c.responseCallbacks = make([]ResponseCallback, 0, 4)
	}
	c.responseCallbacks = append(c.responseCallbacks, f)
	c.lock.Unlock()
}

// OnTAB registers a function. Function will be executed on every JSON
// element matched by the xpath parameter.
func (c *Collector) OnTAB(query string, f TABCallback) {
	c.lock.Lock()
	if c.tabCallbacks == nil {
		c.tabCallbacks = make([]*tabCallbackContainer, 0, 4)
	}
	c.tabCallbacks = append(c.tabCallbacks, &tabCallbackContainer{
		Query:    query,
		Function: f,
	})
	c.lock.Unlock()
}

// OnJSON registers a function. Function will be executed on every JSON
// element matched by the xpath parameter.
func (c *Collector) OnJSON(xPath string, f JSONCallback) {
	c.lock.Lock()
	if c.jsonCallbacks == nil {
		c.jsonCallbacks = make([]*jsonCallbackContainer, 0, 4)
	}
	c.jsonCallbacks = append(c.jsonCallbacks, &jsonCallbackContainer{
		Query:    xPath,
		Function: f,
	})
	c.lock.Unlock()
}

// OnHTML registers a function. Function will be executed on every HTML
// element matched by the GoQuery Selector parameter.
// GoQuery Selector is a selector used by https://github.com/PuerkitoBio/goquery
func (c *Collector) OnHTML(goquerySelector string, f HTMLCallback) {
	c.lock.Lock()
	if c.htmlCallbacks == nil {
		c.htmlCallbacks = make([]*htmlCallbackContainer, 0, 4)
	}
	c.htmlCallbacks = append(c.htmlCallbacks, &htmlCallbackContainer{
		Selector: goquerySelector,
		Function: f,
	})
	c.lock.Unlock()
}

// OnXML registers a function. Function will be executed on every XML
// element matched by the xpath Query parameter.
// xpath Query is used by https://github.com/antchfx/xmlquery
func (c *Collector) OnXML(xpathQuery string, f XMLCallback) {
	c.lock.Lock()
	if c.xmlCallbacks == nil {
		c.xmlCallbacks = make([]*xmlCallbackContainer, 0, 4)
	}
	c.xmlCallbacks = append(c.xmlCallbacks, &xmlCallbackContainer{
		Query:    xpathQuery,
		Function: f,
	})
	c.lock.Unlock()
}

// OnJSONDetach deregister a function. Function will not be execute after detached
func (c *Collector) OnJSONDetach(xPath string) {
	c.lock.Lock()
	deleteIdx := -1
	for i, cc := range c.jsonCallbacks {
		if cc.Query == xPath {
			deleteIdx = i
			break
		}
	}
	if deleteIdx != -1 {
		c.jsonCallbacks = append(c.jsonCallbacks[:deleteIdx], c.jsonCallbacks[deleteIdx+1:]...)
	}
	c.lock.Unlock()
}

// OnTABDetach deregister a function. Function will not be execute after detached
func (c *Collector) OnTABDetach(xPath string) {
	c.lock.Lock()
	deleteIdx := -1
	for i, cc := range c.tabCallbacks {
		if cc.Query == xPath {
			deleteIdx = i
			break
		}
	}
	if deleteIdx != -1 {
		c.tabCallbacks = append(c.tabCallbacks[:deleteIdx], c.tabCallbacks[deleteIdx+1:]...)
	}
	c.lock.Unlock()
}

// OnHTMLDetach deregister a function. Function will not be execute after detached
func (c *Collector) OnHTMLDetach(goquerySelector string) {
	c.lock.Lock()
	deleteIdx := -1
	for i, cc := range c.htmlCallbacks {
		if cc.Selector == goquerySelector {
			deleteIdx = i
			break
		}
	}
	if deleteIdx != -1 {
		c.htmlCallbacks = append(c.htmlCallbacks[:deleteIdx], c.htmlCallbacks[deleteIdx+1:]...)
	}
	c.lock.Unlock()
}

// OnXMLDetach deregister a function. Function will not be execute after detached
func (c *Collector) OnXMLDetach(xpathQuery string) {
	c.lock.Lock()
	deleteIdx := -1
	for i, cc := range c.xmlCallbacks {
		if cc.Query == xpathQuery {
			deleteIdx = i
			break
		}
	}
	if deleteIdx != -1 {
		c.xmlCallbacks = append(c.xmlCallbacks[:deleteIdx], c.xmlCallbacks[deleteIdx+1:]...)
	}
	c.lock.Unlock()
}

// OnError registers a function. Function will be executed if an error
// occurs during the HTTP request.
func (c *Collector) OnError(f ErrorCallback) {
	c.lock.Lock()
	if c.errorCallbacks == nil {
		c.errorCallbacks = make([]ErrorCallback, 0, 4)
	}
	c.errorCallbacks = append(c.errorCallbacks, f)
	c.lock.Unlock()
}

// OnScraped registers a function. Function will be executed after
// OnHTML, as a final part of the scraping.
func (c *Collector) OnScraped(f ScrapedCallback) {
	c.lock.Lock()
	if c.scrapedCallbacks == nil {
		c.scrapedCallbacks = make([]ScrapedCallback, 0, 4)
	}
	c.scrapedCallbacks = append(c.scrapedCallbacks, f)
	c.lock.Unlock()
}

// WithTransport allows you to set a custom http.RoundTripper (transport)
func (c *Collector) WithTransport(transport http.RoundTripper) {
	c.backend.Client.Transport = transport
}

// DisableCookies turns off cookie handling
func (c *Collector) DisableCookies() {
	c.backend.Client.Jar = nil
}

// SetCookieJar overrides the previously set cookie jar
func (c *Collector) SetCookieJar(j *cookiejar.Jar) {
	c.backend.Client.Jar = j
}

// SetRequestTimeout overrides the default timeout (10 seconds) for this collector
func (c *Collector) SetRequestTimeout(timeout time.Duration) {
	c.backend.Client.Timeout = timeout
}

// SetStorage overrides the default in-memory storage.
// Storage stores scraping related data like cookies and visited urls
func (c *Collector) SetStorage(s storage.Storage) error {
	if err := s.Init(); err != nil {
		return err
	}
	c.store = s
	c.backend.Client.Jar = createJar(s)
	return nil
}

// SetProxy sets a proxy for the collector. This method overrides the previously
// used http.Transport if the type of the transport is not http.RoundTripper.
// The proxy type is determined by the URL scheme. "http"
// and "socks5" are supported. If the scheme is empty,
// "http" is assumed.
func (c *Collector) SetProxy(proxyURL string) error {
	proxyParsed, err := url.Parse(proxyURL)
	if err != nil {
		return err
	}

	c.SetProxyFunc(http.ProxyURL(proxyParsed))

	return nil
}

// SetProxyFunc sets a custom proxy setter/switcher function.
// See built-in ProxyFuncs for more details.
// This method overrides the previously used http.Transport
// if the type of the transport is not http.RoundTripper.
// The proxy type is determined by the URL scheme. "http"
// and "socks5" are supported. If the scheme is empty,
// "http" is assumed.
func (c *Collector) SetProxyFunc(p ProxyFunc) {
	t, ok := c.backend.Client.Transport.(*http.Transport)
	if c.backend.Client.Transport != nil && ok {
		t.Proxy = p
	} else {
		c.backend.Client.Transport = &http.Transport{
			Proxy: p,
		}
	}
}

func createEvent(eventType string, requestID, collectorID uint32, kvargs map[string]string) *debug.Event {
	return &debug.Event{
		CollectorID: collectorID,
		RequestID:   requestID,
		Type:        eventType,
		Values:      kvargs,
	}
}

func (c *Collector) handleOnRequest(r *Request) {
	if c.debugger != nil {
		c.debugger.Event(createEvent("request", r.ID, c.ID, map[string]string{
			"url": r.URL.String(),
		}))
	}
	for _, f := range c.requestCallbacks {
		f(r)
	}
}

func (c *Collector) handleOnResponse(r *Response) {
	if c.debugger != nil {
		c.debugger.Event(createEvent("response", r.Request.ID, c.ID, map[string]string{
			"url":    r.Request.URL.String(),
			"status": http.StatusText(r.StatusCode),
		}))
	}
	for _, f := range c.responseCallbacks {
		f(r)
	}
}

var tabContentKeywords []string = []string{"json", "yaml", "xml", "csv", "tsv"}

func isTabularEncoding(header string) (ok bool, format string) {
	for _, match := range tabContentKeywords {
		if strings.Contains(strings.ToLower(header), match) {
			return true, match
		}
	}
	return false, ""
}

func isTabularExtension(url string) (ok bool, format string) {
	for _, match := range tabContentKeywords {
		ext := path.Ext(url)
		if strings.Contains(strings.ToLower(ext), match) {
			return true, match
		}
	}
	return false, ""
}

func isTabular(contentType string, requestURL string) (ok bool, format string) {
	isValid, format := isTabularEncoding(contentType)
	if !isValid {
		isValid, format = isTabularExtension(requestURL)
	}
	return isValid, format
}

func parseSliceQuery(query string, ds *tabular.Dataset) (lower int, upper int, err error) {

	// row[1:5] = dataset slice from row 1 to 5
	parts := strings.Split(query, ":")
	count := len(parts)

	switch {
	case count > 2:
		fallthrough
	case count <= 1:
		return -1, -1, ErrTabularInvalidQuery
	}

	lowerStr := parts[0]
	upperStr := parts[1]

	if lowerStr == "" {
		lower = 0

	} else {
		lower, err = strconv.Atoi(lowerStr)
		if err != nil {
			return -1, -1, err
		}

		// check validity
		switch {
		case lower < 0: // lower limit not inferior to 0
			lower = 0

		case lower > ds.Height(): // lower limit is not out of range
			lower = ds.Height() - 1

		}

	}

	if upperStr == "" {
		upper = ds.Height()
	} else {
		upper, err = strconv.Atoi(upperStr)
		if err != nil {
			return -1, lower, err
		}

		// check validity
		switch {
		case upper <= lower: // upper limit not inferior to lower limit
			return -1, upper, ErrTabularInvalidQuery

		case upper > ds.Height(): // upper limit is not out of range
			upper = ds.Height()

		}

	}

	return
}

/*
	// Loading formats supported:
	JSON (Sets + Books)
	YAML (Sets + Books)
	XML (Sets)
	CSV (Sets)
	TSV (Sets)
*/
func (c *Collector) handleOnTAB(resp *Response) error {

	// check if valid encoding format to load
	isValid, format := isTabular(resp.Headers.Get("Content-Type"), resp.Request.URL.String())

	if c.debugger != nil {
		c.debugger.Event(createEvent("tabular.check", resp.Request.ID, c.ID, map[string]string{
			"content-type":    resp.Headers.Get("Content-Type"),
			"url":             resp.Request.URL.String(),
			"is_valid":        fmt.Sprintf("%t", isValid),
			"format_slug":     format,
			"callbacks_count": fmt.Sprintf("%d", len(c.tabCallbacks)),
		}))
	}

	if len(c.tabCallbacks) == 0 || !isValid {
		return ErrNotValidTabularFormat
	}

	var err error
	var ds *tabular.Dataset
	switch format {
	case "json":

		switch c.UseJsonParser {
		// MXJ: Decode / encode XML to/from map[string]interface{} (or JSON); extract values with dot-notation paths and wildcards.
		case "mxj":
			ds, err = tabular.LoadMXJ(resp.Body)

		case "gjson":
			ds, err = tabular.LoadGJSON(resp.Body)

		case "json":
			fallthrough

		default:
			ds, err = tabular.LoadJSON(resp.Body)
		}

	case "yaml":
		ds, err = tabular.LoadYAML(resp.Body)

	case "xml":
		ds, err = tabular.LoadXML(resp.Body)

	case "csv":
		ds, err = tabular.LoadCSV(resp.Body)

	case "tsv":
		ds, err = tabular.LoadTSV(resp.Body)

	}

	if c.debugger != nil {
		c.debugger.Event(
			createEvent("tabular.dataset", resp.Request.ID, c.ID, map[string]string{
				"valid": fmt.Sprintf("%d", ds.Valid()),
				"cols":  fmt.Sprintf("%d", ds.Width()),
				"rows":  fmt.Sprintf("%d", ds.Height()),
			}))

		if err != nil {
			c.debugger.Event(
				createEvent("tabular.err", resp.Request.ID, c.ID, map[string]string{
					"err": err.Error(),
				}))
		}
	}

	// Invalid dataset if error
	// Note: We can check the dataset validity with `ds.Valid()` method.
	if err != nil {
		return err
	}

	// Note: It requires better query parsing to extract custom columns and rows selection
	for _, cc := range c.tabCallbacks {
		// var isRow, isSlice, isMixed bool
		if strings.Contains(cc.Query, ",") && strings.Contains(cc.Query, ":") {
			// isMixed = true
			return ErrTabularMixedSelectionNotImplemented
		}

		if strings.Contains(cc.Query, ",") {
			// rows, _ := ds.Rows(0, 1) // ([]map[string]interface{}, error) --> need to get Rows selection as a dataset struct
			// isRow = true
			return ErrTabularRowSelectionNotImplemented
		}

		var lowerRow, upperRow int
		var errQuery error
		if strings.Contains(cc.Query, ":") {
			lowerRow, upperRow, errQuery = parseSliceQuery(cc.Query, ds)
			// slice tabular dataset
			if errQuery == nil {
				ds, err = ds.Slice(lowerRow, upperRow)
			}
		}
		e := NewTABElementFromTABNode(resp, cc.Query, ds)

		if c.debugger != nil {
			c.debugger.Event(createEvent("tab", resp.Request.ID, c.ID, map[string]string{
				"selector": cc.Query,
				"url":      resp.Request.URL.String(),
			}))
		}
		cc.Function(e)
	}

	return nil
}

func (c *Collector) handleOnJSON(resp *Response) error {
	if len(c.jsonCallbacks) == 0 || !strings.Contains(strings.ToLower(resp.Headers.Get("Content-Type")), "json") {
		return nil
	}

	doc, err := jsonquery.Parse(bytes.NewBuffer(resp.Body))
	if err != nil {
		return err
	}

	for _, cc := range c.jsonCallbacks {
		for _, n := range jsonquery.Find(doc, cc.Query) {
			e := NewJSONElementFromJSONNode(resp, n)
			if c.debugger != nil {
				c.debugger.Event(createEvent("json", resp.Request.ID, c.ID, map[string]string{
					"selector": cc.Query,
					"url":      resp.Request.URL.String(),
				}))
			}
			cc.Function(e)
		}
	}

	return nil
}

func (c *Collector) handleOnHTML(resp *Response) error {
	if len(c.htmlCallbacks) == 0 || !strings.Contains(strings.ToLower(resp.Headers.Get("Content-Type")), "html") {
		return nil
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewBuffer(resp.Body))
	if err != nil {
		return err
	}
	if href, found := doc.Find("base[href]").Attr("href"); found {
		resp.Request.baseURL, _ = url.Parse(href)
	}
	for _, cc := range c.htmlCallbacks {
		doc.Find(cc.Selector).Each(func(i int, s *goquery.Selection) {
			for _, n := range s.Nodes {
				e := NewHTMLElementFromSelectionNode(resp, s, n)
				if c.debugger != nil {
					c.debugger.Event(createEvent("html", resp.Request.ID, c.ID, map[string]string{
						"selector": cc.Selector,
						"url":      resp.Request.URL.String(),
					}))
				}
				cc.Function(e)
			}
		})
	}
	return nil
}

func (c *Collector) handleOnXML(resp *Response) error {
	if len(c.xmlCallbacks) == 0 {
		return nil
	}
	contentType := strings.ToLower(resp.Headers.Get("Content-Type"))
	if !strings.Contains(contentType, "html") && !strings.Contains(contentType, "xml") {
		return nil
	}

	if strings.Contains(contentType, "html") {
		doc, err := htmlquery.Parse(bytes.NewBuffer(resp.Body))
		if err != nil {
			return err
		}
		if e := htmlquery.FindOne(doc, "//base/@href"); e != nil {
			for _, a := range e.Attr {
				if a.Key == "href" {
					resp.Request.baseURL, _ = url.Parse(a.Val)
					break
				}
			}
		}

		for _, cc := range c.xmlCallbacks {
			htmlquery.FindEach(doc, cc.Query, func(i int, n *html.Node) {
				e := NewXMLElementFromHTMLNode(resp, n)
				if c.debugger != nil {
					c.debugger.Event(createEvent("xml", resp.Request.ID, c.ID, map[string]string{
						"selector": cc.Query,
						"url":      resp.Request.URL.String(),
					}))
				}
				cc.Function(e)
			})
		}
	} else if strings.Contains(contentType, "xml") {
		doc, err := xmlquery.Parse(bytes.NewBuffer(resp.Body))
		if err != nil {
			return err
		}

		for _, cc := range c.xmlCallbacks {
			xmlquery.FindEach(doc, cc.Query, func(i int, n *xmlquery.Node) {
				e := NewXMLElementFromXMLNode(resp, n)
				if c.debugger != nil {
					c.debugger.Event(createEvent("xml", resp.Request.ID, c.ID, map[string]string{
						"selector": cc.Query,
						"url":      resp.Request.URL.String(),
					}))
				}
				cc.Function(e)
			})
		}
	}
	return nil
}

func (c *Collector) handleOnError(response *Response, err error, request *Request, ctx *Context) error {
	if err == nil && (c.ParseHTTPErrorResponse || response.StatusCode < 203) {
		return nil
	}
	if err == nil && response.StatusCode >= 203 {
		err = errors.New(http.StatusText(response.StatusCode))
	}
	if response == nil {
		response = &Response{
			Request: request,
			Ctx:     ctx,
		}
	}
	if c.debugger != nil {
		c.debugger.Event(createEvent("error", request.ID, c.ID, map[string]string{
			"url":    request.URL.String(),
			"status": http.StatusText(response.StatusCode),
		}))
	}
	if response.Request == nil {
		response.Request = request
	}
	if response.Ctx == nil {
		response.Ctx = request.Ctx
	}
	for _, f := range c.errorCallbacks {
		f(response, err)
	}
	return err
}

func (c *Collector) handleOnScraped(r *Response) {
	if c.debugger != nil {
		c.debugger.Event(createEvent("scraped", r.Request.ID, c.ID, map[string]string{
			"url": r.Request.URL.String(),
		}))
	}
	for _, f := range c.scrapedCallbacks {
		f(r)
	}
}

// Limit adds a new LimitRule to the collector
func (c *Collector) Limit(rule *LimitRule) error {
	return c.backend.Limit(rule)
}

// Limits adds new LimitRules to the collector
func (c *Collector) Limits(rules []*LimitRule) error {
	return c.backend.Limits(rules)
}

// SetCookies handles the receipt of the cookies in a reply for the given URL
func (c *Collector) SetCookies(URL string, cookies []*http.Cookie) error {
	if c.backend.Client.Jar == nil {
		return ErrNoCookieJar
	}
	u, err := url.Parse(URL)
	if err != nil {
		return err
	}
	c.backend.Client.Jar.SetCookies(u, cookies)
	return nil
}

// Cookies returns the cookies to send in a request for the given URL.
func (c *Collector) Cookies(URL string) []*http.Cookie {
	if c.backend.Client.Jar == nil {
		return nil
	}
	u, err := url.Parse(URL)
	if err != nil {
		return nil
	}
	return c.backend.Client.Jar.Cookies(u)
}

// Clone creates an exact copy of a Collector without callbacks.
// HTTP backend, robots.txt cache and cookie jar are shared
// between collectors.
func (c *Collector) Clone() (clone *Collector) {

	clone = &Collector{}
	clone.AllowedDomains = c.AllowedDomains
	clone.AllowURLRevisit = c.AllowURLRevisit
	clone.CacheDir = c.CacheDir
	clone.DetectCharset = c.DetectCharset
	clone.DisallowedDomains = c.DisallowedDomains
	clone.ID = atomic.AddUint32(&collectorCounter, 1)
	clone.IgnoreRobotsTxt = c.IgnoreRobotsTxt
	clone.MaxBodySize = c.MaxBodySize
	clone.MaxDepth = c.MaxDepth
	clone.DisallowedURLFilters = c.DisallowedURLFilters
	clone.URLFilters = c.URLFilters
	clone.ParseHTTPErrorResponse = c.ParseHTTPErrorResponse
	clone.UserAgent = c.UserAgent
	clone.Async = c.Async
	clone.DetectTabular = c.DetectTabular

	clone.store = c.store
	clone.backend = c.backend
	clone.debugger = c.debugger
	clone.RedirectHandler = c.RedirectHandler
	clone.errorCallbacks = make([]ErrorCallback, 0, 8)

	clone.htmlCallbacks = make([]*htmlCallbackContainer, 0, 8)
	clone.xmlCallbacks = make([]*xmlCallbackContainer, 0, 8)
	clone.scrapedCallbacks = make([]ScrapedCallback, 0, 8)

	clone.lock = c.lock
	clone.requestCallbacks = make([]RequestCallback, 0, 8)
	clone.responseCallbacks = make([]ResponseCallback, 0, 8)
	clone.robotsMap = c.robotsMap
	clone.wg = c.wg
	return

	/*
		return &Collector{
			store:             c.store,
			backend:           c.backend,
			debugger:          c.debugger,
			Async:             c.Async,
			RedirectHandler:   c.RedirectHandler,
			errorCallbacks:    make([]ErrorCallback, 0, 8),
			htmlCallbacks:     make([]*htmlCallbackContainer, 0, 8),
			xmlCallbacks:      make([]*xmlCallbackContainer, 0, 8),
			scrapedCallbacks:  make([]ScrapedCallback, 0, 8),
			lock:              c.lock,
			requestCallbacks:  make([]RequestCallback, 0, 8),
			responseCallbacks: make([]ResponseCallback, 0, 8),
			robotsMap:         c.robotsMap,
			wg:                c.wg,
		}
	*/
}

func (c *Collector) checkRedirectFunc() func(req *http.Request, via []*http.Request) error {
	return func(req *http.Request, via []*http.Request) error {
		if !c.isDomainAllowed(req.URL.Host) {
			return fmt.Errorf("Not following redirect to %s because its not in AllowedDomains", req.URL.Host)
		}

		if c.RedirectHandler != nil {
			return c.RedirectHandler(req, via)
		}

		// Honor golangs default of maximum of 10 redirects
		if len(via) >= 10 {
			return http.ErrUseLastResponse
		}

		lastRequest := via[len(via)-1]

		// Copy the headers from last request
		for hName, hValues := range lastRequest.Header {
			for _, hValue := range hValues {
				req.Header.Set(hName, hValue)
			}
		}

		// If domain has changed, remove the Authorization-header if it exists
		if req.URL.Host != lastRequest.URL.Host {
			req.Header.Del("Authorization")
		}

		return nil
	}
}

func (c *Collector) parseSettingsFromEnv() {
	for _, e := range os.Environ() {
		if !strings.HasPrefix(e, "COLLY_") {
			continue
		}
		pair := strings.SplitN(e[6:], "=", 2)
		if f, ok := envMap[pair[0]]; ok {
			f(c, pair[1])
		} else {
			log.Println("Unknown environment variable:", pair[0])
		}
	}
}

// SanitizeFileName replaces dangerous characters in a string
// so the return value can be used as a safe file name.
func SanitizeFileName(fileName string) string {
	ext := filepath.Ext(fileName)
	cleanExt := sanitize.BaseName(ext)
	if cleanExt == "" {
		cleanExt = ".unknown"
	}
	return strings.Replace(fmt.Sprintf(
		"%s.%s",
		sanitize.BaseName(fileName[:len(fileName)-len(ext)]),
		cleanExt[1:],
	), "-", "_", -1)
}

func createFormReader(data map[string]string) io.Reader {
	form := url.Values{}
	for k, v := range data {
		form.Add(k, v)
	}
	return strings.NewReader(form.Encode())
}

func createMultipartReader(boundary string, data map[string][]byte) io.Reader {
	dashBoundary := "--" + boundary

	body := []byte{}
	buffer := bytes.NewBuffer(body)

	buffer.WriteString("Content-type: multipart/form-data; boundary=" + boundary + "\n\n")
	for contentType, content := range data {
		buffer.WriteString(dashBoundary + "\n")
		buffer.WriteString("Content-Disposition: form-data; name=" + contentType + "\n")
		buffer.WriteString(fmt.Sprintf("Content-Length: %d \n\n", len(content)))
		buffer.Write(content)
		buffer.WriteString("\n")
	}
	buffer.WriteString(dashBoundary + "--\n\n")
	return buffer
}

// randomBoundary was borrowed from
// github.com/golang/go/mime/multipart/writer.go#randomBoundary
func randomBoundary() string {
	var buf [30]byte
	_, err := io.ReadFull(rand.Reader, buf[:])
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%x", buf[:])
}

func isYesString(s string) bool {
	switch strings.ToLower(s) {
	case "1", "yes", "true", "y":
		return true
	}
	return false
}

func createJar(s storage.Storage) http.CookieJar {
	return &cookieJarSerializer{store: s, lock: &sync.RWMutex{}}
}

func (j *cookieJarSerializer) SetCookies(u *url.URL, cookies []*http.Cookie) {
	j.lock.Lock()
	defer j.lock.Unlock()
	cookieStr := j.store.Cookies(u)

	// Merge existing cookies, new cookies have precedence.
	cnew := make([]*http.Cookie, len(cookies))
	copy(cnew, cookies)
	existing := storage.UnstringifyCookies(cookieStr)
	for _, c := range existing {
		if !storage.ContainsCookie(cnew, c.Name) {
			cnew = append(cnew, c)
		}
	}
	j.store.SetCookies(u, storage.StringifyCookies(cnew))
}

func (j *cookieJarSerializer) Cookies(u *url.URL) []*http.Cookie {
	cookies := storage.UnstringifyCookies(j.store.Cookies(u))
	// Filter.
	now := time.Now()
	cnew := make([]*http.Cookie, 0, len(cookies))
	for _, c := range cookies {
		// Drop expired cookies.
		if c.RawExpires != "" && c.Expires.Before(now) {
			continue
		}
		// Drop secure cookies if not over https.
		if c.Secure && u.Scheme != "https" {
			continue
		}
		cnew = append(cnew, c)
	}
	return cnew
}

func isMatchingFilter(fs []*regexp.Regexp, d []byte) bool {
	for _, r := range fs {
		if r.Match(d) {
			return true
		}
	}
	return false
}
