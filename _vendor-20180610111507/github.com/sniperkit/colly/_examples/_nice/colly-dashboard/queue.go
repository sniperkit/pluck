package main

import (
	"errors"
	"strings"

	//// collector - queue core
	queue "github.com/sniperkit/colly/pkg/queue"
	storage "github.com/sniperkit/colly/pkg/storage"

	//// collector - queue storage backends
	res "github.com/sniperkit/colly/plugins/data/storage/backend/redis"
	// lru "github.com/sniperkit/colly/plugins/data/storage/backend/lru"
	// sq3 "github.com/sniperkit/colly/plugins/data/storage/backend/sqlite3"
	// baq "github.com/sniperkit/colly/plugins/data/storage/backend/badger"
	// stq "github.com/sniperkit/colly/plugins/data/storage/backend/storm"
	// myq "github.com/sniperkit/colly/plugins/data/storage/backend/mysql"
	// moq "github.com/sniperkit/colly/plugins/data/storage/backend/mongo"
	// elq "github.com/sniperkit/colly/plugins/data/storage/backend/elastic"
	// shq "github.com/sniperkit/colly/plugins/data/storage/backend/sphinx"
	// caq "github.com/sniperkit/colly/plugins/data/storage/backend/cassandra"
)

// collector queue
var (
	collectorQueueStorage storage.Storage // Storage interface
	collectorQueue        *queue.Queue    // collector's queue instance
)

// colly queue processing errors
var (
	errInvalidQueueThreads     = errors.New("Invalid queue consumer threads count. Must be superior or equal to 0.")
	errInvalidQueueBackend     = errors.New("Unkown queue storage backend name. Available: inmemory, redis, sqlite3, badger, mysql, postgres.")
	errInvalidQueueMaxSize     = errors.New("Invalid queue max size value. Must be superior or equal to 0.")
	errLocalFileStat           = errors.New("File not found.")
	errLocalFileOpen           = errors.New("Could not open the filepath")
	errInvalidRemoteStatusCode = errors.New("errInvalidRemoteStatusCode")
)

func initCollectorQueue(queueThreads int, maxSize int, storeBackend string) (q *queue.Queue, err error) {
	if queueThreads < 0 {
		err = errInvalidQueueThreads
		return
	}
	storeBackend = strings.ToLower(storeBackend)
	if maxSize < 0 {
		err = errInvalidQueueMaxSize
		return
	}

	switch storeBackend {
	// case "sqlite":
	//	fallthrough

	// case "sqlite3": // Warning! Conflict with Pivot
	//	q, err = queue.New(queueThreads, &sq3.Storage{Filename: "./shared/datastore/queue.db"})

	case "redis":
		q, err = queue.New(queueThreads, &res.Storage{Address: "127.0.0.1:6379", Password: "", DB: 0, Prefix: "job01"})

	case "badger":
		log.Printf("Collector queue storage/backend '%s' is not implemented yet...\n", storeBackend)

	case "diskv":
		log.Printf("Collector queue storage/backend '%s' is not implemented yet...\n", storeBackend)

	case "boltdb":
		log.Printf("Collector queue storage/backend '%s' is not implemented yet...\n", storeBackend)

	case "lru":
		log.Printf("Collector queue storage/backend '%s' is not implemented yet...\n", storeBackend)

	case "inmemory":
		fallthrough

	default:
		q, err = queue.New(queueThreads, &queue.InMemoryQueueStorage{MaxSize: maxSize})

	}
	return
}
