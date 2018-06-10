package sitemap

const (
	DefaultPrintFormat string = "tabular"
)

var (
	AvailablePrintFormats []string = []string{"tabular", "yaml", "json", "xml", "toml", "mysql", "postgres", "csv", "tsv", "markdown"}
)
