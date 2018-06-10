package badgercache

import (
	"encoding/hex"
	"errors"
	"path/filepath"
	"time"

	sync "github.com/sniperkit/xutil/plugin/concurrency/sync/debug"
	// "github.com/sniperkit/xutil/plugin/debug/pp"

	"github.com/dgraph-io/badger"
	"github.com/json-iterator/go"
)

const (
	defaultCacheTTL      time.Duration = time.Duration(30 * 7 * time.Hour)
	defaultCacheValueDir string        = "httpcache"
	defaultCacheDir      string        = "./shared/data/cache/.badger"
)

var defaultCachePath string = filepath.Join(defaultCacheDir, defaultCacheValueDir)

/*
	Refs:
	- github.com/rohanthewiz/robadger
	- https://github.com/rohanthewiz/robadger/blob/master/badger.go
*/

// Cache stores and retrieves data using Badger KV.
type Cache struct {
	mu          sync.RWMutex
	db          *badger.DB
	storagePath string
	bucketName  string
	compress    bool
	debug       bool
	hasTTL      bool
	ttl         time.Duration
}

/*
	refs:
	- https://github.com/dilyevsky/httplru/blob/master/cache/badger_lru.go
	- https://github.com/jonas747/dbstate/blob/master/util.go
	- https://github.com/jonas747/dbstate
	- https://github.com/vicanso/pike/blob/master/cache/cache.go (fasthttp)
	-
*/
/*
// badgerCache implements LRU cache using BadgerDB as high-performant, embedded k/v store.
type badgerCache struct {
	maxEntries int
	ttl        time.Duration

	db *badger.DB
	mu sync.Mutex
	// TODO(dilyevsky): Perhaps LSM or B+ tree would be of use, however
	// as it is, it can already support millions of cached entries pretty
	// efficiently.
	// Anything larger would probably need to be distributed anyway.
	ll    *list.List
	cache map[uint64]*list.Element
}
*/

type Check struct {
	Enabled   bool
	Key       string
	Requests  int
	CreatedAt time.Time
	UpdatedAt time.Time
	ExpiredAt time.Time
	Priority  bool
	Provider  string
}

type Config struct {
	Debug bool

	UseTTL bool
	TTL    time.Duration

	Compress bool

	// 1. Mandatory flags
	// -------------------
	// Directory to store the data in. Should exist and be writable.
	StoragePath string // Dir

	// Directory to store the value log in. Can be the same as Dir. Should
	// exist and be writable.
	ValueDir string

	// 2. Frequently modified flags
	// -----------------------------
	// Sync all writes to disk. Setting this to true would slow down data
	// loading significantly.
	SyncWrites bool

	// 3. Flags that user might want to review
	// ----------------------------------------
	// The following affect all levels of LSM tree.
	MaxTableSize        int64 // Each table (or file) is at most this size.
	LevelSizeMultiplier int   // Equals SizeOf(Li+1)/SizeOf(Li).
	MaxLevels           int   // Maximum number of levels of compaction.

	// If value size >= this threshold, only store value offsets in tree.
	ValueThreshold int

	// Maximum number of tables to keep in memory, before stalling.
	NumMemtables int

	// The following affect how we handle LSM tree L0.
	// Maximum number of Level 0 tables before we start compacting.
	NumLevelZeroTables int

	// If we hit this number of Level 0 tables, we will stall until L0 is
	// compacted away.
	NumLevelZeroTablesStall int

	// Maximum total size for L1.
	LevelOneSize int64

	// Size of single value log file.
	ValueLogFileSize int64

	// Number of compaction workers to run concurrently.
	NumCompactors int

	// 4. Flags for testing purposes
	// ------------------------------
	DoNotCompact bool // Stops LSM tree from compactions.

}

func Mount(client *badger.DB) *Cache {
	return &Cache{db: client}
}

