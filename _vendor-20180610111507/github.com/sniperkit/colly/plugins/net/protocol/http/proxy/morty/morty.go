package morty

import (
	"log"
	"net/url"

	"golang.org/x/net/html"

	"github.com/valyala/fasthttp"
)

const VERSION = "v0.2.0"

const (
	STATE_DEFAULT     int = 0
	STATE_IN_STYLE    int = 1
	STATE_IN_NOSCRIPT int = 2
)

var (
	client *fasthttp.Client = &fasthttp.Client{
		MaxResponseBodySize: 10 * 1024 * 1024, // 10M
	}
)

type Proxy struct {
	Key            []byte
	RequestTimeout time.Duration
}

func (p *Proxy) RequestHandler(ctx *fasthttp.RequestCtx) {

	if appRequestHandler(ctx) {
		return
	}

	requestHash := popRequestParam(ctx, []byte("mortyhash"))

	requestURI := popRequestParam(ctx, []byte("mortyurl"))

	if requestURI == nil {
		p.serveMainPage(ctx, 200, nil)
		return
	}

	if p.Key != nil {
		if !verifyRequestURI(requestURI, requestHash, p.Key) {
			// HTTP status code 403 : Forbidden
			p.serveMainPage(ctx, 403, errors.New(`invalid "mortyhash" parameter`))
			return
		}
	}

	parsedURI, err := url.Parse(string(requestURI))

	if err != nil {
		// HTTP status code 500 : Internal Server Error
		p.serveMainPage(ctx, 500, err)
		return
	}

	// Serve an intermediate page for protocols other than HTTP(S)
	if (parsedURI.Scheme != "http" && parsedURI.Scheme != "https") || strings.HasSuffix(parsedURI.Host, ".onion") {
		p.serveExitMortyPage(ctx, parsedURI)
		return
	}

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.SetConnectionClose()

	requestURIStr := string(requestURI)

	log.Println("getting", requestURIStr)

	req.SetRequestURI(requestURIStr)
	req.Header.SetUserAgentBytes([]byte("Mozilla/5.0 (Windows NT 10.0; WOW64; rv:50.0) Gecko/20100101 Firefox/50.0"))

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	req.Header.SetMethodBytes(ctx.Method())
	if ctx.IsPost() || ctx.IsPut() {
		req.SetBody(ctx.PostBody())
	}

	err = CLIENT.DoTimeout(req, resp, p.RequestTimeout)

	if err != nil {
		if err == fasthttp.ErrTimeout {
			// HTTP status code 504 : Gateway Time-Out
			p.serveMainPage(ctx, 504, err)
		} else {
			// HTTP status code 500 : Internal Server Error
			p.serveMainPage(ctx, 500, err)
		}
		return
	}

	if resp.StatusCode() != 200 {
		switch resp.StatusCode() {
		case 301, 302, 303, 307, 308:
			loc := resp.Header.Peek("Location")
			if loc != nil {
				rc := &RequestConfig{Key: p.Key, BaseURL: parsedURI}
				url, err := rc.ProxifyURI(loc)
				if err == nil {
					ctx.SetStatusCode(resp.StatusCode())
					ctx.Response.Header.Add("Location", url)
					log.Println("redirect to", string(loc))
					return
				}
			}
		}
		error_message := fmt.Sprintf("invalid response: %d (%s)", resp.StatusCode(), requestURIStr)
		p.serveMainPage(ctx, resp.StatusCode(), errors.New(error_message))
		return
	}

	contentTypeBytes := resp.Header.Peek("Content-Type")

	if contentTypeBytes == nil {
		// HTTP status code 503 : Service Unavailable
		p.serveMainPage(ctx, 503, errors.New("invalid content type"))
		return
	}

	contentTypeString := string(contentTypeBytes)

	// decode Content-Type header
	contentType, error := contenttype.ParseContentType(contentTypeString)
	if error != nil {
		// HTTP status code 503 : Service Unavailable
		p.serveMainPage(ctx, 503, errors.New("invalid content type"))
		return
	}

	// content-disposition
	contentDispositionBytes := ctx.Request.Header.Peek("Content-Disposition")

	// check content type
	if !ALLOWED_CONTENTTYPE_FILTER(contentType) {
		// it is not a usual content type
		if ALLOWED_CONTENTTYPE_ATTACHMENT_FILTER(contentType) {
			// force attachment for allowed content type
			contentDispositionBytes = contentDispositionForceAttachment(contentDispositionBytes, parsedURI)
		} else {
			// deny access to forbidden content type
			// HTTP status code 403 : Forbidden
			p.serveMainPage(ctx, 403, errors.New("forbidden content type"))
			return
		}
	}

	// HACK : replace */xhtml by text/html
	if contentType.SubType == "xhtml" {
		contentType.TopLevelType = "text"
		contentType.SubType = "html"
		contentType.Suffix = ""
	}

	// conversion to UTF-8
	var responseBody []byte

	if contentType.TopLevelType == "text" {
		e, ename, _ := charset.DetermineEncoding(resp.Body(), contentTypeString)
		if (e != encoding.Nop) && (!strings.EqualFold("utf-8", ename)) {
			responseBody, err = e.NewDecoder().Bytes(resp.Body())
			if err != nil {
				// HTTP status code 503 : Service Unavailable
				p.serveMainPage(ctx, 503, err)
				return
			}
		} else {
			responseBody = resp.Body()
		}
		// update the charset or specify it
		contentType.Parameters["charset"] = "UTF-8"
	} else {
		responseBody = resp.Body()
	}

	//
	contentType.FilterParameters(ALLOWED_CONTENTTYPE_PARAMETERS)

	// set the content type
	ctx.SetContentType(contentType.String())

	// output according to MIME type
	switch {
	case contentType.SubType == "css" && contentType.Suffix == "":
		sanitizeCSS(&RequestConfig{Key: p.Key, BaseURL: parsedURI}, ctx, responseBody)
	case contentType.SubType == "html" && contentType.Suffix == "":
		sanitizeHTML(&RequestConfig{Key: p.Key, BaseURL: parsedURI}, ctx, responseBody)
	default:
		if contentDispositionBytes != nil {
			ctx.Response.Header.AddBytesV("Content-Disposition", contentDispositionBytes)
		}
		ctx.Write(responseBody)
	}
}

