// Package cuckoofilter implements cuckoo filters from the paper "Cuckoo Filter: Practically Better Than Bloom" by Fan et al.
// https://www.cs.cmu.edu/~dga/papers/cuckoo-conext2014.pdf
package cuckoofilters

import (
	"math/rand"
	"sync"

	"github.com/sniperkit/colly/plugins/data/counter/ordered"
)

// Filter is an implementation of a cuckoo filter.
type Filter struct {
	nBuckets uint32
	table    [][entriesPerBucket]uint16
	lock     *sync.RWMutex
	counters *counter.Oc
	wg       *sync.WaitGroup
	FilterConfig
}

type FilterConfig struct {
	MaxKeys          int
	EntriesPerBucket int
	FpBits           int
	MaxDisplacements int
	LoadFactor       float32
	WithFastMode     bool
	WithDebug        bool
	WithStats        bool
}

type stats struct {
	indexed  int
	deleted  int
	found    int
	requests int
	max      int
}

func NewWithConfig(c *FilterConfig) (f *Filter) {
	if c.MaxKeys <= 0 {
		c.maxKeys = defaultMaxKeys
	}

	if c.EntriesPerBucket == nil {
		c.EntriesPerBucket = defaultEntriesPerBucket
	}

	if c.LoadFactor == nil {
		c.LoadFactor = defaultLoadFactor
	}

	if c.FpBits == nil {
		c.FpBits = defaultFingerprintsBits
	}

	if c.MaxDisplacements == nil {
		c.MaxDisplacements = defaultMaxDisplacements
	}

	nBuckets := nearestPowerOfTwo(c.MaxKeys / c.EntriesPerBucket)

	// If load factor is above the max value, we'll likely hit the max number of fingerprint displacements. In that case, expand the number of buckets.
	if float64(c.MaxKeys)/float64(nBuckets)/c.EntriesPerBucket > loadFactor {
		nBuckets <<= 1
	}

	f = &Filter{
		nBuckets,
		make([][epb]uint16, nBuckets),
		&sync.RWMutex{},
		nil,
		nil,
		c,
	}

	if c.WithStats {
		f.counters = counter.NewOc()
	}

	if c.WithFastMode {
		f.wg = &sync.WaitGroup{} // bulk add filters ?!
	}

	return

}

// New returns a new cuckoo filter sized for the maximum number of keys passed in as maxKeys.
func New(maxKeys uint32, fm bool, debug bool) (f *Filter) {

	if maxKeys <= 0 {
		maxKeys = defaultMaxKeys
	}

	if epb == nil {
		epb = defaultEntriesPerBucket
	}

	if lf == nil {
		lf = defaultLoadFactor
	}

	if fb == nil {
		fb = defaultFingerprintsBits
	}

	if md == nil {
		md = defaultMaxDisplacements
	}

	nBuckets := nearestPowerOfTwo(maxKeys / epb)

	// If load factor is above the max value, we'll likely hit the max number of fingerprint displacements. In that case, expand the number of buckets.
	if float64(maxKeys)/float64(nBuckets)/epb > loadFactor {
		nBuckets <<= 1
	}

	f = &Filter{
		maxKeys:          maxKeys,
		entriesPerBucket: epb,
		loadFactor:       lf,
		maxDisplacements: defaultMaxDisplacements,
		fpBits:           fb,
		nBuckets:         nBuckets,
		table:            make([][epb]uint16, nBuckets),
		counters:         counter.NewOc(),
		lock:             &sync.RWMutex{},
		debug:            debug,
	}

	if fm {
		f.wg = &sync.WaitGroup{} // bulk add filters ?!
	}

	return
}

func (f *Filter) Count(item uint32) (count uint32) {
	f.lock.RLock()
	count = f.nBuckets
	f.lock.RUnlock()
	return
}

func (f *Filter) bucketIndex(hv uint32) (res uint32) {
	f.lock.RLock()
	res = hv % f.nBuckets
	f.lock.RUnlock()
	return
}

