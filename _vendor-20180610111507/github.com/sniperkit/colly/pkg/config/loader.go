package config

import (
	"fmt"
	"os"

	configor "github.com/sniperkit/colly/plugins/data/import/configor"
)

// Snippets:
// - configor.New(&configor.Config{Debug: true, Verbose: true}).Load(&Config, "config.json")

// autoLoad
func autoLoad() {
	var err error
	collectorBaseDir, err = configor.XDGBaseDir()
	if err != nil {
		fmt.Println("Can't find XDG BaseDirectory")
		os.Exit(1)
	}
}

// NewFromFile
func NewFromFile(verbose, debug, esrrorOnUnmatchedKeys bool, files ...string) (*Config, error) {
	globalConfig := &Config{}
	xdgPath, err := getDefaultXDGBaseDirectory()
	if err != nil {
		return nil, err
	}
	globalConfig.Dirs.XDGBaseDir = xdgPath
	configor.New(&configor.Config{Debug: debug, Verbose: verbose, ErrorOnUnmatchedKeys: false}).Load(&globalConfig, files...)

	return globalConfig, nil
}

// Dump
func (c *Config) Dump(formats, nodes []string, prefixPath string) error {
	return configor.Dump(c, nodes, formats, prefixPath)
}

// Dump
func Dump(c interface{}, formats, nodes []string, prefixPath string) error {
	return configor.Dump(c, nodes, formats, prefixPath)
}

func getDefaultXDGBaseDirectory() (string, error) {
	xdgPath, err := configor.XDGBaseDir()
	if err != nil {
		return DEFAULT_BASE_DIR, err
	}
	return xdgPath, nil
}
