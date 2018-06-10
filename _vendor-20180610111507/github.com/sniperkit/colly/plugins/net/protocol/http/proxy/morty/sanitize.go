package morty

import (
	"bytes"
	"fmt"
	"io"

	"golang.org/x/net/html"
)

// Sanitized URI : removes all runes bellow 32 (included) as the begining and end of URI, and lower case the scheme.
// avoid memory allocation (except for the scheme)
func sanitizeURI(uri []byte) ([]byte, string) {
	first_rune_index := 0
	first_rune_seen := false
	scheme_last_index := -1
	buffer := bytes.NewBuffer(make([]byte, 0, 10))

	// remove trailing space and special characters
	uri = bytes.TrimRight(uri, "\x00\x01\x02\x03\x04\x05\x06\x07\x08\x09\x0A\x0B\x0C\x0D\x0E\x0F\x10\x11\x12\x13\x14\x15\x16\x17\x18\x19\x1A\x1B\x1C\x1D\x1E\x1F\x20")

	// loop over byte by byte
	for i, c := range uri {
		// ignore special characters and space (c <= 32)
		if c > 32 {
			// append to the lower case of the rune to buffer
			if c < utf8.RuneSelf && 'A' <= c && c <= 'Z' {
				c = c + 'a' - 'A'
			}

			buffer.WriteByte(c)

			// update the first rune index that is not a special rune
			if !first_rune_seen {
				first_rune_index = i
				first_rune_seen = true
			}

			if c == ':' {
				// colon rune found, we have found the scheme
				scheme_last_index = i
				break
			} else if c == '/' || c == '?' || c == '\\' || c == '#' {
				// special case : most probably a relative URI
				break
			}
		}
	}

	if scheme_last_index != -1 {
		// scheme found
		// copy the "lower case without special runes scheme" before the ":" rune
		scheme_start_index := scheme_last_index - buffer.Len() + 1
		copy(uri[scheme_start_index:], buffer.Bytes())
		// and return the result
		return uri[scheme_start_index:], buffer.String()
	} else {
		// scheme NOT found
		return uri[first_rune_index:], ""
	}
}

func sanitizeCSS(rc *RequestConfig, out io.Writer, css []byte) {
	// TODO
	urlSlices := CSS_URL_REGEXP.FindAllSubmatchIndex(css, -1)
	if urlSlices == nil {
		out.Write(css)
		return
	}
	startIndex := 0
	for _, s := range urlSlices {
		urlStart := s[4]
		urlEnd := s[5]

		if uri, err := rc.ProxifyURI(css[urlStart:urlEnd]); err == nil {
			out.Write(css[startIndex:urlStart])
			out.Write([]byte(uri))
			startIndex = urlEnd
		} else {
			log.Println("cannot proxify css uri:", string(css[urlStart:urlEnd]))
		}
	}
	if startIndex < len(css) {
		out.Write(css[startIndex:len(css)])
	}
}

func sanitizeHTML(rc *RequestConfig, out io.Writer, htmlDoc []byte) {
	r := bytes.NewReader(htmlDoc)
	decoder := html.NewTokenizer(r)
	decoder.AllowCDATA(true)
	unsafeElements := make([][]byte, 0, 8)
	state := STATE_DEFAULT
	for {
		token := decoder.Next()
		if token == html.ErrorToken {
			err := decoder.Err()
			if err != io.EOF {
				log.Println("failed to parse HTML:")
			}
			break
		}
		if len(unsafeElements) == 0 {
			switch token {
			case html.StartTagToken, html.SelfClosingTagToken:
				tag, hasAttrs := decoder.TagName()
				safe := !inArray(tag, UNSAFE_ELEMENTS)
				if !safe {
					if !inArray(tag, SELF_CLOSING_ELEMENTS) {
						var unsafeTag []byte = make([]byte, len(tag))
						copy(unsafeTag, tag)
						unsafeElements = append(unsafeElements, unsafeTag)
					}
					break
				}
				if bytes.Equal(tag, []byte("base")) {
					for {
						attrName, attrValue, moreAttr := decoder.TagAttr()
						if bytes.Equal(attrName, []byte("href")) {
							parsedURI, err := url.Parse(string(attrValue))
							if err == nil {
								rc.BaseURL = parsedURI
							}
						}
						if !moreAttr {
							break
						}
					}
					break
				}
				if bytes.Equal(tag, []byte("noscript")) {
					state = STATE_IN_NOSCRIPT
					break
				}
				var attrs [][][]byte
				if hasAttrs {
					for {
						attrName, attrValue, moreAttr := decoder.TagAttr()
						attrs = append(attrs, [][]byte{
							attrName,
							attrValue,
							[]byte(html.EscapeString(string(attrValue))),
						})
						if !moreAttr {
							break
						}
					}
				}
				if bytes.Equal(tag, []byte("link")) {
					sanitizeLinkTag(rc, out, attrs)
					break
				}

				if bytes.Equal(tag, []byte("meta")) {
					sanitizeMetaTag(rc, out, attrs)
					break
				}

				fmt.Fprintf(out, "<%s", tag)
				if hasAttrs {
					sanitizeAttrs(rc, out, attrs)
				}

				if token == html.SelfClosingTagToken {
					fmt.Fprintf(out, " />")
				} else {
					fmt.Fprintf(out, ">")
					if bytes.Equal(tag, []byte("style")) {
						state = STATE_IN_STYLE
					}
				}

				if bytes.Equal(tag, []byte("head")) {
					fmt.Fprintf(out, HTML_HEAD_CONTENT_TYPE)
				}

				if bytes.Equal(tag, []byte("form")) {
					var formURL *url.URL
					for _, attr := range attrs {
						if bytes.Equal(attr[0], []byte("action")) {
							formURL, _ = url.Parse(string(attr[1]))
							formURL = mergeURIs(rc.BaseURL, formURL)
							break
						}
					}
					if formURL == nil {
						formURL = rc.BaseURL
					}
					urlStr := formURL.String()
					var key string
					if rc.Key != nil {
						key = hash(urlStr, rc.Key)
					}
					fmt.Fprintf(out, HTML_FORM_EXTENSION, urlStr, key)

				}

			case html.EndTagToken:
				tag, _ := decoder.TagName()
				writeEndTag := true
				switch string(tag) {
				case "body":
					fmt.Fprintf(out, HTML_BODY_EXTENSION, rc.BaseURL.String())
				case "style":
					state = STATE_DEFAULT
				case "noscript":
					state = STATE_DEFAULT
					writeEndTag = false
				}
				// skip noscript tags - only the tag, not the content, because javascript is sanitized
				if writeEndTag {
					fmt.Fprintf(out, "</%s>", tag)
				}

			case html.TextToken:
				switch state {
				case STATE_DEFAULT:
					fmt.Fprintf(out, "%s", decoder.Raw())
				case STATE_IN_STYLE:
					sanitizeCSS(rc, out, decoder.Raw())
				case STATE_IN_NOSCRIPT:
					sanitizeHTML(rc, out, decoder.Raw())
				}

			case html.CommentToken:
				// ignore comment. TODO : parse IE conditional comment

			case html.DoctypeToken:
				out.Write(decoder.Raw())
			}
		} else {
			switch token {
			case html.StartTagToken:
				tag, _ := decoder.TagName()
				if inArray(tag, UNSAFE_ELEMENTS) {
					unsafeElements = append(unsafeElements, tag)
				}

			case html.EndTagToken:
				tag, _ := decoder.TagName()
				if bytes.Equal(unsafeElements[len(unsafeElements)-1], tag) {
					unsafeElements = unsafeElements[:len(unsafeElements)-1]
				}
			}
		}
	}
}

