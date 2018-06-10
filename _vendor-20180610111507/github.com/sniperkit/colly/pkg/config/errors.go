package config

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
)

// cache storage errors
var (
	ErrInvalidCacheDuration = errors.New("Invalid cache duration...")
	ErrInvalidCacheTTL      = errors.New("Invalid cache TTL...")
	ErrInvalidCacheTimeUnit = errors.New("Invalid cache time unit. available: second, hour, day")
	// ErrInvalidCacheBackend  = errors.New("Invalid cache backend. Available: " + strings.Join(DefaultCacheBackends, ",") + ".")
)

// colly queue processing errors
var (
	ErrInvalidQueueThreads     = errors.New("Invalid queue consumer threads count. Must be superior or equal to 0.")
	ErrInvalidQueueBackend     = errors.New("Unkown queue storage backend name. Available: inmemory, redis, sqlite3, badger, mysql, postgres.")
	ErrInvalidQueueMaxSize     = errors.New("Invalid queue max size value. Must be superior or equal to 0.")
	ErrLocalFileStat           = errors.New("File not found.")
	ErrLocalFileOpen           = errors.New("Could not open the filepath")
	ErrInvalidRemoteStatusCode = errors.New("errInvalidRemoteStatusCode")
)

// e returns an error, prefixed with the name of the function that triggered it. Originally by StackOverflow user svenwltr:
// http://stackoverflow.com/a/38551362/199475
func e(err error) error {
	pc, _, _, _ := runtime.Caller(2)

	fr := runtime.CallersFrames([]uintptr{pc})
	namer, _ := fr.Next()
	name := namer.Function

	if !fullyQualifiedPath {
		fn := strings.Split(name, "/")
		if len(fn) > 0 {
			return fmt.Errorf("%s: %s", fn[len(fn)-1], err.Error())
		}
	}

	return fmt.Errorf("%s: %s", name, err.Error())
}

// Err consumes an error, a string, or nil, and produces an error message prefixed with the name of the function that called it (or nil).
func Err(err interface{}) error {
	switch o := err.(type) {
	case string:
		return e(fmt.Errorf("%s", o))
	case error:
		return e(o)
	default:
		return nil
	}
}
