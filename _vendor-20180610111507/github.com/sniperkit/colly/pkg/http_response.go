// Copyright 2018 Adam Tauber
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package colly

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"
	"strings"
	"time"

	"github.com/saintfish/chardet"
	"golang.org/x/net/html/charset"
)

// Response is the representation of a HTTP response made by a Collector
type Response struct {
	// StatusCode is the status code of the Response
	StatusCode int
	// Body is the content of the Response
	Body []byte
	// Ctx is a context between a Request and a Response
	Ctx *Context
	// Request is the Request object of the response
	Request *Request
	// Headers contains the Response's HTTP headers
	Headers *http.Header
	//
	StartTime time.Time
	//
	EndTime time.Time
	//
	ContentType string
	//
	Err error
}

// Save writes response body to disk
func (r *Response) Save(fileName string) error {
	return ioutil.WriteFile(fileName, r.Body, 0644)
}

// GetBody returns response body as byte array
func (r *Response) GetBody() []byte {
	return r.Body
}

func (r *Response) GetSize() int {
	return len(r.Body)
}

func (r *Response) GetContentType() string {
	return r.ContentType
}

func (r *Response) GetStatusCode() int {
	return r.StatusCode
}

func (r *Response) GetStartTime() time.Time {
	return r.StartTime
}

func (r *Response) GetEndTime() time.Time {
	return r.EndTime
}

func (r *Response) IsHTML() bool {
	return strings.HasPrefix(r.ContentType, "text/html")
}

// to finish properly
func (r *Response) IsXML() bool {
	return strings.HasPrefix(r.ContentType, "application/xml")
}

func (r *Response) IsCSV() bool {
	return strings.HasPrefix(r.ContentType, "application/csv")
}

func (r *Response) IsJSON() bool {
	return strings.HasPrefix(r.ContentType, "application/json")
}

func (r *Response) IsYAML() bool {
	return strings.HasPrefix(r.ContentType, "application/yaml")
}

// FileName returns the sanitized file name parsed from "Content-Disposition"
// header or from URL
func (r *Response) FileName() string {
	_, params, err := mime.ParseMediaType(r.Headers.Get("Content-Disposition"))
	if fName, ok := params["filename"]; ok && err == nil {
		return SanitizeFileName(fName)
	}
	if r.Request.URL.RawQuery != "" {
		return SanitizeFileName(fmt.Sprintf("%s_%s", r.Request.URL.Path, r.Request.URL.RawQuery))
	}
	return SanitizeFileName(strings.TrimPrefix(r.Request.URL.Path, "/"))
}

func (r *Response) fixCharset(detectCharset bool, defaultEncoding string) error {
	if defaultEncoding != "" {
		tmpBody, err := encodeBytes(r.Body, "text/plain; charset="+defaultEncoding)
		if err != nil {
			return err
		}
		r.Body = tmpBody
		return nil
	}
	contentType := strings.ToLower(r.Headers.Get("Content-Type"))
	if !strings.Contains(contentType, "charset") {
		if !detectCharset {
			return nil
		}
		d := chardet.NewTextDetector()
		r, err := d.DetectBest(r.Body)
		if err != nil {
			return err
		}
		contentType = "text/plain; charset=" + r.Charset
	}
	if strings.Contains(contentType, "utf-8") || strings.Contains(contentType, "utf8") {
		return nil
	}
	tmpBody, err := encodeBytes(r.Body, contentType)
	if err != nil {
		return err
	}
	r.Body = tmpBody
	return nil
}

func encodeBytes(b []byte, contentType string) ([]byte, error) {
	r, err := charset.NewReader(bytes.NewReader(b), contentType)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(r)
}