func New(config *Config) (*Cache, error) {

	if config.Debug {
		LogWithFields(Fields{
			"config": config,
		}).Warnf("badgercache.New()")
	}

	badgerConfig := badger.DefaultOptions
	if config == nil {
		badgerConfig.Dir = defaultCacheDir
		badgerConfig.ValueDir = defaultCacheValueDir
		badgerConfig.SyncWrites = true
	} else {
		badgerConfig.Dir = config.StoragePath
		badgerConfig.ValueDir = filepath.Join(config.StoragePath, config.ValueDir)
		badgerConfig.SyncWrites = config.SyncWrites

		if config.TTL > 0 {
			config.UseTTL = true
		}
		if config.UseTTL && config.TTL <= 0 {
			config.TTL = defaultCacheTTL
			// warning !
		}

	}

	// pp.Println("config=", config)
	// log.Fatalln("config.UseTTL=", config.UseTTL, "config.TTL=", config.TTL)

	client, err := badger.Open(badgerConfig)
	if err != nil {
		LogWithFields(Fields{
			"config": config,
		}).Fatalln("badgercache.New().badger.Open(), ERROR: ", err)
		return nil, err
	}

	return &Cache{
		db:       client,
		ttl:      config.TTL,
		hasTTL:   config.UseTTL,
		debug:    config.Debug,
		compress: config.Compress,
	}, nil
}

func (c *Cache) Get(key string) (resp []byte, ok bool) {
	c.mu.Lock()
	start := time.Now()

	err := c.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			if c.debug {
				LogWithFields(Fields{
					"duration": time.Since(start),
					"key":      key,
				}).Warnln("badgercache.Get().View(), ERROR: ", err)
			}
			return err
		}
		resp, err = item.Value()
		if err != nil {
			if c.debug {
				LogWithFields(Fields{
					"duration": time.Since(start),
					"key":      key,
				}).Warnln("badgercache.Get(), ERROR: ", err)
			}
			return err
		}
		if c.compress {
			var err error
			resp, err = Decompress(resp)
			if err != nil {
				LogWithFields(Fields{
					"duration": time.Since(start),
					"resp":     string(resp),
					"key":      key,
					"err":      err.Error(),
				}).Fatalln("badgercache.Get().Decompress(), ERROR: ", err)
				return err
			}
		}
		return nil
	})

	if c.debug {
		LogWithFields(Fields{
			"duration": time.Since(start),
			"key":      key,
			"ok":       err == nil,
		}).Info("badgercache.Get()")
	}
	c.mu.Unlock()
	return resp, err == nil
}

func (c *Cache) prettyView(key string, resp []byte, ttl time.Duration) error {
	return c.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.AllVersions = true
		iter := txn.NewIterator(badger.DefaultIteratorOptions)
		count := 1
		for iter.Rewind(); iter.Valid(); iter.Next() {
			item := iter.Item()
			val, err := item.Value()
			if err != nil {
				return err
			}
			LogWithFields(Fields{
				"Key":      item.Key(),
				"Value":    hex.EncodeToString(val[:10]),
				"UserMeta": item.UserMeta(),
				"Version":  item.Version(),
			}).Println("badgercache.prettyView().All(), count=", count)
			count++
		}
		return nil
	})
}

func (c *Cache) setWithTTL(key string, resp []byte, ttl time.Duration) error {
	c.mu.Lock()
	start := time.Now()

	err := c.db.Update(func(txn *badger.Txn) error {
		if c.compress {
			var err error
			resp, err = Compress(resp)
			if err != nil {
				LogWithFields(Fields{
					"duration": time.Since(start),
					"resp":     string(resp),
					"key":      key,
					"err":      err.Error(),
				}).Fatalln("badgercache.Get().Compress(), ERROR: ", err)
				return errors.New("error while compressing content...")
			}
		}

		err := txn.SetWithTTL([]byte(key), resp, ttl)
		if err != nil {
			LogWithFields(Fields{
				"duration": time.Since(start),
				"key":      key,
			}).Fatalln("badgercache.Set(), ERROR: ", err)
		}

		return err
	})

	if c.debug {
		LogWithFields(Fields{
			"duration": time.Since(start),
			"key":      key,
			"ok":       err == nil,
		}).Warnln("badgercache.Set()")
	}
	c.mu.Unlock()
	return err
}