func sanitizeLinkTag(rc *RequestConfig, out io.Writer, attrs [][][]byte) {
	exclude := false
	for _, attr := range attrs {
		attrName := attr[0]
		attrValue := attr[1]
		if bytes.Equal(attrName, []byte("rel")) {
			if !inArray(attrValue, LINK_REL_SAFE_VALUES) {
				exclude = true
				break
			}
		}
		if bytes.Equal(attrName, []byte("as")) {
			if bytes.Equal(attrValue, []byte("script")) {
				exclude = true
				break
			}
		}
	}

	if !exclude {
		out.Write([]byte("<link"))
		for _, attr := range attrs {
			sanitizeAttr(rc, out, attr[0], attr[1], attr[2])
		}
		out.Write([]byte(">"))
	}
}

func sanitizeMetaTag(rc *RequestConfig, out io.Writer, attrs [][][]byte) {
	var http_equiv []byte
	var content []byte

	for _, attr := range attrs {
		attrName := attr[0]
		attrValue := attr[1]
		if bytes.Equal(attrName, []byte("http-equiv")) {
			http_equiv = bytes.ToLower(attrValue)
			// exclude some <meta http-equiv="..." ..>
			if !inArray(http_equiv, LINK_HTTP_EQUIV_SAFE_VALUES) {
				return
			}
		}
		if bytes.Equal(attrName, []byte("content")) {
			content = attrValue
		}
		if bytes.Equal(attrName, []byte("charset")) {
			// exclude <meta charset="...">
			return
		}
	}

	out.Write([]byte("<meta"))
	urlIndex := bytes.Index(bytes.ToLower(content), []byte("url="))
	if bytes.Equal(http_equiv, []byte("refresh")) && urlIndex != -1 {
		contentUrl := content[urlIndex+4:]
		// special case of <meta http-equiv="refresh" content="0; url='example.com/url.with.quote.outside'">
		if len(contentUrl) >= 2 && (contentUrl[0] == byte('\'') || contentUrl[0] == byte('"')) {
			if contentUrl[0] == contentUrl[len(contentUrl)-1] {
				contentUrl = contentUrl[1 : len(contentUrl)-1]
			}
		}
		// output proxify result
		if uri, err := rc.ProxifyURI(contentUrl); err == nil {
			fmt.Fprintf(out, ` http-equiv="refresh" content="%surl=%s"`, content[:urlIndex], uri)
		}
	} else {
		if len(http_equiv) > 0 {
			fmt.Fprintf(out, ` http-equiv="%s"`, http_equiv)
		}
		sanitizeAttrs(rc, out, attrs)
	}
	out.Write([]byte(">"))
}

func sanitizeAttrs(rc *RequestConfig, out io.Writer, attrs [][][]byte) {
	for _, attr := range attrs {
		sanitizeAttr(rc, out, attr[0], attr[1], attr[2])
	}
}

func sanitizeAttr(rc *RequestConfig, out io.Writer, attrName, attrValue, escapedAttrValue []byte) {
	if inArray(attrName, SAFE_ATTRIBUTES) {
		fmt.Fprintf(out, " %s=\"%s\"", attrName, escapedAttrValue)
		return
	}
	switch string(attrName) {
	case "src", "href", "action":
		if uri, err := rc.ProxifyURI(attrValue); err == nil {
			fmt.Fprintf(out, " %s=\"%s\"", attrName, uri)
		} else {
			log.Println("cannot proxify uri:", string(attrValue))
		}
	case "style":
		cssAttr := bytes.NewBuffer(nil)
		sanitizeCSS(rc, cssAttr, attrValue)
		fmt.Fprintf(out, " %s=\"%s\"", attrName, html.EscapeString(string(cssAttr.Bytes())))
	}
}
