package filter

// Whitelist
type Whitelist struct {

	// Domains
	Domains []string `json:"domains" yaml:"domains" toml:"domains" xml:"domains" ini:"domains" csv:"Domains"`

	// URLs
	URLs []Config `json:"urls" yaml:"urls" toml:"urls" xml:"urls" ini:"urls" csv:"urls"`

	// FileExtensions
	FileExtensions []string `json:"file_extensions" yaml:"file_extensions" toml:"file_extensions" xml:"fileExtensions" ini:"fileExtensions" csv:"FileExtensions"`

	// Headers
	Headers []Config `json:"headers" yaml:"headers" toml:"headers" xml:"headers" ini:"headers" csv:"headers"`

	// MimeTypes
	MimeTypes []string `json:"mime_types" yaml:"mime_types" toml:"mime_types" xml:"mimeTypes" ini:"mimeTypes" csv:"MimeTypes"`

	// Responses
	Responses []Config `json:"responses" yaml:"responses" toml:"responses" xml:"responses" ini:"responses" csv:"responses"`

	//-- END
}
