package storage

import (
	"fmt"
	"sort"
	"strings"

	"github.com/sniperkit/colly/plugins/data/storage/backend"
)

// Initialize creates a new Store object, initializing the client
type Initialize func(addrs []string, options *backends.Config) (backends.Store, error)

var (
	// Backend initializers
	initializers = make(map[backends.Backend]Initialize)

	supportedBackend = func() string {
		keys := make([]string, 0, len(initializers))
		for k := range initializers {
			keys = append(keys, string(k))
		}
		sort.Strings(keys)
		return strings.Join(keys, ", ")
	}()
)

// NewStore creates an instance of store
func NewStore(backend backends.Backend, addrs []string, options *backends.Config) (backends.Store, error) {
	if init, exists := initializers[backend]; exists {
		return init(addrs, options)
	}

	return nil, fmt.Errorf("%s %s", backends.ErrBackendNotSupported.Error(), supportedBackend)
}

// AddStore adds a new store backend to storage
func AddStore(store backends.Backend, init Initialize) {
	initializers[store] = init
}
