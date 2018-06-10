package tachymeter

import (
	"fmt"
	"math"
	"time"
)

// timeslice holds time.Duration values.
type timeSlice []time.Duration

// Satisfy sort for timeSlice.
func (p timeSlice) Len() int           { return len(p) }
func (p timeSlice) Less(i, j int) bool { return int64(p[i]) < int64(p[j]) }
func (p timeSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// hgram returns a histogram of event durations in
// b bins, along with the bin size.
func (ts timeSlice) hgram(b int) (*Histogram, time.Duration) {
	res := time.Duration(1000)
	// Interval is the time range / n bins.
	interval := time.Duration(int64(ts.srange()) / int64(b))
	high := ts.min() + interval
	low := ts.min()
	max := ts.max()
	hgram := &Histogram{}
	pos := 1 // Bin position.

	bstring := fmt.Sprintf("%s - %s", low/res*res, high/res*res)
	bin := map[string]uint64{}

	for _, v := range ts {
		// If v fits in the current bin,
		// increment the bin count.
		if v <= high {
			bin[bstring]++
		} else {
			// If not, prepare the next bin.
			*hgram = append(*hgram, bin)
			bin = map[string]uint64{}

			// Update the high/low range values.
			low = high + time.Nanosecond

			high += interval
			// if we're going into the
			// last bin, set high to max.
			if pos == b-1 {
				high = max
			}

			bstring = fmt.Sprintf("%s - %s", low/res*res, high/res*res)

			// The value didn't fit in the previous
			// bin, so the new bin count should
			// be incremented.
			bin[bstring]++

			pos++
		}
	}

	*hgram = append(*hgram, bin)

	return hgram, interval
}

// These should be self-explanatory:

func (ts timeSlice) cumulative() time.Duration {
	var total time.Duration
	for _, t := range ts {
		total += t
	}

	return total
}

func (ts timeSlice) hMean() time.Duration {
	var total float64

	for _, t := range ts {
		total += (1 / float64(t))
	}

	return time.Duration(float64(ts.Len()) / total)
}

func (ts timeSlice) avg() time.Duration {
	var total time.Duration
	for _, t := range ts {
		total += t
	}
	return time.Duration(int(total) / ts.Len())
}

func (ts timeSlice) p(p float64) time.Duration {
	return ts[int(float64(ts.Len())*p+0.5)-1]
}

func (ts timeSlice) stdDev() time.Duration {
	m := ts.avg()
	s := 0.00

	for _, t := range ts {
		s += math.Pow(float64(m-t), 2)
	}

	msq := s / float64(ts.Len())

	return time.Duration(math.Sqrt(msq))
}

func (ts timeSlice) long5p() time.Duration {
	set := ts[int(float64(ts.Len())*0.95+0.5):]

	if len(set) <= 1 {
		return ts[ts.Len()-1]
	}

	var t time.Duration
	var i int
	for _, n := range set {
		t += n
		i++
	}

	return time.Duration(int(t) / i)
}

func (ts timeSlice) short5p() time.Duration {
	set := ts[:int(float64(ts.Len())*0.05+0.5)]

	if len(set) <= 1 {
		return ts[0]
	}

	var t time.Duration
	var i int
	for _, n := range set {
		t += n
		i++
	}

	return time.Duration(int(t) / i)
}

func (ts timeSlice) min() time.Duration {
	return ts[0]
}

func (ts timeSlice) max() time.Duration {
	return ts[ts.Len()-1]
}

func (ts timeSlice) srange() time.Duration {
	return ts.max() - ts.min()
}