// Set stores a response to the cache at the given key.
func (c *Cache) Set(key string, resp []byte) {
	c.mu.Lock()
	start := time.Now()

	err := c.db.Update(func(txn *badger.Txn) error {
		if c.compress {
			var err error
			resp, err = Compress(resp)
			if err != nil {
				LogWithFields(Fields{
					"duration": time.Since(start),
					"resp":     string(resp),
					"key":      key,
					"err":      err.Error(),
				}).Fatalln("badgercache.Get().Compress(), ERROR: ", err)
				return errors.New("error while compressing content...")
			}
		}
		if c.debug {
			log.Debugln("cache.badger.SetWithTTL()=", c.hasTTL, ", cache.ttl=", c.ttl, "key=", key)
		}
		var err error
		if c.hasTTL {
			err = txn.SetWithTTL([]byte(key), resp, c.ttl)
		} else {
			err = txn.Set([]byte(key), resp)
		}

		if err != nil {
			LogWithFields(Fields{
				"duration": time.Since(start),
				"key":      key,
			}).Fatalln("badgercache.Set(), ERROR: ", err)
		}
		return err
	})
	if c.debug {
		LogWithFields(Fields{
			"duration": time.Since(start),
			"key":      key,
			"ok":       err == nil,
		}).Warnln("badgercache.Set()")
	}
	c.mu.Unlock()
	return
}

func (c *Cache) Delete(key string) {
	c.mu.Lock()
	start := time.Now()
	err := c.db.Update(func(txn *badger.Txn) error {
		err := txn.Delete([]byte(key))
		if err != nil {
			LogWithFields(Fields{
				"key":      key,
				"duration": time.Since(start),
			}).Fatalln("badgercache.Delete(), ERROR: ", err)
		}
		return err
	})
	if c.debug {
		LogWithFields(Fields{
			"duration": time.Since(start),
			"key":      key,
			"ok":       err == nil,
		}).Warnln("badgercache.Delete()")
	}
	c.mu.Unlock()
	return
}

func (c *Cache) Action(name string, args ...interface{}) (map[string]*interface{}, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if name == "getKeys" {
		keys := c.keys()
		resp := make(map[string]*interface{})
		for _, v := range keys {
			resp[v] = nil
		}
		return resp, nil
	}

	return nil, errors.New("Action not implemented yet")
}

func (c *Cache) Debug(action string) {
	LogWithFields(Fields{
		"action": action,
	}).Warnln("badgercache.ListAll()")
	// c.mu.Lock()
	// c.listAll()
	// c.keys()
	c.compressor()

	log.Fatal("badgercache.ListAll()")

	// c.db.PurgeOlderVersions()
	// c.purgeOlderVersions()
	// c.updates()
	// c.seekPrefix()
	// c.mu.Unlock()
}

