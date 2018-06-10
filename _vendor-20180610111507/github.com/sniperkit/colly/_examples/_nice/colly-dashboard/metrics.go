package main

import (
	"fmt"
	"net/http"
	"runtime"
	"strings"
	"time"

	// stats - core
	"github.com/sniperkit/colly/plugins/data/collection/stats"

	// stats - collectors
	"github.com/sniperkit/colly/plugins/data/collection/stats/httpstats"

	// stats - remote clients
	"github.com/sniperkit/colly/plugins/data/collection/stats/datadog"
	"github.com/sniperkit/colly/plugins/data/collection/stats/influxdb"
	// "github.com/sniperkit/colly/plugins/data/collection/stats/prometheus"

	metric "github.com/sniperkit/colly/pkg/metric"
	tachymeter "github.com/sniperkit/colly/plugins/data/collection/stats/tachymeter"
)

// stats - tachymeter params
var (
	startedAt            time.Time
	wallTimeStart        time.Time
	cTachymeter          chan *tachymeter.Tachymeter
	xTachy               *tachymeter.Tachymeter
	xTachyResults        *tachymeter.Metrics
	xTachyTimeline       tachymeter.Timeline
	isTachymeterParallel bool
)

// stats - dashboard
var (
	collectorStats           *metric.Statistics
	collectorResponseMetrics chan metric.Response
)

// stats - httpstats
var (
	// core
	enableCollectorStatsTransport bool = true
	collectorStatsEngine          *stats.Engine
	collectorStatsTags            []*stats.Tag
	// influxDB
	collectorStatsClientInfluxDB     *influxdb.Client
	collectorStatsClientInfluxDBConf *influxdb.ClientConfig
	// datadog
	collectorStatsClientDatadog     *datadog.Client
	collectorStatsClientDatadogConf *datadog.ClientConfig
)

func initTachymeter() {

	// Create a Tachymeter
	cTachymeter = make(chan *tachymeter.Tachymeter)
	xTachyTimeline = tachymeter.Timeline{}
	xTachy = tachymeter.New(
		&tachymeter.Config{
			// Tachymeter
			SampleSize: appConfig.Debug.Tachymeter.SampleSize,
			HBins:      appConfig.Debug.Tachymeter.HistogramBins,
			Export: &tachymeter.Export{
				// Exports
				Encoding:   appConfig.Debug.Tachymeter.Export.Encoding,
				Basename:   appConfig.Debug.Tachymeter.Export.Basename,
				PrefixPath: appConfig.Debug.Tachymeter.Export.PrefixPath,
				SplitLimit: appConfig.Debug.Tachymeter.Export.SplitAt,
				BufferSize: appConfig.Debug.Tachymeter.Export.BufferSize,
				Overwrite:  appConfig.Debug.Tachymeter.Export.Overwrite,
				BackupMode: appConfig.Debug.Tachymeter.Export.BackupMode,
			},
		},
	)

	switch appConfig.Collector.CurrentMode {
	case "async":
		fallthrough
	case "queue":
		appConfig.Debug.Tachymeter.Async = true
		wallTimeStart = time.Now()
	}

}

func newTachymeterWithConfig(cfg *tachymeter.Config) *tachymeter.Tachymeter {
	// Create a new Tachymeter instance
	newTachymeter := tachymeter.New(
		&tachymeter.Config{
			// Tachymeter
			SafeMode:   true, // deprecated
			SampleSize: 50,
			HBins:      10,
			Export: &tachymeter.Export{
				// Exports
				Encoding:   "tsv",
				Basename:   "golanglibs_tachymter_%d",
				PrefixPath: "./shared/exports/stats/tachymeter/",
				SplitLimit: 2500,
				BufferSize: 20000,
				Overwrite:  true,
				BackupMode: true,
			},
		},
	)

	switch appConfig.Collector.CurrentMode {
	case "async":
		fallthrough
	case "queue":
		appConfig.Debug.Tachymeter.Async = true
		wallTimeStart = time.Now()
	}

	return newTachymeter
}

func closeTachymeter(outputFile string) error {
	// Write out an HTML page with the histogram for all iterations.
	err := xTachyTimeline.WriteHTML(outputFile)
	if err != nil {
		return err
	}
	return nil
}

func addHttpStatsTransport(rt http.RoundTripper) http.RoundTripper {
	return httpstats.NewTransport(rt)
}

func newStatsEngine(backend string) {
	switch backend {
	case "datadog":
		collectorStatsEngine = nil
	case "influxdb":
		fallthrough
	default:
		collectorStatsEngine = nil
	}
}

type funcMetrics struct {
	calls struct {
		count  int           `metric:"count" type:"counter"`
		failed int           `metric:"failed" type:"counter"`
		time   time.Duration `metric:"time"  type:"histogram"`
	} `metric:"func.calls"`
}

func funcTrack(start time.Time) {
	return
	function, file, line, _ := runtime.Caller(1)
	go func() {
		elapsed := time.Since(start)
		if log != nil {
			log.Printf("main().funcTrack() %s took %s", fmt.Sprintf("%s:%s:%d", runtime.FuncForPC(function).Name(), chopPath(file), line), elapsed)
		}

	}()
}

// return the source filename after the last slash
func chopPath(original string) string {
	i := strings.LastIndex(original, "/")
	if i == -1 {
		return original
	} else {
		return original[i+1:]
	}
}

func addMetrics(start time.Time, incr int, failed bool) {
	callTime := time.Now().Sub(start)
	m := &funcMetrics{}
	m.calls.count = incr
	m.calls.time = callTime
	if failed {
		m.calls.failed = incr
	}
	stats.Report(m)
}

func statsWithTags() {}

/*
	if len(config.Dogstatsd.Address) != 0 {
		stats.Register(datadog.NewClientWith(datadog.ClientConfig{
			Address:    config.Dogstatsd.Address,
			BufferSize: config.Dogstatsd.BufferSize,
		}))
	}
*/

/*
func newCacheTransportWithStats(engine string, prefixPath string) (httpcache.Cache, *httpcache.Transport) {
	defer funcTrack(time.Now())

	backendCache, err := newCacheBackend(engine, prefixPath)
	if err != nil {
		log.Fatal("cache err", err.Error())
	}

	var httpTransport = http.DefaultTransport
	httpTransport = httpstats.NewTransport(httpTransport)
	http.DefaultTransport = httpTransport

	cachingTransport := httpcache.NewTransportFrom(backendCache, httpTransport) // httpcache.NewMemoryCacheTransport()
	cachingTransport.MarkCachedResponses = true

	return backendCache, cachingTransport
}
*/
