package service

import (
	"github.com/sniperkit/colly/plugins/data/aggregate/web/core/memcache"
)

func NewCache() *memcache.MemcacheClient {
	return memcache.NewMemcacheClient()
}