// ListAll lists all the pairs KV of a given type
func (c *Cache) compressor() (err error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var key string
	in := []byte(`HTTP/1.1 200 OK\r\nTransfer-Encoding: chunked\r\nAccess-Control-Allow-Origin: *\r\nAccess-Control-Expose-Headers: ETag, Link, Retry-After, X-GitHub-OTP, X-RateLimit-Limit, X-RateLimit-Remaining, X-RateLimit-Reset, X-OAuth-Scopes, X-Accepted-OAuth-Scopes, X-Poll-Interval\r\nCache-Control: private, max-age=60, s-maxage=60\r\nContent-Security-Policy: default-src 'none'\r\nContent-Type: application/json; charset=utf-8\r\nDate: Sun, 17 Dec 2017 09:01:26 GMT\r\nEtag: W/\"97327ebe7afdfe040e366a8f35cf8c72\"\r\nLast-Modified: Fri, 26 May 2017 20:57:47 GMT\r\nServer: GitHub.com\r\nStatus: 200 OK\r\nStrict-Transport-Security: max-age=31536000; includeSubdomains; preload\r\nVary: Accept, Authorization, Cookie, X-GitHub-OTP\r\nX-Accepted-Oauth-Scopes: \r\nX-Content-Type-Options: nosniff\r\nX-Frame-Options: deny\r\nX-Github-Media-Type: github.v3; format=json\r\nX-Github-Request-Id: EA31:2902B:1281CE8:23CD027:5A363265\r\nX-Oauth-Scopes: repo, user\r\nX-Ratelimit-Limit: 5000\r\nX-Ratelimit-Remaining: 2441\r\nX-Ratelimit-Reset: 1513504311\r\nX-Runtime-Rack: 0.049262\r\nX-Varied-Accept: application/vnd.github.v3+json\r\nX-Varied-Authorization: Bearer 63814c0ef8a9a7a273e828d1cc4d410b4f449a9f\r\nX-Xss-Protection: 1; mode=block\r\n\r\n347\r\n{\"name\":\"README.md\",\"path\":\"README.md\",\"sha\":\"723b022f9b66c690b95239ef7de83f9dc9d24290\",\"size\":60,\"url\":\"https://api.github.com/repos/AaronTL/TapNews/contents/README.md?ref=master\",\"html_url\":\"https://github.com/AaronTL/TapNews/blob/master/README.md\",\"git_url\":\"https://api.github.com/repos/AaronTL/TapNews/git/blobs/723b022f9b66c690b95239ef7de83f9dc9d24290\",\"download_url\":\"https://raw.githubusercontent.com/AaronTL/TapNews/master/README.md\",\"type\":\"file\",\"content\":\"IyBUYXBOZXdzClJlYWwgVGltZSBOZXdzIFNjcmFwaW5nIGFuZCBSZWNvbW1l\\nbmRhdGlvbiBTeXN0ZW0K\\n\",\"encoding\":\"base64\",\"_links\":{\"self\":\"https://api.github.com/repos/AaronTL/TapNews/contents/README.md?ref=master\",\"git\":\"https://api.github.com/repos/AaronTL/TapNews/git/blobs/723b022f9b66c690b95239ef7de83f9dc9d24290\",\"html\":\"https://github.com/AaronTL/TapNews/blob/master/README.md\"}}\r\n0\r\n\r\n`)

	key = "test_raw"
	err = c.db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(key), in)
	})
	LogWithFields(Fields{
		"key": key,
		// "in":        string(in),
		// "contains": Match("application/json", in),
		"isJSON? ": SubString("application/json", in),
		"isGZIP? ": IsGZIP(in),
	}).Warnln("badgercache.compressor()")
	Match("application/json", in)

	key = "test_gzipped"
	err = c.db.Update(func(txn *badger.Txn) error {
		in, err = gzipData(in)
		if err != nil {
			return err // errors.New("error while compressing content...")
		}
		return txn.Set([]byte(key), in)
	})
	LogWithFields(Fields{
		"key": key,
		// "in":        string(in),
		// "contains":  Match("application/json", in),
		"isJSON? ": SubString("application/json", in),
		"isGZIP? ": IsGZIP(in),
	}).Warnln("badgercache.compressor()")

	return err
}

// ListAll lists all the pairs KV of a given type
func (c *Cache) listAll() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	err := c.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		// opts.PrefetchSize = 10000
		it := txn.NewIterator(opts)
		i := 0
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			k := item.Key()
			kStr := string(k)
			v, err := item.Value()
			if err != nil {
				LogWithFields(Fields{"key": kStr}).Fatalln("badgercache.ListAll(), ERROR: ", err)
				return err
			}
			vStr := string(v)
			/*
				vStr, err := strconv.Unquote(string(v))
				if err != nil {
					log.WithFields(log.Fields{
						"v": string(v),
					}).Info("badgercache.ListAll().strconv.Unquote()")
					log.WithFields(log.Fields{"key": kStr}).Fatalln("badgercache.ListAll().strconv.Unquote(), ERROR: ", err)
					return err
				}
			*/
			ling := DetectLang(vStr)
			lang, safe := DetectType(string(k), vStr)
			LogWithFields(Fields{
				"val": vStr,
			}).Info("badgercache.ListAll().Value()")
			LogWithFields(Fields{
				"iter":     i,
				"key":      kStr,
				"safe":     safe,
				"lang":     lang,
				"ling":     ling,
				"IsYAML? ": IsYAML([]byte(vStr)),
				"IsJSON? ": IsJSON([]byte(vStr)),
				"IsGZIP? ": IsGZIP([]byte(vStr)),
			}).Warn("badgercache.ListAll().Detections()")
			i++
		}
		return nil
	})
	if c.debug {
		LogWithFields(Fields{
			"ok": err == nil,
		}).Warnln("badgercache.ListAll()")
	}

	return err
}

