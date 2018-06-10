package collector

import (
	"time"
)

// Config
type Config struct {

	// Parallelism
	Parallelism int `default:"3" flag:"async-parallelism" yaml:"parallelism" toml:"parallelism" xml:"parallelism" ini:"parallelism" csv:"Parallelism" json:"parallelism" yaml:"parallelism" toml:"parallelism" xml:"parallelism" ini:"parallelism" csv:"Parallelism"`

	// DomainGlob
	DomainGlob string `default:"*" flag:"async-domain-glob" yaml:"domain_glob" toml:"domain_glob" xml:"domainGlob" ini:"domainGlob" csv:"DomainGlob" json:"domain_glob" yaml:"domain_glob" toml:"domain_glob" xml:"domainGlob" ini:"domainGlob" csv:"DomainGlob"`

	// RandomDelay
	RandomDelay time.Duration `default:"5" flag:"async-random-delay" yaml:"random_delay" toml:"random_delay" xml:"randomDelay" ini:"randomDelay" csv:"RandomDelay" json:"random_delay" yaml:"random_delay" toml:"random_delay" xml:"randomDelay" ini:"randomDelay" csv:"RandomDelay"`

	// MaxSize
	MaxSize int `default:"100000" flag:"async-max-size" yaml:"max_size" toml:"max_size" xml:"maxSize" ini:"maxSize" csv:"MaxSize" json:"max_size" yaml:"max_size" toml:"max_size" xml:"maxSize" ini:"maxSize" csv:"MaxSize"`

	//-- END
}
