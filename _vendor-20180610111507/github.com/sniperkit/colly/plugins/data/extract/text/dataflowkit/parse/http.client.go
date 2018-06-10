package parse

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/slotix/dataflowkit/scrape"
)

// NewHTTPClient returns an Parse Service backed by an HTTP server living at the
// remote instance. We expect instance to come from a service discovery system,
// so likely of the form "host:port". We bake-in certain middlewares,
// implementing the client library pattern.
func NewHTTPClient(instance string) (Service, error) {
	// Quickly sanitize the instance string.
	if !strings.HasPrefix(instance, "http") {
		instance = "http://" + instance
	}
	u, err := url.Parse(instance)
	if err != nil {
		return nil, err
	}

	// Each individual endpoint is an http/transport.Client (which implements
	// endpoint.Endpoint) that gets wrapped with various middlewares. If you
	// made your own client library, you'd do this work there, so your server
	// could rely on a consistent set of client behavior.
	var parseEndpoint endpoint.Endpoint
	{
		parseEndpoint = httptransport.NewClient(
			"POST",
			copyURL(u, "/parse"),
			encodeParseRequest,
			decodeParseResponse,
		).Endpoint()
	}

	// Returning the endpoint.Set as a service.Service relies on the
	// endpoint.Set implementing the Service methods. That's just a simple bit
	// of glue code.
	return Endpoints{
		ParseEndpoint: parseEndpoint,
	}, nil
}

// encodeParseRequest is a transport/http.EncodeRequestFunc that
// JSON-encodes any request to the request body. Primarily useful in a client.
func encodeParseRequest(ctx context.Context, r *http.Request, request interface{}) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return err
	}
	r.Body = ioutil.NopCloser(&buf)
	return nil
}

// decodeSplashFetcherContent is a transport/http.DecodeResponseFunc that decodes a
// JSON-encoded splash fetcher response from the HTTP response body. If the response has a
// non-200 status code, we will interpret that as an error and attempt to decode
// the specific error message from the response body. Primarily useful in a
// client.
func decodeParseResponse(ctx context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, errors.New(r.Status)
	}
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func copyURL(base *url.URL, path string) *url.URL {
	next := *base
	next.Path = path
	return &next
}

// func (e Endpoints) Fetch(req FetchRequester) (io.ReadCloser, error) {
// 	r, err := e.Response(req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return r.GetHTML()
// }

func (e Endpoints) Parse(p scrape.Payload) (io.ReadCloser, error) {
	ctx := context.Background()
	resp, err := e.ParseEndpoint(ctx, p)
	if err != nil {
		return nil, err
	}
	readCloser := ioutil.NopCloser(bytes.NewReader(resp.([]byte)))
	//	response := resp.(io.ReadCloser)
	return readCloser, nil

}
