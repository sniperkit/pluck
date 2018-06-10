package debug

// Config
type Config struct {

	// LoadVerbose
	LoadVerbose bool `default:"false" flag:"config-verbose" yaml:"load_verbose" toml:"load_verbose" xml:"loadVerbose" ini:"loadVerbose" csv:"LoadVerbose" json:"load_verbose" yaml:"load_verbose" toml:"load_verbose" xml:"loadVerbose" ini:"loadVerbose" csv:"LoadVerbose"`

	// LoadDebug
	LoadDebug bool `default:"false" flag:"config-debug" yaml:"load_debug" toml:"load_debug" xml:"loadDebug" ini:"loadDebug" csv:"LoadDebug" json:"load_debug" yaml:"load_debug" toml:"load_debug" xml:"loadDebug" ini:"loadDebug" csv:"LoadDebug"`

	// LoadErrorOnUnmatchedKeys
	LoadErrorOnUnmatchedKeys bool `default:"false" flag:"with-error-unmatched-keys" yaml:"load_error_on_unmatched_keys" toml:"load_error_on_unmatched_keys" xml:"loadErrorOnUnmatchedKeys" ini:"loadErrorOnUnmatchedKeys" csv:"LoadErrorOnUnmatchedKeys" json:"load_error_on_unmatched_keys" yaml:"load_error_on_unmatched_keys" toml:"load_error_on_unmatched_keys" xml:"loadErrorOnUnmatchedKeys" ini:"loadErrorOnUnmatchedKeys" csv:"LoadErrorOnUnmatchedKeys"`

	// ExportDisabled
	ExportEnabled bool `default:"true" flag:"config-export" yaml:"export_enabled" toml:"export_enabled" xml:"exportEnabled" ini:"exportEnabled" csv:"ExportEnabled" json:"export_enabled" yaml:"export_enabled" toml:"export_enabled" xml:"exportEnabled" ini:"exportEnabled" csv:"ExportEnabled"`

	// ExportSections
	ExportSections []string `json:"export_sections" yaml:"export_sections" toml:"export_sections" xml:"ExportSections" ini:"ExportSections" csv:"ExportSections"`

	// ExportSchemaOnly
	ExportSchemaOnly bool `default:"false" flag:"config-schema-only" yaml:"export_schema_only" toml:"export_schema_only" xml:"exportSchemaOnly" ini:"exportSchemaOnly" csv:"ExportSchemaOnly" json:"export_schema_only" yaml:"export_schema_only" toml:"export_schema_only" xml:"exportSchemaOnly" ini:"exportSchemaOnly" csv:"ExportSchemaOnly"`

	// ExportPrefixPath
	ExportPrefixPath string `default:"./shared/exports/config/dump" flag:"config-export-prefix-path" yaml:"export_prefix_path" toml:"export_prefix_path" xml:"exportPrefixPath" ini:"exportPrefixPath" csv:"ExportPrefixPath" json:"export_prefix_path" yaml:"export_prefix_path" toml:"export_prefix_path" xml:"exportPrefixPath" ini:"exportPrefixPath" csv:"ExportPrefixPath"`

	// ExportFormat
	ExportFormat []string `json:"export_formats" yaml:"export_formats" toml:"export_formats" xml:"exportFormats" ini:"exportFormats" csv:"ExportFormats"`

	//-- END
}
