package config

import (
	"fmt"
	"os"
	// pp "github.com/sniperkit/colly/plugins/app/debug/pp"
)

const (
	// DEFAULT_BASE_DIR
	DEFAULT_BASE_DIR = "."
)

var (
	// DefaultXDGBaseDirectory
	DefaultXDGBaseDirectory string = ""

	// DefaultConfigBasenameList
	DefaultConfigBasenameList []string = []string{"config", "pluck", "plucker"}

	// DefaultConfigFormatList
	DefaultConfigFormatList []string = []string{"yaml", "yml", "toml", "ini", "xml", "json"}

	// DefaultConfigPrefixPathList
	DefaultConfigPrefixPathList []string = []string{"conf", "plucker", "pluck", ".plucker", ".pluck"}

	// DefaultConfigFilepaths
	DefaultConfigFilepaths []string = []string{"conf/config.toml", "~/.plucker/config", "~/.plucker/config.yaml"}

	//-- Ends
)

func init() {
	var err error
	DefaultXDGBaseDirectory, err = getDefaultXDGBaseDirectory()
	if err != nil {
		fmt.Println("error while trying to get the default xdgb base directory")
		os.Exit(1)
	}
}

// GenerateExpectedFilepaths
func GenerateExpectedFilepaths(pp string) {
	var defaultConfigPrefixPathList []string

	defaultConfigPrefixPathList = append(defaultConfigPrefixPathList, DefaultConfigPrefixPathList...)
	for _, p := range DefaultConfigPrefixPathList {
		newPath := fmt.Sprintf("%s/%s", DefaultXDGBaseDirectory, p)
		for _, b := range DefaultConfigBasenameList {
			newPath := fmt.Sprintf("%s/%s", newPath, b)
			for _, f := range DefaultConfigFormatList {
				newPath := fmt.Sprintf("%s.%s", newPath, f)
				defaultConfigPrefixPathList = append(defaultConfigPrefixPathList, newPath)
			}
		}
	}
	if len(defaultConfigPrefixPathList) > 0 {
		DefaultConfigFilepaths = defaultConfigPrefixPathList
	}
	// pp.Println(defaultConfigPrefixPathList)
}

// AddConfigPaths
func AddConfigPaths(pp ...string) {
	for _, p := range pp {
		DefaultConfigPrefixPathList = append(DefaultConfigPrefixPathList, p)
	}
}
