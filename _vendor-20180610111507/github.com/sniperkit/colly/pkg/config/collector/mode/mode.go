package collector

import (
	"time"

	mode_async "github.com/sniperkit/colly/pkg/config/collector/mode/async"
	mode_queue "github.com/sniperkit/colly/pkg/config/collector/mode/queue"
)

// Mode
type Mode struct {

	// Default
	Default struct {

		// RandomDelay
		RandomDelay time.Duration `default:"5" flag:"random-delay" yaml:"random_delay" toml:"random_delay" xml:"randomDelay" ini:"randomDelay" csv:"RandomDelay" json:"random_delay" yaml:"random_delay" toml:"random_delay" xml:"randomDelay" ini:"randomDelay" csv:"RandomDelay"`

		//-- END
	} `json:"default" yaml:"default" toml:"default" xml:"default" ini:"default" csv:"default"`

	// Async
	Async *mode_async.Config `json:"async" yaml:"async" toml:"async" xml:"async" ini:"async" csv:"async"`

	// Queue
	Queue *mode_queue.Config `json:"queue" yaml:"queue" toml:"queue" xml:"queue" ini:"queue" csv:"queue"`

	//-- END
}
