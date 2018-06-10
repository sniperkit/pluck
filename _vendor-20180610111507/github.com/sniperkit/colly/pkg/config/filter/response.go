package filter

// Response
type Response struct {

	// ParseHTTPErrorResponse allows parsing HTTP responses with non 2xx status codes.
	// By default, Colly parses only successful HTTP responses. Set ParseHTTPErrorResponse to true to enable it.
	ParseHTTPErrorResponse bool `default:"true" flag:"parse-http-error-response" yaml:"parse_http_error_response" toml:"parse_http_error_response" xml:"parseHTTPErrorResponse" ini:"parseHTTPErrorResponse" csv:"parseHTTPErrorResponse" json:"parse_http_error_response" yaml:"parse_http_error_response" toml:"parse_http_error_response" xml:"parseHTTPErrorResponse" ini:"parseHTTPErrorResponse" csv:"parseHTTPErrorResponse"`

	// DetectCharset can enable character encoding detection for non-utf8 response bodies
	// without explicit charset declaration. This feature uses https://github.com/saintfish/chardet
	DetectCharset bool `default:"true" flag:"detect-charset" yaml:"detect_charset" toml:"detect_charset" xml:"detectCharset" ini:"detectCharset" csv:"DetectCharset" json:"detect_charset" yaml:"detect_charset" toml:"detect_charset" xml:"detectCharset" ini:"detectCharset" csv:"DetectCharset"`

	// DetectMimeType
	DetectMimeType bool `default:"true" flag:"detect-mime-type" yaml:"detect_mime_type" toml:"detect_mime_type" xml:"detectMimeType" ini:"detectMimeType" csv:"detectMimeType" json:"detect_mime_type" yaml:"detect_mime_type" toml:"detect_mime_type" xml:"detectMimeType" ini:"detectMimeType" csv:"detectMimeType"`

	// DetectTabular
	DetectTabular bool `default:"true" flag:"detect-tabular-data" yaml:"detect_tabular_data" toml:"detect_tabular_data" xml:"detectTabularData" ini:"detectTabularData" csv:"DetectTabularData" json:"detect_tabular_data" yaml:"detect_tabular_data" toml:"detect_tabular_data" xml:"detectTabularData" ini:"detectTabularData" csv:"DetectTabularData"`

	// MaxBodySize is the limit of the retrieved response body in bytes.
	// 0 means unlimited.
	// The default value for MaxBodySize is 10MB (10 * 1024 * 1024 bytes).
	MaxBodySize int `default:"0" flag:"max-body-size" yaml:"max_body_size" toml:"max_body_size" xml:"maxBodySize" ini:"maxBodySize" csv:"maxBodySize" json:"max_body_size" yaml:"max_body_size" toml:"max_body_size" xml:"maxBodySize" ini:"maxBodySize" csv:"maxBodySize"`

	//-- END
}
