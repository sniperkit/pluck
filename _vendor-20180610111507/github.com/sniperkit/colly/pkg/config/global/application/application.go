package application

// Config
type Config struct {

	// ID is the unique identifier of a collector
	ID uint32 `default:"colly" flag:"identifier" yaml:"identifier" toml:"identifier" xml:"identifier" ini:"identifier" csv:"identifier"`

	// Title/name of the current crawling campaign
	Title string `default:"Colly - Web Scraper" flag:"title" yaml:"title" toml:"title" xml:"title" ini:"title" csv:"title"`

	// DebugMode
	DebugMode bool `default:"false" flag:"debug" yaml:"debug" toml:"debug" xml:"debugMode" ini:"debugMode" csv:"DebugMode"`

	// VerboseMode
	VerboseMode bool `default:"false" flag:"verbose" yaml:"verbose" toml:"verbose" xml:"verboseMode" ini:"verboseMode" csv:"VerboseMode"`

	// IsDashboard
	DashboardMode bool `default:"true" flag:"dashboard" yaml:"dashboard" toml:"dashboard" xml:"dashboard" ini:"dashboardMode" csv:"dashboardMode"`

	//-- END
}
