package convert

import (
	"log"
	"strings"
	"sync"

	"github.com/sniperkit/colly/plugins/data/format/sitemap"
	"github.com/sniperkit/colly/plugins/data/transform/tabular"
)

type convert struct {
	// exported attributes
	// url           string `default:'sitemap'`
	// maxEntries    int    `default:100000`
	// forceEncoding string `default:'xml'` // available xml or txt
	*Config
	*stats

	// not exported

	datasets          map[string]*ds.Dataset
	databooks         map[string]*ds.Databooks
	sitemaps          []*sitemap.Sitemap
	contentSize       int
	compressionFormat string
	protocol          string
	content           string
	mimeType          string
	extension         string
	basename          string
	filename          string
	isCompressed      bool
	isLocalFile       bool // if false, expecting sitemap as []byte or string
	isValid           bool
	isDone            bool
	isPrepared        bool

	lock *sync.Mutex
	wg   *sync.WaitGroup
}

var entry struct {
	URL string `name:"url"`
}

func New(url string) *convert {
	c := &convert{}
	c.url = url
	return c
}

func NewWithConfig(cfg *Config) *convert {

	c := &convert{
		lock: &sync.RWMutex{},
		wg:   &sync.WaitGroup{},
	}
	c.Stats = &stats{extensions: make(map[string]int, 0)}

	// sets of sitemaps
	ssList := strings.Split(",", cfg.Datasets)
	if len(ssList) > 0 {
		c.datasets = make(map[string]*ds.Dataset, len(ssList))
	}

	// books of sitemaps
	bsList := strings.Split(",", cfg.Databooks)
	if len(bsList) > 0 {
		c.databooks = make(map[string]*ds.Dataset, len(bsList))
	}

	root_sitemap := &sitemap.Sitemap{
		ExportEntries: true,
		Enabled:       true,
	}

	c.sitemaps = append(c.sitemaps, sitemaps)

	return c
}

func (c *convert) LoadFileCSV(f string) {
	ds, err := tablib.LoadFileCSV(f)
	if err != nil {
		log.Fatal(err)
	}
}

func (c *convert) AppendDynamicColumns(cols []string) error {

	for _, col := range cols {
		ds.AppendDynamicColumn(col, DEFAULT_PRIORITY)
	}

	ds.AppendDynamicColumn("priority", DEFAULT_PRIORITY)
	// ds.AppendDynamicColumn("loc", DEFAULT_ENTRY_LOC)
	ds.AppendDynamicColumn("lastmod", DEFAULT_LAST_MOD)
	ds.AppendDynamicColumn("changefreq", DEFAULT_CHANGE_FREQ)

	xml, err := ds.XML()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(xml)

}

func (c *convert) TXT2XML(cols []string) error {

	ds.AppendDynamicColumn("priority", DEFAULT_PRIORITY)
	ds.AppendDynamicColumn("lastmod", DEFAULT_LAST_MOD)
	ds.AppendDynamicColumn("changefreq", DEFAULT_CHANGE_FREQ)

	xml, err := ds.XML()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(xml)

}
