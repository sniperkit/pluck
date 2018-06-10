package export

// Config
type Config struct {

	// Enabled
	Enabled bool `default:"false" json:"enabled" yaml:"enabled" toml:"enabled" xml:"enabled" ini:"enabled"`

	// BackupMode
	BackupMode bool `default:"true" json:"backup_mode" yaml:"backup_mode" toml:"backupMode" xml:"backupMode" ini:"BackupMode"`

	// BackupSuffix
	BackupSuffix string `default:"_backupSuffix" json:"backup_suffix" yaml:"backup_suffix" toml:"backupSuffix" xml:"backupSuffix" ini:"BackupSuffix"`

	// BackupPrefix
	BackupPrefix string `default:"backupPrefix_" json:"backup_prefix" yaml:"backup_prefix" toml:"backupPrefix" xml:"backupPrefix" ini:"BackupPrefix"`

	// Overwrite
	Overwrite bool `default:"false" json:"overwrite" yaml:"overwrite" toml:"overwrite" xml:"overwrite" ini:"overwrite"`

	// Encoding sets..
	Encoding string `default:"csv" json:"encoding" yaml:"encoding" toml:"encoding" xml:"encoding" ini:"encoding"`

	// Format sets..
	UseTemplate string `json:"format" yaml:"format" toml:"format" xml:"format" ini:"format"`

	// Basename
	Basename string `default:"export_%s_%d" json:"basename" yaml:"basename" toml:"basename" xml:"basename" ini:"basename"`

	// SplitAt
	SplitAt int `json:"split_at" yaml:"split_at" toml:"split_at" xml:"splitAt" ini:"splitAt"`

	// BufferSize
	BufferSize int `json:"buffer_size" yaml:"buffer_size" toml:"buffer_size" xml:"bufferSize" ini:"bufferSize"`

	// PrefixPath sets..
	PrefixPath string `default:"./shared" json:"prefix_path" yaml:"prefix_path" toml:"prefix_path" xml:"prefixPath" ini:"prefixPath"`

	// ExportDir
	ExportDir string `default:"./storage/export" json:"export_dir" yaml:"export_dir" toml:"export_dir" xml:"exportDir" ini:"exportDir"`

	// ForceDir specifies that the program will try to create missing storage directories.
	EnsureDirs bool `default:"true" json:"ensure_dir" yaml:"ensure_dir" toml:"ensure_dir" xml:"ensureDirs" ini:"ensureDirs"`

	// ForceDirRecursive specifies that the program will try to create missing storage directories recursively.
	ForceDirRecursive bool `default:"true" json:"ensure_dir_recursively" yaml:"ensure_dir_recursively" toml:"ensure_dir_recursively" xml:"ensureDirRecursively" ini:"ensureDirRecursively"`

	//- End
}
