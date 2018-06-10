package sitemap

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"compress/zlib"
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"

	"path/filepath"
	"strings"
	"sync"
	"time"

	colly "github.com/sniperkit/colly/pkg"
	queue "github.com/sniperkit/colly/pkg/queue"
)

type SitemapCollector struct {
	Name         xml.Name `xml:"urlset,sitemapindex"`
	NS           string   `xml:"xmlns,attr"`
	Indices      []url.URL
	URLs         []url.URL
	href         url.URL
	converted    map[string]string
	content      []byte
	contentType  string
	contentSize  int
	statusCode   int
	prefixPath   string
	localPath    string
	total_size   int
	startTime    time.Time
	endTime      time.Time
	responseTime float64
	collector    *colly.Collector
	cqueue       *queue.Queue
	lock         *sync.RWMutex
	wg           *sync.WaitGroup
	// Indices []Index `json:"loc" xml:"sitemap"`
	// URLs    []URL   `json:"url" xml:"url"`
}

type URL struct {
	href       url.URL
	Loc        string `json:"loc" xml:"loc" csv:"loc"`
	LastMod    string `json:"lastmod" xml:"lastmod" csv:"lastmod"`
	ChangeFreq string `default:"weekly" json:"changefreq" xml:"changefreq" csv:"changefreq"`
	Priority   string `default:"1.0" json:"priority" xml:"priority" csv:"priority"`
}

type Index struct {
	Loc     string `json:"loc" xml:"loc"`
	LastMod string `json:"lastmod" xml:"lastmod"`
}

type Indices struct {
	Sitemaps []Index `xml:"sitemap" json:"sitemap"`
}

type IndexError struct {
	message string
}

type XMLSitemap struct {
	URLs []URL `xml:"url"`
}

type XmlSitemapError struct {
	message string
}

type TXTSitemap struct {
	URLs []URL `csv:"url"`
}

type TxtSitemapError struct {
	message string
}

func New(inputURL string) (*SitemapCollector, error) {
	sitemapURL, err := url.Parse(inputURL)
	if err != nil {
		return nil, err
	}

	ouput, err := readURL(*sitemapURL)
	if err != nil {
		return nil, err
	}

	s := &SitemapCollector{
		href:         *sitemapURL,
		content:      ouput.Body,
		statusCode:   ouput.StatusCode,
		contentType:  ouput.ContentType,
		contentSize:  len(ouput.Body),
		responseTime: ouput.EndTime.Sub(ouput.StartTime).Seconds(),
		// lock:         &sync.RWMutex{},
		// wg:           &sync.WaitGroup{},
	}

	s.getURLs()
	return s, nil
}

func (s *SitemapCollector) IsValid() bool {
	return s.href.String() != ""
}

func (s *SitemapCollector) Read() error {
	return errInvalidContent
}

func (s *SitemapCollector) Print(format string) error {
	return errInvalidContent
}

func getXMLSitemap(xmlSitemapURL url.URL) (XMLSitemap, error) {
	response, readErr := readURL(xmlSitemapURL)
	if readErr != nil {
		return XMLSitemap{}, readErr
	}

	if !strings.Contains(string(response.GetBody()), "</urlset>") {
		return XMLSitemap{}, XmlSitemapError{"Invalid content"}
	}

	var urlSet XMLSitemap
	unmarshalError := xml.Unmarshal(response.GetBody(), &urlSet)
	if unmarshalError != nil {
		return XMLSitemap{}, unmarshalError
	}
	return urlSet, nil
}

func (sitemapIndexError XmlSitemapError) Error() string {
	return sitemapIndexError.message
}

func isInvalidSitemapContent(err error) bool {
	if err == nil {
		return false
	}
	return err.Error() == "Invalid content"
}

func checkSitemap(loc string) bool {
	return true
}

