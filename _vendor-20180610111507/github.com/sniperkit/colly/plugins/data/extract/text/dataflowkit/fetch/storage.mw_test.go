package fetch

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/slotix/dataflowkit/storage"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	mw storageMiddleware
	s  storage.Store
)

func init() {
	var svc Service
	svc = FetchService{}
	s = storage.NewStore(storage.Diskv)
	mw = storageMiddleware{
		storage: s,
		Service: svc,
	}
}

func Test_storageMiddleware(t *testing.T) {
	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Conent-Type", "text/html")
		w.Write(IndexContent)
	})
	ts := httptest.NewServer(r)
	defer ts.Close()

	req := BaseFetcherRequest{
		//URL:      "http://example.com",
		URL:      ts.URL,
		FormData: "",
	}
	//Loading from remote server
	start := time.Now()
	resp, err := mw.Response(req)
	assert.Nil(t, err, "Expected no error")
	assert.Equal(t, 200, resp.(*BaseFetcherResponse).StatusCode, "Expected Fetcher returns 200 status code")
	elapsed1 := time.Since(start)
	t.Log("Loading from remote server... ", elapsed1)

	//Loading from cached storage
	start = time.Now()
	resp, err = mw.Response(req)
	assert.Nil(t, err, "Expected no error")
	assert.Equal(t, 200, resp.(*BaseFetcherResponse).StatusCode, "Expected Splash server returns 200 status code")
	elapsed2 := time.Since(start)
	t.Log("Loading from remote server... ", elapsed2)
	//assert.Equal(t, true, elapsed1 > elapsed2, "it takes longer to load a webpage from remote server")

	err = s.DeleteAll()
	assert.Nil(t, err, "Expected no error")
}

func Test_IGNORE_CACHE_INFO(t *testing.T) {
	viper.Set("IGNORE_CACHE_INFO", true)
	viper.Set("STORAGE_TYPE", "Diskv")

	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Conent-Type", "text/html")
		w.Write(IndexContent)
	})
	ts := httptest.NewServer(r)
	defer ts.Close()
	req := BaseFetcherRequest{
		//URL:      "http://google.com",
		URL:       ts.URL,
		FormData:  "",
		UserToken: "12345",
	}
	//Loading from remote server
	start := time.Now()
	resp, err := mw.Response(req)
	assert.Nil(t, err, "Expected no error")
	assert.Equal(t, 200, resp.(*BaseFetcherResponse).StatusCode, "Expected Splash server returns 200 status code")
	elapsed1 := time.Since(start)
	t.Log("Loading from remote server... ", elapsed1)

	//Loading from cached storage
	start = time.Now()
	resp, err = mw.Response(req)
	assert.Nil(t, err, "Expected no error")
	assert.Equal(t, 200, resp.(*BaseFetcherResponse).StatusCode, "Expected Splash server returns 200 status code")
	elapsed2 := time.Since(start)
	t.Log("Loading from remote server... ", elapsed2)
	assert.Equal(t, true, elapsed1 > elapsed2, "it takes longer to load a webpage from remote server")

	err = s.DeleteAll()
	assert.Nil(t, err, "Expected no error")
}
