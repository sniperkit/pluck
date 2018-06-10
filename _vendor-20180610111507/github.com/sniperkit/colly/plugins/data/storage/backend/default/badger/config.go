package badgercache

import (
	"github.com/jinzhu/configor"

	"github.com/sniperkit/vipertags"
	"github.com/sniperkit/xconfig"
	// cache "github.com/sniperkit/xcache/pkg"

	"github.com/sniperkit/xutil/plugin/debug/pp"
)

type badgercacheConfig struct {
	Provider       string        `json:"provider" config:"database.provider"`
	Endpoints      []string      `json:"endpoints" config:"database.endpoints"`
	MaxConnections int           `json:"max_connections" config:"database.max_connections" default:"0"`
	done           chan struct{} `json:"-" config:"-"`
}

// Config ...
var (
	PluginConfig = &badgercacheConfig{
		done: make(chan struct{}),
	}
)

// ConfigName ...
func (badgercacheConfig) ConfigName() string {
	return "BadgerKV"
}

// SetDefaults ...
func (a *badgercacheConfig) SetDefaults() {
	vipertags.SetDefaults(a)
}

// Read ...
func (a *badgercacheConfig) Read() {
	defer close(a.done)
	vipertags.Fill(a)
	if a.Provider == "" {
		a.Provider = a.ConfigName()
	}
	//if a.MaxConnections == 0 {
	//	a.MaxConnections = cache.DefaultMaxConnections
	//}
}

// Read several config files (yaml, json or env variables)
func (a *badgercacheConfig) Configor(files []string) {
	configor.Load(&PluginConfig, files...)
}

// Wait ...
func (c badgercacheConfig) Wait() {
	<-c.done
}

// String ...
func (c badgercacheConfig) String() string {
	return pp.Sprintln(c)
}

// Debug ...
func (c badgercacheConfig) Debug() {
	// log.Debug("BadgerKV PluginConfig = ", c)
}

func init() {
	config.Register(PluginConfig)
}