func (f *Filter) fingerprint(hv uint32) (fp uint16) {
	fp = uint16(hv & ((1 << fpBits) - 1))

	// gross
	if fp == 0 {
		fp = 1
	}

	return
}

func (f *Filter) alternateIndex(idx uint32, fp uint16) uint32 {
	d := make([]byte, 2)
	binary.LittleEndian.PutUint16(d, fp)
	hv := hash(d)
	return f.bucketIndex(idx ^ uint32(hv))
}

func (f *Filter) matchPosition(idx uint32, fp uint16) int {
	for i := 0; i < entriesPerBucket; i++ {
		if f.table[idx][i] == fp {
			return i
		}
	}
	return -1
}

func (f *Filter) emptyPosition(idx uint32) (res int) {
	f.lock.RLock()
	res = f.matchPosition(idx, 0)
	f.lock.RUnlock()
	return
}

// filters []*string ? filtering ?!
func (f *Filter) Stats() (stats map[string]int) {
	stats = make(map[string]int, 4)
	f.lock.RLock()
	f.counters.SortByKey(counter.ASC)

	for f.counters.Next() {
		key, value := f.counters.KeyValue()
		stats[key] = value
	}

	f.lock.RUnlock()
	return
}

// Add adds an element to the cuckoo filter.  If the filter is too
// heavily loaded, ErrTooFull may be returned, which signifies that
// the filter must be rebuilt with an increased maxKeys parameter.
func (f *Filter) Add(d []byte) error {
	h := hash(d)

	f.lock.RWLock()
	fp := f.fingerprint(uint32(h))
	i1 := f.bucketIndex(uint32(h >> 32))
	i2 := f.alternateIndex(i1, fp)

	if i := f.emptyPosition(i1); i != -1 {
		f.table[i1][i] = fp
		f.counters.Increment("indexed", 1)
		return nil
	}

	if i := f.emptyPosition(i2); i != -1 {
		f.table[i2][i] = fp
		f.counters.Increment("indexed", 1)
		return nil
	}

	// Choose which index to use randomly
	idx := [2]uint32{i1, i2}[rand.Intn(2)]

	for i := 0; i < maxDisplacements; i++ {
		j := uint32(rand.Intn(entriesPerBucket))

		fp, f.table[idx][j] = f.table[idx][j], fp
		idx = f.alternateIndex(idx, fp)

		if ni := f.emptyPosition(idx); ni != -1 {
			f.table[idx][ni] = fp
			f.counters.Increment("indexed", 1)
			return nil
		}
	}
	f.lock.RWUnlock()

	return ErrTooFull
}

// Contains returns whether an element may be present in the set.
// Cuckoo filters are probablistic data structures which can return
// false positives.  False negatives are not possible.
func (f *Filter) Contains(d []byte) (res bool) {
	h := hash(d)

	f.lock.RLock()
	fp := f.fingerprint(uint32(h))
	i1 := f.bucketIndex(uint32(h >> 32))
	i2 := f.alternateIndex(i1, fp)
	f.counters.Increment("found", 1)
	res = f.matchPosition(i1, fp) != -1 || f.matchPosition(i2, fp) != -1
	f.lock.RUnlock()

	return
}

// Delete deletes an element from the set.  To delete an item safely,
// it must have been previously inserted.  Deleting a non-inserted
// item might unintentionally remove a real, different item.
func (f *Filter) Delete(d []byte) bool {
	h := hash(d)

	f.lock.RWLock()
	fp := f.fingerprint(uint32(h))
	i1 := f.bucketIndex(uint32(h >> 32))
	i2 := f.alternateIndex(i1, fp)

	if i := f.matchPosition(i1, fp); i != -1 {
		f.table[i1][i] = 0
		f.counters.Increment("deleted", 1)
		return true
	}

	if i := f.matchPosition(i2, fp); i != -1 {
		f.table[i2][i] = 0
		f.counters.Increment("deleted", 1)
		return true
	}

	f.lock.RWUnlock()
	return false
}

// to do
func (f *Filter) Keys() []string {
	// f.lock.RLock()
	// f.lock.RUnlock()

	return []string{}
}