func (p *Proxy) serveExitMortyPage(ctx *fasthttp.RequestCtx, uri *url.URL) {
	ctx.SetContentType("text/html")
	ctx.SetStatusCode(403)
	ctx.Write([]byte(MORTY_HTML_PAGE_START))
	ctx.Write([]byte("<h2>You are about to exit MortyProxy</h2>"))
	ctx.Write([]byte("<p>Following</p><p><a href=\""))
	ctx.Write([]byte(html.EscapeString(uri.String())))
	ctx.Write([]byte("\" rel=\"noreferrer\">"))
	ctx.Write([]byte(html.EscapeString(uri.String())))
	ctx.Write([]byte("</a></p><p>the content of this URL will be <b>NOT</b> sanitized.</p>"))
	ctx.Write([]byte(MORTY_HTML_PAGE_END))
}

func (p *Proxy) serveMainPage(ctx *fasthttp.RequestCtx, statusCode int, err error) {
	ctx.SetContentType("text/html; charset=UTF-8")
	ctx.SetStatusCode(statusCode)
	ctx.Write([]byte(MORTY_HTML_PAGE_START))
	if err != nil {
		log.Println("error:", err)
		ctx.Write([]byte("<h2>Error: "))
		ctx.Write([]byte(html.EscapeString(err.Error())))
		ctx.Write([]byte("</h2>"))
	}
	if p.Key == nil {
		ctx.Write([]byte(`
		<form action="post">
		Visit url: <input placeholder="https://url.." name="mortyurl" autofocus />
		<input type="submit" value="go" />
		</form>`))
	} else {
		ctx.Write([]byte(`<h3>Warning! This instance does not support direct URL opening.</h3>`))
	}
	ctx.Write([]byte(MORTY_HTML_PAGE_END))
}
