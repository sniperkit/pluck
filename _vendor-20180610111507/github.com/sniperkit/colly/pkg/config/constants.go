package config

const (

	// DEFAULT_APP_NAME
	DEFAULT_APP_NAME string = "X-Colly - Web Crawler"

	// DEFAULT_APP_VERSION
	DEFAULT_APP_VERSION string = "1.0.0"

	// DEFAULT_BASE_DIR
	DEFAULT_BASE_DIR string = "./shared"

	// DefaultDetectMimeType
	DefaultDetectMimeType bool = true

	// DefaultDetectCharset
	DefaultDetectCharset bool = true

	// DefaultParseHTTPErrorResponse
	DefaultParseHTTPErrorResponse bool = true

	// DefaultForceDir
	DefaultForceDir bool = true

	// DefaultForceDirRecursive
	DefaultForceDirRecursive bool = true

	// DefaultIgnoreRobotsTxt
	DefaultIgnoreRobotsTxt bool = true

	// DefaultAllowURLRevisit
	DefaultAllowURLRevisit bool = false

	// DefaultRandomUserAgent
	DefaultRandomUserAgent bool = false

	// DefaultSummarizeContent
	DefaultSummarizeContent bool = false

	// DefaultTopicModelling
	DefaultTopicModelling bool = false

	// DefaultAnalyzeContent
	DefaultAnalyzeContent bool = false

	// DefaultDebugMode
	DefaultDebugMode bool = false

	// DefaultVerboseMode
	DefaultVerboseMode bool = false

	// DefaultMaxDepth
	DefaultMaxDepth int = 0

	// DefaultMaxBodySize
	DefaultMaxBodySize int = 10 * 1024 * 1024

	// DefaultConfigFilepath
	DefaultConfigFilepath string = "./colly.yaml"

	// DefaultSitemapXpath
	DefaultSitemapXpath string = "//urlset/url/loc"

	// DefaultSitemapFilename
	DefaultSitemapFilename string = "sitemap.xml"

	// DefaultUserAgent
	DefaultUserAgent string = "X-Colly - Alpha"

	// DefaultStorageDir
	DefaultStorageDir string = "./shared/storage"

	// DefaultCacheDir
	DefaultCacheDir string = "http/raw"

	// DefaultEnvPrefix
	DefaultEnvPrefix string = "COLLY_"
)
