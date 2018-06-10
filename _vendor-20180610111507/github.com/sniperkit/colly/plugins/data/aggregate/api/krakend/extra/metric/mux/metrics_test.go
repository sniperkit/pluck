package mux

import (
	"bytes"
	"context"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	krakendmetrics "github.com/devopsfaith/krakend-metrics"
	"github.com/devopsfaith/krakend/config"
	"github.com/devopsfaith/krakend/logging"
	"github.com/devopsfaith/krakend/proxy"
	"github.com/devopsfaith/krakend/router/mux"
	"github.com/rcrowley/go-metrics"
)

func TestNew(t *testing.T) {
	rand.Seed(time.Now().Unix())

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	buf := bytes.NewBuffer(make([]byte, 1024))
	l, _ := logging.NewLogger("DEBUG", buf, "")
	metric := New(ctx, 100*time.Millisecond, l)

	response := proxy.Response{Data: map[string]interface{}{}, IsComplete: true}
	max := 1000
	min := 1
	p := func(_ context.Context, _ *proxy.Request) (*proxy.Response, error) {
		time.Sleep(time.Microsecond * time.Duration(rand.Intn(max-min)+min))
		return &response, nil
	}
	hf := metric.NewHTTPHandlerFactory(mux.EndpointHandler)
	cfg := &config.EndpointConfig{
		Endpoint: "/test/{var}",
		Timeout:  10 * time.Second,
		CacheTTL: time.Second,
		Method:   "GET",
	}
	// engine.GET("/test", hf(cfg, p))
	// engine.GET("/__stats", metric.NewExpHandler())

	for i := 0; i < 100; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test/something", ioutil.NopCloser(strings.NewReader("")))
		hf(cfg, p)(w, req)
	}

	metric.Router.Aggregate()
	snapshot := metric.TakeSnapshot()

	expected := map[string]int64{
		"krakend.router.response./test/{var}.status.200.count": 100,
		"krakend.router.connected":                             0,
		"krakend.router.disconnected":                          0,
		"krakend.router.connected-total":                       100,
		"krakend.router.disconnected-total":                    100,
		"krakend.router.response./test/{var}.status":           0,
	}
	for k, v := range snapshot.Counters {
		if exp, ok := expected[k]; !ok || int(exp) != int(v) {
			t.Errorf("unexpected metric: got [%s: %d] want [%s: %d]", k, v, k, exp)
		}
	}

	if _, ok := snapshot.Histograms["krakend.router.response./test/{var}.size"]; !ok {
		t.Error("expected histogram not present")
	}

	expected = map[string]int64{
		"krakend.router.connected-gauge":    100,
		"krakend.router.disconnected-gauge": 100,
	}
	for k, exp := range expected {
		if v, ok := snapshot.Gauges[k]; !ok || int(exp) != int(v) {
			t.Errorf("unexpected metric: got [%s: %d] want [%s: %d]", k, v, k, exp)
		}
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/__stats", ioutil.NopCloser(strings.NewReader("")))
	metric.NewExpHandler().ServeHTTP(w, req)

	if w.Result().StatusCode != 200 {
		t.Errorf("unexpected status code: %d\n", w.Result().StatusCode)
	}
}

func TestNewHTTPHandler(t *testing.T) {
	registry := metrics.NewRegistry()

	rm := krakendmetrics.NewRouterMetrics(&registry)
	assertion := func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Millisecond)
		w.Header().Set("x-test", "ok")
		w.WriteHeader(200)
		w.Write([]byte("okidoki"))
	}
	h := NewHTTPHandler("test", http.HandlerFunc(assertion), rm)
	ts := httptest.NewServer(h)

	for i := 0; i < 10; i++ {
		resp, err := http.Get(ts.URL)
		if err != nil {
			t.Error(err)
		}
		if resp.Header.Get("x-test") != "ok" {
			t.Errorf("unexpected header: %s\n", resp.Header.Get("x-test"))
		}
		if resp.StatusCode != 200 {
			t.Errorf("unexpected status code: %d\n", resp.StatusCode)
		}
	}
	rm.Aggregate()
	ts.Close()

	expected := map[string]struct{}{
		"router.connected":                      {},
		"router.disconnected":                   {},
		"router.connected-gauge":                {},
		"router.disconnected-gauge":             {},
		"router.connected-total":                {},
		"router.disconnected-total":             {},
		"router.response.test.status.200.count": {},
		"router.response.test.time":             {},
		"router.response.test.size":             {},
		"router.response.test.status":           {},
	}
	tracked := []string{}
	registry.Each(func(k string, _ interface{}) {
		tracked = append(tracked, k)
	})
	if len(tracked) != len(expected) {
		t.Error("unexpected size of the tracked list", tracked)
	}
	for _, k := range tracked {
		if _, ok := expected[k]; !ok {
			t.Error("the key", k, " has not been tracked")
		}
	}

	ts = httptest.NewServer(NewExpHandler(&registry))

	resp, err := http.Get(ts.URL)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("unexpected status code: %d\n", resp.StatusCode)
	}

	ts.Close()
}
