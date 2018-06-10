package morty

import (
	"bytes"
	"fmt"
	"net/url"
)

type RequestConfig struct {
	Key     []byte
	BaseURL *url.URL
}

func (rc *RequestConfig) ProxifyURI(uri []byte) (string, error) {
	// sanitize URI
	uri, scheme := sanitizeURI(uri)

	// remove javascript protocol
	if scheme == "javascript:" {
		return "", nil
	}

	// TODO check malicious data: - e.g. data:script
	if scheme == "data:" {
		if bytes.HasPrefix(uri, []byte("data:image/png")) ||
			bytes.HasPrefix(uri, []byte("data:image/jpeg")) ||
			bytes.HasPrefix(uri, []byte("data:image/pjpeg")) ||
			bytes.HasPrefix(uri, []byte("data:image/gif")) ||
			bytes.HasPrefix(uri, []byte("data:image/webp")) {
			// should be safe
			return string(uri), nil
		} else {
			// unsafe data
			return "", nil
		}
	}

	// parse the uri
	u, err := url.Parse(string(uri))
	if err != nil {
		return "", err
	}

	// get the fragment (with the prefix "#")
	fragment := ""
	if len(u.Fragment) > 0 {
		fragment = "#" + u.Fragment
	}

	// reset the fragment: it is not included in the mortyurl
	u.Fragment = ""

	// merge the URI with the document URI
	u = mergeURIs(rc.BaseURL, u)

	// simple internal link ?
	// some web pages describe the whole link https://same:auth@same.host/same.path?same.query#new.fragment
	if u.Scheme == rc.BaseURL.Scheme &&
		(rc.BaseURL.User == nil || (u.User != nil && u.User.String() == rc.BaseURL.User.String())) &&
		u.Host == rc.BaseURL.Host &&
		u.Path == rc.BaseURL.Path &&
		u.RawQuery == rc.BaseURL.RawQuery {
		// the fragment is the only difference between the document URI and the uri parameter
		return fragment, nil
	}

	// return full URI and fragment (if not empty)
	morty_uri := u.String()

	if rc.Key == nil {
		return fmt.Sprintf("./?mortyurl=%s%s", url.QueryEscape(morty_uri), fragment), nil
	}
	return fmt.Sprintf("./?mortyhash=%s&mortyurl=%s%s", hash(morty_uri, rc.Key), url.QueryEscape(morty_uri), fragment), nil
}

func appRequestHandler(ctx *fasthttp.RequestCtx) bool {
	// serve robots.txt
	if bytes.Equal(ctx.Path(), []byte("/robots.txt")) {
		ctx.SetContentType("text/plain")
		ctx.Write([]byte("User-Agent: *\nDisallow: /\n"))
		return true
	}
	// server favicon.ico
	if bytes.Equal(ctx.Path(), []byte("/favicon.ico")) {
		ctx.SetContentType("image/png")
		ctx.Write(FAVICON_BYTES)
		return true
	}
	return false
}

func popRequestParam(ctx *fasthttp.RequestCtx, paramName []byte) []byte {
	param := ctx.QueryArgs().PeekBytes(paramName)
	if param == nil {
		param = ctx.PostArgs().PeekBytes(paramName)
		if param != nil {
			ctx.PostArgs().DelBytes(paramName)
		}
	} else {
		ctx.QueryArgs().DelBytes(paramName)
	}
	return param
}
