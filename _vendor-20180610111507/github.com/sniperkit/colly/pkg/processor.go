package colly

import (
	// striphtml "github.com/sniperkit/colly/plugins/data/encoding/html/striphtml"
	plucker "github.com/sniperkit/colly/plugins/data/extract/text/pluck"
	plucker_config "github.com/sniperkit/colly/plugins/data/extract/text/pluck/config"
)

// Extractor
type Extractor struct {
	ID      string `json:"id" yaml:"id" toml:"id" xml:"id" ini:"id"`
	plucker *plucker.Plucker
	// expression
}

type ExtractorConfig struct {
	units []*plucker_config.Config
}

/*
// NewExtractorWithConfig
func NewExtractorWithConfig(config *Extractor) (*Extractor, error) {
	p, err := plucker.New()
	return p, err
}

// NewExtractorWithConfigFile
func NewExtractorWithConfigFile(configFile string) (*Extractor, error) {
	p, err := plucker.New()
	return p, err
}

*/

// NewExtractor
// func NewExtractor() {

/*
	p.Add(pluck.Config{
		Name:        "",         // string - the key in the returned map, after completion
		Activators:  []string{}, // []string
		Deactivator: "",         // string
		Limit:       2,          // int - specifies the number of times capturing can occur
		Sanitize:    true,       // bool
		Finisher:    "",         // string - finishes capturing this pluck
		Permanent:   1,          // int - number of activators that stay permanently (counted from left to right)
	})
*/

// PluckFile takes a file as input
// and uses the specified parameters and generates
// a map (p.result) with the finished results. The streaming
// can be enabled by setting it to true.
//
// err = p.PluckFile(c.GlobalString("file"))

// PluckURL takes a URL as input
// and uses the specified parameters and generates
// a map (p.result) with the finished results
//
// err = p.PluckURL(c.GlobalString("url"))

// Pluck takes a buffered reader stream and
// extracts the text from it. This spawns a thread for
// each plucker and copies the entire buffer to memory,
// so that each plucker works in parallel.
//
// err = p.Pluck(r *bufio.Reader)

// PluckString takes a string as input
// and uses the specified parameters and generates
// a map (p.result) with the finished results.
// The streaming can be enabled by setting it to true.
//
// err = p.PluckString(s string, stream ...bool)

// PluckStream takes a buffered reader stream and streams one
// byte at a time and processes all pluckers serially and
// simultaneously.
//
// err = p.PluckStream(r *bufio.Reader)

// Result returns the raw result
// err = p.Result()

// }
