package config

import (
	"time"

	// global config
	global_app "github.com/sniperkit/colly/pkg/config/global/application"
	global_dir "github.com/sniperkit/colly/pkg/config/global/directory"

	// application helpers
	debug "github.com/sniperkit/colly/pkg/config/debug"
	debug_tachymeter "github.com/sniperkit/colly/pkg/config/debug/tachymeter"

	// collector - core
	collector "github.com/sniperkit/colly/pkg/config/collector"

	// collector - helpers
	filter "github.com/sniperkit/colly/pkg/config/filter"

	// collector - data storage
	store_collection "github.com/sniperkit/colly/pkg/config/store/collection"
)

// public vars
var (
	// AutoLoad enables to load the default collector
	AutoLoad = false
)

// private vars
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

// init
func init() {
	if AutoLoad {
		autoLoad()
	}
}

// Config
type Config struct {

	// createdAt is set when...
	createdAt time.Time

	// startedAt is set when...
	startedAt time.Time

	// App
	App global_app.Config `json:"app" yaml:"app" toml:"app" xml:"app" ini:"app" csv:"App"`

	// Debug
	Debug struct {

		// Config
		Config debug.Config `json:"config" yaml:"config" toml:"config" xml:"config" ini:"config" csv:"Config"`

		// Tachymeter
		Tachymeter debug_tachymeter.Config `json:"tachymeter" yaml:"tachymeter" toml:"tachymeter" xml:"tachymeter" ini:"tachymeter" csv:"tachymeter"`

		//-- END
	} `json:"debug" yaml:"debug" toml:"debug" xml:"debug" ini:"debug" csv:"Debug"`

	// Collector
	Collector collector.Config `json:"collector" yaml:"collector" toml:"collector" xml:"collector" ini:"collector" csv:"collector"`

	// Collectors (create several collectors... eg master and slaves...)
	// Collectors []*collector.Config `json:"collector" yaml:"collector" toml:"collector" xml:"collector" ini:"collector" csv:"collector"`

	// Filters
	Filters *filter.Filter `json:"filters" yaml:"filters" toml:"filters" xml:"filters" ini:"filters" csv:"filters"`

	// Collection
	Collection *store_collection.Config `json:"collection" yaml:"collection" toml:"collection" xml:"collection" ini:"collection" csv:"collection"`

	// Dirs
	Dirs global_dir.Config `json:"outputs" yaml:"outputs" toml:"outputs" xml:"outputs" ini:"outputs" csv:"outputs"`
}
