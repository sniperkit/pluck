package convert

type Config struct {
	CacheDir  string `default:'./shared/storage/cache/sitemaps'`
	ExportDir string `default:'./shared/storage/export/sitemaps'`

	Datasets  string `default:'default_dataset'`
	Databooks string `default:''`

	AddDynamicColumns string `default:'priority,changefreq,lastmod'`
	AddColumns        string `default:''`
	DeleteColumns     string `default:''`
	DeleteRows        string `default:''`
	SelectRows        string `default:''`
	ForceEncoding     string `default:'xml'` // available xml or txt
	DetectEncoding    bool   `default:'true'`
	MaxEntries        int    `default:100000`
	createdAt         time.Time
}

type DatasetConfig struct {
	Headers string `default:'column_1'`
	MaxCols int    `default:'50'`
	MaxRows int    `default:'2500'`
	SplitAt int    `default:'2500'`
	headers []string
}

type DatabookConfig struct {
	Headers string `default:'column_1'`
	MaxCols int    `default:'50'`
	MaxRows int    `default:'2500'`
	SplitAt int    `default:'2500'`
	headers []string
}
