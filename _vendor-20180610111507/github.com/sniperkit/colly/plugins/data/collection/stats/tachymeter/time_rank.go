package tachymeter

import (
	"fmt"
	"math"
	"time"
)

type ranking struct {
	label     string
	startedAt time.Time
	endedAt   time.Time
	duration  time.Duration
	err       bool
}

// timeRank holds time.Duration values.
type timeRank []ranking

// Satisfy sort for timeRank.
func (p timeRank) Len() int           { return len(p) }
func (p timeRank) Less(i, j int) bool { return int64(p[i].duration) < int64(p[j].duration) } //  p[i].duration.Before(p[j].duration) }
func (p timeRank) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// These should be self-explanatory:
func (tr timeRank) hMean() time.Duration {
	var total float64
	for _, t := range tr {
		total += (1 / float64(t.duration))
	}
	return time.Duration(float64(tr.Len()) / total)
}

func (tr timeRank) cumulative() time.Duration {
	var total time.Duration
	for _, t := range tr {
		total += t.duration
	}
	return total
}

func (tr timeRank) avg() time.Duration {
	var total time.Duration
	for _, t := range tr {
		total += t.duration
	}
	return time.Duration(int(total) / tr.Len())
}

func (tr timeRank) p(p float64) time.Duration {
	return tr[int(float64(tr.Len())*p+0.5)-1].duration
}

func (tr timeRank) long5p() time.Duration {
	set := tr[int(float64(tr.Len())*0.95+0.5):]
	if len(set) <= 1 {
		return tr[tr.Len()-1].duration
	}
	var t time.Duration
	var i int
	for _, n := range set {
		t += n.duration
		i++
	}

	return time.Duration(int(t) / i)
}

func (tr timeRank) short5p() time.Duration {
	set := tr[:int(float64(tr.Len())*0.05+0.5)]
	if len(set) <= 1 {
		return tr[0].duration
	}
	var t time.Duration
	var i int
	for _, n := range set {
		t += n.duration
		i++
	}
	return time.Duration(int(t) / i)
}

func (tr timeRank) srange() time.Duration {
	return tr.max() - tr.min()
}

func (tr timeRank) p50() time.Duration {
	k := tr.Len() / 2
	return tr[k].duration
}

func (tr timeRank) min() time.Duration {
	return tr[0].duration
}

func (tr timeRank) stdDev() time.Duration {
	m := tr.avg()
	s := 0.00
	for _, t := range tr {
		s += math.Pow(float64(m-t.duration), 2)
	}
	msq := s / float64(tr.Len())
	return time.Duration(math.Sqrt(msq))
}

func (tr timeRank) minStr() string {
	return fmt.Sprintf("label=%s, duration=%s, len=%d", tr[0].label, tr[0].duration, len(tr))
}

func (tr timeRank) max() time.Duration {
	k := tr.Len() - 1
	return tr[k].duration
}

func (tr timeRank) maxStr() string {
	k := tr.Len() - 1
	return fmt.Sprintf("label=%s, duration=%s, len=%d", tr[k].label, tr[k].duration, len(tr))
}

// hgram returns a histogram of event durations in b buckets, along with the bucket size.
func (tr timeRank) hgram(b int) (*Histogram, time.Duration) {
	res := time.Duration(1000)
	// Interval is the time range / n buckets.
	interval := time.Duration(int64(tr.srange()) / int64(b))
	high := tr.min() + interval
	low := tr.min()
	max := tr.max()
	hgram := &Histogram{}
	pos := 1 // Bucket position.

	bstring := fmt.Sprintf("%s - %s", low/res*res, high/res*res)
	bucket := map[string]uint64{}

	for _, v := range tr {
		// If v fits in the current bucket,
		// increment the bucket count.
		if v.duration <= high {
			bucket[bstring]++
		} else {
			// If not, prepare the next bucket.
			*hgram = append(*hgram, bucket)
			bucket = map[string]uint64{}

			// Update the high/low range values.
			low = high + time.Nanosecond

			high += interval
			// if we're going into the
			// last bucket, set high to max.
			if pos == b-1 {
				high = max
			}

			bstring = fmt.Sprintf("%s - %s", low/res*res, high/res*res)

			// The value didn't fit in the previous
			// bucket, so the new bucket count should
			// be incremented.
			bucket[bstring]++

			pos++
		}
	}

	*hgram = append(*hgram, bucket)

	return hgram, interval
}
