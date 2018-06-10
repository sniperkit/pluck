package config

import (
	"fmt"
	"os"

	// config
	// "github.com/BurntSushi/toml"
	// "github.com/jinzhu/configor"

	configor "github.com/sniperkit/colly/plugins/data/import/configor"
)

/*
 - Snippets:
   - configor.New(&configor.Config{Debug: true, Verbose: true}).Load(&Config, "config.json")
*/

var (
	// AutoLoad
	AutoLoad bool = false
)

func init() {
	if AutoLoad {
		autoLoad()
	}
}

// autoLoad
func autoLoad() {
	var err error
	DefaultXDGBaseDirectory, err = configor.XDGBaseDir()
	if err != nil {
		fmt.Println("Can't find XDG BaseDirectory")
		os.Exit(1)
	}
}

// NewFromFile
func NewFromFile(verbose, debug, esrrorOnUnmatchedKeys bool, files ...string) (*Configs, error) {
	globalConfig := &Configs{}
	xdgPath, err := getDefaultXDGBaseDirectory()
	if err != nil {
		return nil, err
	}
	globalConfig.XDGBaseDir = xdgPath
	configor.New(&configor.Config{Debug: debug, Verbose: verbose, ErrorOnUnmatchedKeys: false}).Load(&globalConfig, files...)

	return globalConfig, nil
}

// Dump
func (c *Configs) Dump(formats, nodes []string, prefixPath string) error {
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
