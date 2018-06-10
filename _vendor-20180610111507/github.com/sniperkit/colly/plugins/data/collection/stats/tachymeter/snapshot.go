package tachymeter

import (
	"math"
	"sort"
	"sync/atomic"
)

// Calc summarizes Tachymeter sample data and returns it in the form of a *Metrics.
func (m *Tachymeter) Calc() *Metrics {
	metrics := &Metrics{}
	if atomic.LoadUint64(&m.count) == 0 {
		return metrics
	}

	m.Lock()

	metrics.Samples = int(math.Min(float64(atomic.LoadUint64(&m.count)), float64(m.size)))
	metrics.Count = int(atomic.LoadUint64(&m.count))
	times := make(timeSlice, metrics.Samples)

	copy(times, m.times[:metrics.Samples])
	sort.Sort(times)

	metrics.Time.Cumulative = times.cumulative()
	var rateTime float64
	if m.wallTime != 0 {
		rateTime = float64(metrics.Count) / float64(m.wallTime)
	} else {
		rateTime = float64(metrics.Samples) / float64(metrics.Time.Cumulative)
	}

	metrics.Rate.Second = rateTime * 1e9
	m.Unlock()

	metrics.Time.Avg = times.avg()
	metrics.Time.HMean = times.hMean()
	metrics.Time.P50 = times[times.Len()/2]
	metrics.Time.P75 = times.p(0.75)
	metrics.Time.P95 = times.p(0.95)
	metrics.Time.P99 = times.p(0.99)
	metrics.Time.P999 = times.p(0.999)
	metrics.Time.Long5p = times.long5p()
	metrics.Time.Short5p = times.short5p()
	metrics.Time.Min = times.min()
	metrics.Time.Max = times.max()
	metrics.Time.Range = times.srange()
	metrics.Time.StdDev = times.stdDev()

	metrics.Histogram, metrics.HistogramBinSize = times.hgram(m.hBins)
	return metrics
}

// Snapshot summarizes Tachymeter sample data and returns it in the form of a *Metrics.
func (m *Tachymeter) Snapshot() *Metrics {

	metrics := &Metrics{}
	if atomic.LoadUint64(&m.count) == 0 {
		return metrics
	}

	m.Lock()

	metrics.Samples = int(math.Min(float64(atomic.LoadUint64(&m.count)), float64(m.size)))
	metrics.Count = int(atomic.LoadUint64(&m.count))
	metrics.Wall = m.wallTime

	ranks := make(timeRank, metrics.Samples)
	copy(ranks, m.ranks[:metrics.Samples])

	// GO 1.8 or above:
	// sort.Slice(ranks)
	// sort.Sort(ranks)

	sort.Slice(ranks, func(i, j int) bool { return int64(ranks[i].duration) < int64(ranks[j].duration) })

	metrics.Rank.Cumulative = ranks.cumulative()
	var rateTime float64
	if m.wallTime != 0 {
		rateTime = float64(metrics.Count) / float64(m.wallTime)
	} else {
		rateTime = float64(metrics.Samples) / float64(metrics.Time.Cumulative)
	}

	metrics.Rate.Second = rateTime * 1e9

	m.Unlock()

	metrics.Rank.Avg = ranks.avg()
	metrics.Rank.HMean = ranks.hMean()
	metrics.Rank.P50 = ranks.p50() // ranks.p(0.50) //[ranks.Len()/2]
	metrics.Rank.P75 = ranks.p(0.75)
	metrics.Rank.P95 = ranks.p(0.95)
	metrics.Rank.P99 = ranks.p(0.99)
	metrics.Rank.P999 = ranks.p(0.999)
	metrics.Rank.Long5p = ranks.long5p()
	metrics.Rank.Short5p = ranks.short5p()
	metrics.Rank.Max = ranks.maxStr()
	metrics.Rank.Min = ranks.minStr()
	metrics.Rank.Range = ranks.srange()
	metrics.Time.StdDev = ranks.stdDev()
	metrics.Histogram, metrics.HistogramBinSize = ranks.hgram(m.hBins)

	return metrics
}