func (c *Cache) addJSON(key string, payload interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	err := c.db.Update(func(txn *badger.Txn) error {
		payload, err := jsoniter.Marshal(payload)
		if err != nil {
			LogWithFields(Fields{
				"key": key,
			}).Errorln("badgercache.AddJSON(), ERROR: ", err)
			return err
		}
		if c.compress {
			payload, err = gzipData(payload)
			if err != nil {
				return err // errors.New("error while compressing content...")
			}
		}
		/*
			err := txn.Set([]byte(key), payload)
			if c.debug && err != nil {
				log.WithFields(log.Fields{
					"key": key,
				}).Warnln("badgercache.AddJSON(), ERROR: ", err)
			}
		*/
		return txn.Set([]byte(key), payload)
	})

	if c.debug {
		LogWithFields(Fields{
			"key": key,
			"ok":  err == nil,
		}).Warnln("badgercache.Set()")
	}
	return err
}

func (c *Cache) keys() (keys []string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	err := c.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			k := item.Key()
			keys = append(keys, string(k))
			// log.WithFields(log.Fields{"key": string(k)}).Info("badgercache.Keys()")
		}
		return nil
	})
	if c.debug {
		LogWithFields(Fields{
			"ok": err == nil,
		}).Warnln("badgercache.ListAll()")
	}
	return // keys
}

/*
func (c *Cache) purgeOlderVersions() {
	c.mu.Lock()
	defer c.mu.Unlock()

	log.Info("badgercache.purgeOlderVersions(), START")
	c.db.PurgeOlderVersions()
	log.Info("badgercache.purgeOlderVersions(), END")
}
*/

func (c *Cache) seekPrefix(value string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		it := txn.NewIterator(opts)
		prefix := []byte(value)
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			k := item.Key()
			v, err := item.Value()
			if err != nil {
				LogWithFields(Fields{"key": k, "val": v}).Errorln("badgercache.SeekPrefix(), ERROR: ", err)
				return err
			}
			LogWithFields(Fields{"key": k, "val": v}).Info("badgercache.SeekPrefix()")
		}
		return nil
	})
}

/*
func (c *Cache) firstIndex() (uint64, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	tx := c.db.NewTransaction(false)
	defer tx.Discard()
	iter := tx.NewIterator(iterAscOpt)
	iter.Rewind()
	item := iter.Item()
	if item == nil {
		return 0, nil
	}

	return bytesToUint64(item.Key()), nil
}

func (c *Cache) lastIndex() (uint64, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	tx := c.db.NewTransaction(false)
	defer tx.Discard()
	iter := tx.NewIterator(iterDescOpt)
	iter.Rewind()
	item := iter.Item()
	if item == nil {
		return 0, nil
	}
	return bytesToUint64(item.Key()), nil
}

func (c *Cache) deleteRange(min, max uint64) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	tx := c.db.NewTransaction(true)
	defer tx.Discard()
	minKey := uint64ToBytes(min)
	iter := tx.NewIterator(iterAscOpt)
	for iter.Seek(minKey); iter.Valid(); iter.Next() {
		item := iter.Item()
		if item == nil {
			break
		}
		curKey := safeKey(item)
		if bytesToUint64(curKey) > max {
			break
		}
		if err := tx.Delete(curKey); err != nil {
			return err
		}
	}
	if err := tx.Commit(nil); err != nil {
		return err
	}
	return nil
}
*/

/*
func (c *Cache) updates(updates map[string]interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	txn := db.NewTransaction(true)
	for k,v := range updates {
	  if err := txn.Set(byte[](k),byte[](v)); err == ErrTxnTooBig {
	    _ = txn.Commit()
	    txn = db.NewTransaction(..)
	    _ = txn.Set(k,v)
	  }
	}
	_ = txn.Commit()

}
*/

// Close closes the underlying boltdb database.
func (c *Cache) Close() error {
	return c.db.Close()
}
