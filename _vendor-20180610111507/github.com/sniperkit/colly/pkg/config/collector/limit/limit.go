package limit

import (
	"time"
)

// Config
type Config struct {

	// Parallelism
	Parallelism int `default:"2" flag:"limit-parallelism" yaml:"parallelism" toml:"parallelism" xml:"parallelism" ini:"parallelism" csv:"Parallelism"`

	// DomainGlob
	DomainGlob string `default:"*" flag:"limit-domain-glob" yaml:"domain_glob" toml:"domain_glob" xml:"domainGlob" ini:"domainGlob" csv:"DomainGlob"`

	// MaxQueue
	MaxQueue int `default:"10000" flag:"limit-max-queue" yaml:"max_queue" toml:"max_queue" xml:"maxQueue" ini:"maxQueue" csv:"MaxQueue"`

	// RandomDelay
	RandomDelay time.Duration `json:"random_delay" flag:"limit-random-delay" yaml:"random_delay" toml:"random_delay" xml:"randomDelay" ini:"randomDelay" csv:"RandomDelay"`

	//-- End
}
