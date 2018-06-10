package directory

// Config
type Config struct {

	// XDGBaseDir
	XDGBaseDir string `json:"xdg_base_dir" yaml:"xdg_base_dir" toml:"xdg_base_dir" xml:"xdgBaseDir" ini:"xdgBaseDir" csv:"XDGBaseDir"`

	// BaseDirectory
	BaseDir string `flag:"base-dir" json:"base_dir" yaml:"base_dir" toml:"base_dir" xml:"baseDir" ini:"baseDir" csv:"BaseDir"`

	// LogsDirectory
	LogsDir string `flag:"logs-dir" json:"logs_dir" yaml:"logs_dir" toml:"logs_dir" xml:"logsDir" ini:"logsDir" csv:"LogsDir"`

	// CacheDir specifies a location where GET requests are cached as files.
	// When it"s not defined, caching is disabled.
	CacheDir string `default:"./shared/storage/cache/http/backends/internal" flag:"cache-dir" yaml:"cache_dir" toml:"cache_dir" xml:"cacheDir" ini:"cacheDir" csv:"CacheDir" json:"cache_dir" yaml:"cache_dir" toml:"cache_dir" xml:"cacheDir" ini:"cacheDir" csv:"CacheDir"`

	// ExportDir
	ExportDir string `default:"./shared/exports" flag:"export-dir" yaml:"export_dir" toml:"export_dir" xml:"exportDir" ini:"exportDir" csv:"ExportDir" json:"export_dir" yaml:"export_dir" toml:"export_dir" xml:"exportDir" ini:"exportDir" csv:"ExportDir"`

	// ForceDir specifies that the program will try to create missing storage directories.
	ForceDir bool `default:"true" flag:"force-dir" yaml:"force_dir" toml:"force_dir" xml:"forceDir" ini:"forceDir" csv:"ForceDir" json:"force_dir" yaml:"force_dir" toml:"force_dir" xml:"forceDir" ini:"forceDir" csv:"ForceDir"`

	// ForceDirRecursive specifies that the program will try to create missing storage directories recursively.
	ForceDirRecursive bool `default:"true" flag:"force-dir-recursive" yaml:"force_dir_recursive" toml:"force_dir_recursive" xml:"forceDirRecursive" ini:"forceDirRecursive" csv:"ForceDirRecursive" json:"force_dir_recursive" yaml:"force_dir_recursive" toml:"force_dir_recursive" xml:"forceDirRecursive" ini:"forceDirRecursive" csv:"ForceDirRecursive"`

	//-- END
}
