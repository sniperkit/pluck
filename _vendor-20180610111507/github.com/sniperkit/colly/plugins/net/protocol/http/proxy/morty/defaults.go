package morty

import (
	"regexp"

	"github.com/sniperkit/colly/plugins/data/encoding/html/contenttype"
)

var (
	CSS_URL_REGEXP *regexp.Regexp = regexp.MustCompile("url\\((['\"]?)[ \\t\\f]*([\u0009\u0021\u0023-\u0026\u0028\u002a-\u007E]+)(['\"]?)\\)?")

	ALLOWED_CONTENTTYPE_FILTER contenttype.Filter = contenttype.NewFilterOr([]contenttype.Filter{
		// html
		contenttype.NewFilterEquals("text", "html", ""),
		contenttype.NewFilterEquals("application", "xhtml", "xml"),
		// css
		contenttype.NewFilterEquals("text", "css", ""),
		// images
		contenttype.NewFilterEquals("image", "gif", ""),
		contenttype.NewFilterEquals("image", "png", ""),
		contenttype.NewFilterEquals("image", "jpeg", ""),
		contenttype.NewFilterEquals("image", "pjpeg", ""),
		contenttype.NewFilterEquals("image", "webp", ""),
		contenttype.NewFilterEquals("image", "tiff", ""),
		contenttype.NewFilterEquals("image", "vnd.microsoft.icon", ""),
		contenttype.NewFilterEquals("image", "bmp", ""),
		contenttype.NewFilterEquals("image", "x-ms-bmp", ""),
		// fonts
		contenttype.NewFilterEquals("application", "font-otf", ""),
		contenttype.NewFilterEquals("application", "font-ttf", ""),
		contenttype.NewFilterEquals("application", "font-woff", ""),
		contenttype.NewFilterEquals("application", "vnd.ms-fontobject", ""),
	})

	ALLOWED_CONTENTTYPE_ATTACHMENT_FILTER contenttype.Filter = contenttype.NewFilterOr([]contenttype.Filter{
		// texts
		contenttype.NewFilterEquals("text", "csv", ""),
		contenttype.NewFilterEquals("text", "tab-separated-value", ""),
		contenttype.NewFilterEquals("text", "plain", ""),
		// API
		contenttype.NewFilterEquals("application", "json", ""),
		// Documents
		contenttype.NewFilterEquals("application", "x-latex", ""),
		contenttype.NewFilterEquals("application", "pdf", ""),
		contenttype.NewFilterEquals("application", "vnd.oasis.opendocument.text", ""),
		contenttype.NewFilterEquals("application", "vnd.oasis.opendocument.spreadsheet", ""),
		contenttype.NewFilterEquals("application", "vnd.oasis.opendocument.presentation", ""),
		contenttype.NewFilterEquals("application", "vnd.oasis.opendocument.graphics", ""),
		// Compressed archives
		contenttype.NewFilterEquals("application", "zip", ""),
		contenttype.NewFilterEquals("application", "gzip", ""),
		contenttype.NewFilterEquals("application", "x-compressed", ""),
		contenttype.NewFilterEquals("application", "x-gtar", ""),
		contenttype.NewFilterEquals("application", "x-compress", ""),
		// Generic binary
		contenttype.NewFilterEquals("application", "octet-stream", ""),
	})

	ALLOWED_CONTENTTYPE_PARAMETERS map[string]bool = map[string]bool{
		"charset": true,
	}

	UNSAFE_ELEMENTS [][]byte = [][]byte{
		[]byte("applet"),
		[]byte("canvas"),
		[]byte("embed"),
		//[]byte("iframe"),
		[]byte("math"),
		[]byte("script"),
		[]byte("svg"),
	}

	SAFE_ATTRIBUTES [][]byte = [][]byte{
		[]byte("abbr"),
		[]byte("accesskey"),
		[]byte("align"),
		[]byte("alt"),
		[]byte("as"),
		[]byte("autocomplete"),
		[]byte("charset"),
		[]byte("checked"),
		[]byte("class"),
		[]byte("content"),
		[]byte("contenteditable"),
		[]byte("contextmenu"),
		[]byte("dir"),
		[]byte("for"),
		[]byte("height"),
		[]byte("hidden"),
		[]byte("hreflang"),
		[]byte("id"),
		[]byte("lang"),
		[]byte("media"),
		[]byte("method"),
		[]byte("name"),
		[]byte("nowrap"),
		[]byte("placeholder"),
		[]byte("property"),
		[]byte("rel"),
		[]byte("spellcheck"),
		[]byte("tabindex"),
		[]byte("target"),
		[]byte("title"),
		[]byte("translate"),
		[]byte("type"),
		[]byte("value"),
		[]byte("width"),
	}

	SELF_CLOSING_ELEMENTS [][]byte = [][]byte{
		[]byte("area"),
		[]byte("base"),
		[]byte("br"),
		[]byte("col"),
		[]byte("embed"),
		[]byte("hr"),
		[]byte("img"),
		[]byte("input"),
		[]byte("keygen"),
		[]byte("link"),
		[]byte("meta"),
		[]byte("param"),
		[]byte("source"),
		[]byte("track"),
		[]byte("wbr"),
	}

	LINK_REL_SAFE_VALUES [][]byte = [][]byte{
		[]byte("alternate"),
		[]byte("archives"),
		[]byte("author"),
		[]byte("copyright"),
		[]byte("first"),
		[]byte("help"),
		[]byte("icon"),
		[]byte("index"),
		[]byte("last"),
		[]byte("license"),
		[]byte("manifest"),
		[]byte("next"),
		[]byte("pingback"),
		[]byte("prev"),
		[]byte("publisher"),
		[]byte("search"),
		[]byte("shortcut icon"),
		[]byte("stylesheet"),
		[]byte("up"),
	}

	LINK_HTTP_EQUIV_SAFE_VALUES [][]byte = [][]byte{
		// X-UA-Compatible will be added automaticaly, so it can be skipped
		[]byte("date"),
		[]byte("last-modified"),
		[]byte("refresh"), // URL rewrite
		// []byte("location"), TODO URL rewrite
		[]byte("content-language"),
	}
)
