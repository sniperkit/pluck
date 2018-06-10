package convert

/*
import (
	"log"
)

var logger *logger.Logger

*/

// Default allowed patterns
var (
	AllowedEncodingFormats []string = []string{"xml", "txt"}
	AllowedExtensions      []string = []string{"xml", "txt", "gz"}
	AllowedBasename        []string = []string{"sitemap_index", "sitemap"}
)

// Default Vars
var (
	DefaultSitemapFilename string   = DEFAULT_FILE_BASENAME + DEFAULT_FORMAT
	DefaultSitemapColumns  []string = []string{"priority", "loc", "lastmod", "changefreq"}
)

// Default Constants
const (
	// Default sitemap parameters
	DEFAULT_FORMAT                 string = "xml"
	DEFAULT_FILE_EXTENSION         string = "xml"
	DEFAULT_FILE_BASENAME          string = "sitemap"
	DEFAULT_LOCAL_FILEPATH         string = "./shared/storage/cache/sitemaps/sitemap" + DEFAULT_FILE_EXTENSION
	DEFAULT_STORAGE_TMP_DIR        string = "./shared/tmp"
	DEFAULT_STORAGE_EXPORT_DIR     string = "./shared/storage/cache/sitemaps"
	DEFAULT_STORAGE_EXPORT_PREFIX  string = "export_entries"
	DEFAULT_STORAGE_CONVERT_PREFIX string = "new_"

	// Default entry attributes
	DEFAULT_ENTRY_LOC   string  = ""
	DEFAULT_LAST_MOD    string  = "2005-01-01"
	DEFAULT_CHANGE_FREQ string  = "monthly"
	DEFAULT_MAX_ENTRIES int     = 100000
	DEFAULT_PRIORITY    float32 = 0.8
)

func urlLen(row []interface{}) interface{} {
	return len(row[0].(string))
}

func priority(row []interface{}) interface{} {
	return DEFAULT_PRIORITY
}

func loc(row []interface{}) interface{} {
	return row[0].(string)
}

func lastmod(row []interface{}) interface{} {
	return DEFAULT_LASTMOD
}

func changefreq(row []interface{}) interface{} {
	return DEFAULT_CHANGE_FREQ
}
