package html

/*
import (
	striphtml "github.com/sniperkit/colly/plugins/data/encoding/html/striphtml"
)
*/

// Config
type Config struct {

	// Enabled
	Enabled bool `default:"false" flag:"with-strip-html" yaml:"enabled" toml:"enabled" xml:"enabled" ini:"enabled" csv:"Enabled"`

	//-- END
}
