package colly

import (
	"sync"

	"github.com/sniperkit/colly/pkg/storage"
)

type assetSerializer struct {
	store storage.Storage
	lock  *sync.RWMutex
}

type bodySerializer struct {
	store storage.Storage
	lock  *sync.RWMutex
}

type responseSerializer struct {
	store storage.Storage
	lock  *sync.RWMutex
}

type cookieJarSerializer struct {
	store storage.Storage
	lock  *sync.RWMutex
}