func readURL(url url.URL) (colly.Response, error) {

	startTime := time.Now().UTC()
	resp, fetchErr := http.Get(url.String())
	if fetchErr != nil {
		return colly.Response{}, fetchErr
	}

	defer resp.Body.Close()

	var body []byte
	var errReader error

	// content type
	contentType := resp.Header.Get("Content-Type")

	body, errReader = ioutil.ReadAll(resp.Body)
	if errReader != nil {
		if log != nil {
			log.Fatalln("error.ReadAll:", contentType, ", msg=", errReader)
		}
		return colly.Response{}, errReader
	}

	//if contentType == "" {
	contentType = http.DetectContentType(body)
	//}

	pathExtension := path.Ext(url.String())
	contentEncoding := resp.Header.Get("Content-Encoding")

	if log != nil {
		log.Println("URL:", resp.Request.URL.String(), "Content-Type:", contentType, "StatusCode:", resp.StatusCode)
	}

	if resp.StatusCode == 404 {
		return colly.Response{}, errReader
	}

	switch pathExtension {
	case ".gz":
		contentType = "application/x-gzip"
	case ".txt":
		contentType = "text/plain"
	}

	switch contentEncoding {
	case "gzip":
		contentType = "application/x-gzip"
	case "deflate":
		contentType = "application/x-deflate"
	case "zlib":
		contentType = "application/x-zlib"
	}

	switch contentType {
	// "application/octet-stream", "application/x-tar"
	case "application/x-gzip", "application/gzip":
		gr, _ := gzip.NewReader(bytes.NewBuffer(body))
		defer gr.Close()

		body, errReader = ioutil.ReadAll(gr)

	case "application/x-deflate", "application/deflate":
		rdata := flate.NewReader(bytes.NewBuffer(body))
		body, errReader = ioutil.ReadAll(rdata)

	case "application/x-zlib", "application/zlib":
		var readCloser io.ReadCloser
		readCloser, errReader = zlib.NewReader(bytes.NewBuffer(body))
		if errReader != nil {
			if log != nil {
				log.Fatalln("readCloser.error:", contentType, ", msg=", errReader)
			}
			return colly.Response{}, errReader
		}
		body, errReader = ioutil.ReadAll(readCloser)

	}

	if errReader != nil {
		if log != nil {
			log.Fatalln("error.ReadAll:", contentType, ", msg=", errReader)
		}
		return colly.Response{}, errReader
	}

	endTime := time.Now().UTC()

	return colly.Response{
		Body:        body,
		StatusCode:  resp.StatusCode,
		StartTime:   startTime,
		EndTime:     endTime,
		ContentType: contentType,
	}, nil
}

// String return the string format of the sitemap
func (s *SitemapCollector) String() string {
	var items []string
	for _, item := range s.URLs {
		items = append(items, item.String())
	}
	return fmt.Sprintf(SitemapXML, strings.Join(items, `
`))
}

// ToFile saves a sitemap to a file with either extension .xml or .gz.
// If extension is .gz, the file will be gzipped.
func (s *SitemapCollector) ToFile(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	ext := filepath.Ext(file.Name())
	if ext != ".xml" && ext != ".gz" {
		return fmt.Errorf("filename %s does not have extension .xml or .gz, extension %s given", file.Name(), ext)
	}

	// Gzip
	if ext == ".gz" {
		zip := gzip.NewWriter(file)
		defer zip.Close()

		_, err = zip.Write([]byte(s.String()))
		if err != nil {
			return err
		}
	} else {
		_, err = file.Write([]byte(s.String()))
		if err != nil {
			return err
		}
	}

	return nil
}

func (sitemapIndexError TxtSitemapError) Error() string {
	return sitemapIndexError.message
}

/*
func getTXTSitemap(txtSitemapURL url.URL) (TXTSitemap, error) {
	_, readErr := readURL(txtSitemapURL)
	if readErr != nil {
		return TXTSitemap{}, readErr
	}

	var urlSet TXTSitemap
	//unmarshalError := csv.Unmarshal(response.GetBody(), &urlSet)
	//if unmarshalError != nil {
	//	return TXTSitemap{}, unmarshalError
	//}
	return urlSet, nil
}
*/

func isRemoteURL(url string) bool {
	return strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")
}

// Streams a CSV Reader into a returned channel.  Each CSV row is streamed along with the header.
// "true" is sent to the `done` channel when the file is finished.

var (
	SitemapTXT_SplitAt             int = 2500
	SitemapTXT_ReaderBuffer        int = 20000
	SitemapTXT_StreamBuffer        int = 20000
	SitemapTXT_StreamOuputMaxLines int = 1000
	// csvStreamOuputColumns  []string = []string{"domain", "loc", "created_at", "duration", "duration_time", "finished_at"}
)

type sitemapTXT_Line struct {
	header []string
	line   []string
}

func (s *sitemapTXT_Line) GetByKey(k int) (value string) {
	return strings.TrimSpace(s.line[k])
}

