package main

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	// colly - core
	colly "github.com/sniperkit/colly/pkg"

	// cache - core
	"github.com/gregjones/httpcache"

	// cache - advanced backends
	"github.com/sniperkit/colly/plugins/data/storage/backend/default/badger"
	"github.com/sniperkit/colly/plugins/data/storage/backend/default/diskv"

	// general helpers
	"github.com/sniperkit/colly/plugins/system/fs"
)

// cache related objects
var (
	collector_http_client          *http.Client
	collector_http_round_tripper   http.RoundTripper
	collector_http_cache_transport *httpcache.Transport
	collector_http_cache_backend   httpcache.Cache
)

func new_collector_with_http_cache(c *colly.Collector) *colly.Collector {
	_, t, err := new_http_cache_with_transport("badger", "./shared/storage/cache/http")
	if err != nil {
		log.Warnln("cache err", err.Error())
		return c
	}
	c.WithTransport(t)
	return c
}

func new_http_cache(cacheStorageName, cacheStoragePath string) (httpcache.Cache, error) {
	cacheStoreHTTP, err := new_cache_store(cacheStorageName, cacheStoragePath)
	if err != nil {
		// log.Fatal("cache err", err.Error())
		return nil, err
	}
	return cacheStoreHTTP, nil
}

func new_http_cache_with_transport(engine, prefixPath string) (httpcache.Cache, *httpcache.Transport, error) {
	cacheBackend, err := new_cache_store(engine, prefixPath)
	if err != nil {
		return nil, nil, err
	}
	cacheTransport := httpcache.NewTransport(cacheBackend)
	cacheTransport.MarkCachedResponses = true
	return cacheBackend, cacheTransport, nil
}

func new_http_cache_round_tripper(c httpcache.Cache, markCachedResponses bool) http.RoundTripper {
	t := httpcache.NewTransport(c)
	t.MarkCachedResponses = markCachedResponses
	return t
}

func new_cache_store(engine, prefixPath string) (backend httpcache.Cache, err error) {
	fsutil.EnsureDir(prefixPath)

	engine = strings.ToLower(engine)
	switch engine {
	case "diskv":
		appConfig.Collector.Transport.Http.Cache.Store.Directory = filepath.Join(prefixPath, "httpcache.diskv")
		fsutil.EnsureDir(appConfig.Collector.Transport.Http.Cache.Store.Directory)
		backend = diskcache.New(appConfig.Collector.Transport.Http.Cache.Store.Directory)

	case "badger":
		appConfig.Collector.Transport.Http.Cache.Store.Directory = filepath.Join(prefixPath, "httpcache.badger")
		fsutil.EnsureDir(appConfig.Collector.Transport.Http.Cache.Store.Directory)
		backend, err = badgercache.New(
			&badgercache.Config{
				ValueDir:    "golanglibs.com",
				StoragePath: appConfig.Collector.Transport.Http.Cache.Store.Directory,
				SyncWrites:  false,
				Debug:       false,
				Compress:    true,
				TTL:         time.Duration(120 * 24 * time.Hour),
			},
		)

	case "memory":
		fallthrough

	default:
		backend = httpcache.NewMemoryCache()

	}
	return
}

func set_cache(key string, obj map[string]interface{}) {
	collector_http_cache_backend.Set(key, toBytes(mapToString(obj)))
}

func toBytes(input string) []byte {
	return []byte(input)
}

func mapToString(input map[string]interface{}) string {
	return toString(input)
}

func toString(obj interface{}) string {
	return fmt.Sprintf("%v", obj)
}

/*
func loadSkiplist(filepath string) {
	fp, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer fp.Close()

	csv := csv.NewReader(fp)
	lines := streamCsv(csv, csvReaderBuffer)
	for line := range lines {
		expiresAt, err := parseTimeStamp(line.GetByName("task_expired_timestamp"))
		if err != nil {
			log.Errorln("[SKIP-ERROR] taskInfo, service=", line.GetByName("service"), "topic=", line.GetByName("topic"), "expiresTimestamp", line.GetByName("task_expired_timestamp"))
			continue
		}
		now := time.Now()
		if now.After(expiresAt.Add(cacheTTL)) {
			log.Infoln("[TSK-ALLOW] task info, service=", line.GetByName("service"), "topic=", line.GetByName("topic"), "expiresAt=", expiresAt)
			continue
		}
		cuckflt.InsertUnique([]byte(line.GetByName("topic")))
	}
	log.Warnln("[TSK-EXCLUDED] taskInfo, count=", cuckflt.Count())
}
*/
