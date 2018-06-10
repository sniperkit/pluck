package metric

import (
	"sync"
)

type MetricCollector struct {
	Pause       bool `default:"true"`
	DebugMode   bool `default:"true"`
	VerboseMode bool `default:"true"`
	Payload     string
	wg          *sync.WaitGroup
	lock        *sync.RWMutex
}

// NewMetricCollector creates a new Collector instance with cfg.Default configuration
func NewMetricCollector(options ...func(*MetricCollector)) *MetricCollector {
	mc := &MetricCollector{}
	mc.Init()
	for _, f := range options {
		f(mc)
	}
	// mc.parseSettingsFromEnv()
	return mc
}

// Init initializes the MetricCollector's private variables and sets default configuration for the MetricCollector
func (c *MetricCollector) Init() {
	c.wg = &sync.WaitGroup{}
	c.lock = &sync.RWMutex{}
	c.Pause = true
}

// SetPayload enables ...
func SetPayload(payload string) func(*MetricCollector) {
	return func(c *MetricCollector) {
		c.Payload = payload
	}
}

// Pause enables ...
func Pause() func(*MetricCollector) {
	return func(c *MetricCollector) {
		c.Pause = true
	}
}

func Resume() func(*MetricCollector) {
	return func(c *MetricCollector) {
		c.Pause = false
	}
}

// DebugMode enables ...
func DebugMode(status bool) func(*MetricCollector) {
	return func(c *MetricCollector) {
		c.DebugMode = status
	}
}

// VerboseMode enables ...
func VerboseMode(status bool) func(*MetricCollector) {
	return func(c *MetricCollector) {
		c.VerboseMode = status
	}
}
