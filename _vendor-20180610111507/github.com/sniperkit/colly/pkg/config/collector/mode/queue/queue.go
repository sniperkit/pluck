package collector

import (
	"time"
)

// Config
type Config struct {

	// Workers
	WorkersCount int `default:"3" flag:"queue-workers" yaml:"workers_count" toml:"workers_count" xml:"workersCount" ini:"workersCount" csv:"WorkersCount" json:"workers_count" yaml:"workers_count" toml:"workers_count" xml:"workersCount" ini:"workersCount" csv:"WorkersCount"`

	// MaxSize
	MaxSize int `default:"100000" flag:"queue-max-size" yaml:"max_size" toml:"max_size" xml:"maxSize" ini:"maxSize" csv:"MaxSize" json:"max_size" yaml:"max_size" toml:"max_size" xml:"maxSize" ini:"maxSize" csv:"MaxSize"`

	// RandomDelay
	RandomDelay time.Duration `default:"5" flag:"queue-random-delay" yaml:"random_delay" toml:"random_delay" xml:"randomDelay" ini:"randomDelay" csv:"RandomDelay" json:"random_delay" yaml:"random_delay" toml:"random_delay" xml:"randomDelay" ini:"randomDelay" csv:"RandomDelay"`

	//-- END
}