func (s *sitemapTXT_Line) GetByName(key string) (value string) {
	x := -1
	for i, value := range s.header {
		if value == key {
			x = i
			break
		}
	}
	if x == -1 {
		return ""
	}
	return strings.TrimSpace(s.line[x])
}

type sitemapTXT_Stream struct {
	path         string   `required:'true'` // Filepath to Local CSV File
	splitAt      int      `default:'2500'`
	buffer       int      `default:'20000'`
	selectorType string   `default:'name'` // column selector type, availble: by_key, by_name (default: column_name)
	columnsKeys  []int    // default: 0
	columnsNames []string // default: "url"
	debug        bool
	isRemote     bool
	reader       *csv.Reader
	lock         *sync.Mutex
	wg           *sync.WaitGroup
	*sitemapTXT_Line
}

type sitemapTXT_FlowTable []*sitemapTXT_FlowLine

type sitemapTXT_FlowLine struct {
	loc        string
	priority   string `default:"1.0"`
	changefreq string `default:"weekly"`
}

func getTXTSitemap(txtSitemapURL url.URL) (TXTSitemap, error) {

	var urlSet TXTSitemap
	var reader *csv.Reader

	txtSitemapLink := txtSitemapURL.String()
	if !isRemoteURL(txtSitemapLink) {
		if _, err := os.Stat(txtSitemapLink); os.IsNotExist(err) {
			return urlSet, err
		}
		file, err := os.Open(txtSitemapLink)
		if err != nil {
			return urlSet, err
		}
		defer file.Close()
		reader = csv.NewReader(file)

	} else {

		resp, err := http.Get(txtSitemapLink)
		if err != nil {
			return urlSet, err
		}
		defer resp.Body.Close()

		if resp.StatusCode < 200 || resp.StatusCode >= 400 {
			return urlSet, err
		}
		reader = csv.NewReader(resp.Body)

	}

	// urlSet = make(TXTSitemap, 0)
	lines := streamSitemapTXT(reader, SitemapTXT_StreamBuffer)
	for line := range lines {
		if loc := line.GetByKey(0); loc != "" {
			href, err := url.Parse(loc)
			if err == nil {
				// fmt.Println("url=", href.String())
				urlSet.URLs = append(urlSet.URLs, URL{href: *href, Loc: loc})
			}
		}
	}

	return urlSet, nil
}

// Args
//  csv    - The csv.Reader that will be read from.
//  buffer - The "lines" buffer factor.  Send "0" for an unbuffered channel.
func streamSitemapTXT(csv *csv.Reader, buffer int) (lines chan *sitemapTXT_Line) {
	lines = make(chan *sitemapTXT_Line, buffer)
	go func() {
		header, err := csv.Read()
		if err != nil {
			close(lines)
			return
		}
		i := 0
		for {
			line, err := csv.Read()
			if len(line) > 0 {
				i++
				lines <- &sitemapTXT_Line{
					header: header,
					line:   line,
				}
			}
			if err != nil {
				close(lines)
				return
			}
		}
	}()
	return
}

func printSitemapTXT_FlowTable(lines chan *sitemapTXT_FlowLine) (done chan int) {
	done = make(chan int)
	go func() {
		table := sitemapTXT_FlowTable{}
		i := 0
		for line := range lines {
			i++
			table = append(table, line)
			if len(table) >= SitemapTXT_StreamOuputMaxLines {
				table.Send()
				table = sitemapTXT_FlowTable{}
			}
		}
		if len(table) > 0 {
			table.Send()
		}
		done <- i
	}()
	return
}

func printSitemapTXT_Entry(txtLines chan *sitemapTXT_Line) (lines chan *sitemapTXT_FlowLine) {
	lines = make(chan *sitemapTXT_FlowLine, SitemapTXT_ReaderBuffer)
	go func() {
		var flowLine *sitemapTXT_FlowLine
		for line := range txtLines {
			flowLine, _ = streamSitemapTXT_Entry(line)
			lines <- flowLine
		}
		close(lines)
	}()
	return
}

func streamSitemapTXT_Entry(line *sitemapTXT_Line) (*sitemapTXT_FlowLine, error) {
	sfl := sitemapTXT_FlowLine{}
	sfl.priority = line.GetByName("priority")
	sfl.changefreq = line.GetByName("changefreq")
	sfl.loc = line.GetByName("loc")
	return &sfl, nil
}

func (sft *sitemapTXT_FlowTable) Send() {
	// code to send to the database here.
	fmt.Printf("----\nSending %d lines\n%s", len(*sft), *sft)
}
