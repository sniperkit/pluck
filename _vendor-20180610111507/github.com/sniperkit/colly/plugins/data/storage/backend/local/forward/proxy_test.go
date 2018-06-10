/*
Copyright 2016 Mike Gleason jr Couturier.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package forwardcache

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/sniperkit/xcache"
)

func TestProxy(t *testing.T) {
	var origin = func(path string) string {
		return url.QueryEscape(origin.URL + path)
	}

	tests := []struct {
		method string
		req    string
		status int
		body   string
	}{
		{"GET", "/p", http.StatusBadGateway, ""},
		{"GET", "/proxy?url=" + origin("/small.js"), http.StatusBadGateway, ""},
		{"GET", "/proxy?q=" + origin("/small.js"), http.StatusOK, "console.log('test');"},
		{"GET", "/proxy?q=" + origin("/unknown.html"), http.StatusNotFound, "404 page not found\n"},
		{"GET", "/proxy?q=" + "%25", http.StatusBadGateway, ""},
	}

	cache := httpcache.NewMemoryCache()
	handler := newProxy("/proxy", cache, http.DefaultTransport)

	for _, test := range tests {
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest(test.method, test.req, nil)
		handler.ServeHTTP(rr, req)

		if rr.Code != test.status {
			t.Errorf("proxy sent wrong status: got %d want %d", rr.Code, test.status)
		}

		if rr.Body.String() != test.body {
			t.Errorf("proxy returned unexpected body: got %s want %s", rr.Body.String(), test.body)
		}
	}
}

func BenchmarkProxy(b *testing.B) {
	b.ReportAllocs()

	cache := httpcache.NewMemoryCache()
	handler := newProxy("/proxy", cache, http.DefaultTransport)
	discard := discarder{}
	req, err := http.NewRequest("GET", "/proxy?q="+url.QueryEscape(origin.URL+"/jquery-3.1.1.js"), nil)

	if err != nil {
		b.Error(err)
		return
	}

	for i := 0; i < b.N; i++ {
		handler.ServeHTTP(discard, req)
	}
}

type discarder struct{}

var header = make(map[string][]string)

func (discarder) Header() http.Header         { return header }
func (discarder) Write(b []byte) (int, error) { return len(b), nil }
func (discarder) WriteHeader(int)             {}
