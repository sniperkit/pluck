package main

import (
	// collector - plugins/addons
	cuckoo "github.com/seiflotfy/cuckoofilter"                         // Important! To fork this package
	cmap "github.com/sniperkit/colly/plugins/data/structure/map/multi" // Concurrent multi-map helper
)

var (
	// concurrent multi map store
	collectorConcurrentMap             *cmap.ConcurrentMap             // concurrent map
	collectorShardedConcurrentMap      *cmap.ShardedConcurrentMap      // concurrent map with shards. (Default: 16)
	collectorConcurrentMultiMap        *cmap.ConcurrentMultiMap        // concurrent multi-map
	collectorShardedConcurrentMultiMap *cmap.ShardedConcurrentMultiMap // concurrent multi-map with shards
	collectorConcurrentCheckInsert     bool                            = false

	// cuckofilter store
	collectorCuckooFilterCapacity uint = 20000 // default: 1000000
	collectorCuckooFilter         *cuckoo.CuckooFilter
)

func AutoLoad() {
	// Concurrent map
	if collectorConcurrentMap == nil {
		collectorConcurrentMap, _ = newConcurrentMap()
	}

	// Sharded concurrent map. (Default: 16)
	if collectorShardedConcurrentMap == nil {
		collectorShardedConcurrentMap, _ = newShardedConcurrentMap(16)
	}

	// Concurrent multi-map
	if collectorConcurrentMultiMap == nil {
		collectorConcurrentMultiMap, _ = newConcurrentMultiMap()
	}

	// Concurrent multi-map with shards
	if collectorShardedConcurrentMultiMap == nil {
		collectorShardedConcurrentMultiMap, _ = newShardedConcurrentMultiMap(16)
	}

	// Concurrent cuckoofilter map
	if collectorCuckooFilter == nil {
		collectorCuckooFilter, _ = newCuckoofilter(collectorCuckooFilterCapacity)
	}
}

func newConcurrentMap() (*cmap.ConcurrentMap, bool) {
	c := cmap.NewConcurrentMap()
	return c, c != nil
}

func newShardedConcurrentMap(shards uint32) (*cmap.ShardedConcurrentMap, bool) {
	if shards <= 0 {
		shards = 16
	}
	c := cmap.NewShardedConcurrentMap(cmap.WithNumberOfShards(shards))
	return c, c != nil
}

func newConcurrentMultiMap() (*cmap.ConcurrentMultiMap, bool) {
	c := cmap.NewConcurrentMultiMap()
	return c, c != nil
}

func newShardedConcurrentMultiMap(shards uint32) (*cmap.ShardedConcurrentMultiMap, bool) {
	c := cmap.NewShardedConcurrentMultiMap()
	return c, c != nil
}

func newCuckoofilter(capacity uint) (*cuckoo.CuckooFilter, bool) {
	c := cuckoo.NewCuckooFilter(capacity)
	return c, c != nil
}

/*
var (
		// collector store mapping with cli arguments
		collectorStoreListers = map[string]func(ctx context.Context, args []string){
			"cmap":         listCollectorConcurrentMap,
			"cmmap":        listCollectorConcurrentMultiMap,
			"cuckoofilter": listCollectorCuckooFilter,
		}
)

// just for test purpose
// func listCollectorConcurrentMap()      {}
// func listCollectorConcurrentMultiMap() {}
// func listlistCollectorCuckooFilter()   {}
*/
