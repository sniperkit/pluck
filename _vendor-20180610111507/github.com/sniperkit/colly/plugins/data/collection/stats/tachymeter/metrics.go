package tachymeter

import (
	"fmt"
	// "sync"
	"time"
)

// Metrics holds the calculated outputs
// produced from a Tachymeter sample set.
type MetricsLight struct {
	Time struct { // All values under Time are selected entirely from events within the sample window.
		Cumulative time.Duration // Cumulative time of all sampled events.
		HMean      time.Duration // Event duration harmonic mean.
		Avg        time.Duration // Event duration average.
		P50        time.Duration // Event duration nth percentiles ..
		P75        time.Duration
		P95        time.Duration
		P99        time.Duration
		P999       time.Duration
		Long5p     time.Duration // Average of the longest 5% event durations.
		Short5p    time.Duration // Average of the shortest 5% event durations.
		Max        time.Duration // Highest event duration.
		Min        time.Duration // Lowest event duration.
		StdDev     time.Duration // Standard deviation.
		Range      time.Duration // Event duration range (Max-Min).
	}
	Rate struct {
		// Per-second rate based on event duration avg. via Metrics.Cumulative / Metrics.Samples.
		// If SetWallTime was called, event duration avg = wall time / Metrics.Count
		Second float64
	}
	Histogram        *Histogram    // Frequency distribution of event durations in len(Histogram) bins of HistogramBinSize.
	HistogramBinSize time.Duration // The width of a histogram bin in time.
	Samples          int           // Number of events included in the sample set.
	Count            int           // Total number of events observed.
}

// Metrics holds the calculated outputs produced from a Tachymeter sample set.
type Metrics struct {
	// sync.Mutex
	Time struct { // All values under Time are selected entirely from events within the sample window.
		Cumulative time.Duration // Cumulative time of all sampled events.
		HMean      time.Duration // Event duration harmonic mean.
		Avg        time.Duration // Event duration average.
		P50        time.Duration // Event duration nth percentiles ..
		P75        time.Duration
		P95        time.Duration
		P99        time.Duration
		P999       time.Duration
		Long5p     time.Duration // Average of the longest 5% event durations.
		Short5p    time.Duration // Average of the shortest 5% event durations.
		Max        time.Duration // Highest event duration.
		Min        time.Duration // Lowest event duration.
		StdDev     time.Duration // Standard deviation.
		Range      time.Duration // Event duration range (Max-Min).
	}

	Rank struct {
		Cumulative time.Duration // Cumulative time of all sampled events.
		HMean      time.Duration // Event duration harmonic mean.
		Avg        time.Duration // Event duration average.
		P50        time.Duration // Event duration nth percentiles ..
		P75        time.Duration
		P95        time.Duration
		P99        time.Duration
		P999       time.Duration
		Long5p     time.Duration // Average of the longest 5% event durations.
		Short5p    time.Duration // Average of the shortest 5% event durations.
		Max        string
		Min        string
		Range      time.Duration // Event duration range (Max-Min).
	}

	Rate struct {
		// Per-second rate based on event duration avg. via Metrics.Cumulative / Metrics.Samples.
		// If SetWallTime was called, event duration avg = wall time / Metrics.Count
		Second float64
	}

	Abuse struct {
		Cumulative  time.Duration // Cumulative time of all sampled events.
		HMean       time.Duration // Event duration harmonic mean.
		Avg         time.Duration // Event duration average.
		TriggeredAt time.Time
		Second      float64
		Count       int
	}

	Events           map[string]bool
	Histogram        *Histogram    // Frequency distribution of event durations in len(Histogram) bins of HistogramBinSize.
	HistogramBinSize time.Duration // The width of a histogram bin in time.
	Samples          int           // Number of events included in the sample set.
	Count            int           // Total number of events observed.
	Wall             time.Duration
}

// WriteHTML writes a histograph html file to the cwd.
func (m *Metrics) WriteHTML(p string) error {
	w := Timeline{}
	w.AddEvent(m)
	return w.WriteHTML(p)
}

// Dump prints a formatted Metrics output to console.
func (m *Metrics) Dump() {
	fmt.Println(m.String())
}

// String returns a formatted Metrics string.
func (m *Metrics) String() string {
	return fmt.Sprintf(metricTmplTXT,
		m.Samples,
		m.Count,
		m.Wall.String(),
		m.Rank.Cumulative,
		m.Rank.HMean,
		m.Rank.Avg,
		m.Rank.P50,
		m.Rank.P75,
		m.Rank.P95,
		m.Rank.P99,
		m.Rank.P999,
		m.Rank.Long5p,
		m.Rank.Short5p,
		m.Rank.Max,
		m.Rank.Max,
		m.Rank.Min,
		m.Rank.Min,
		m.Rank.Range,
		m.Time.StdDev,
		m.Rate.Second)
}

// JSON returns a *Metrics as a JSON string.
func (m *Metrics) JSON() string {
	j, _ := json.Marshal(m)
	return string(j)
}

// MarshalJSON defines the output formatting
// for the JSON() method. This is exported as a
// requirement but not intended for end users.
func (m *Metrics) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Time struct {
			Cumulative string
			HMean      string
			Avg        string
			P50        string
			P75        string
			P95        string
			P99        string
			P999       string
			Long5p     string
			Short5p    string
			Max        string
			Min        string
			Range      string
			StdDev     string
		}
		Rate struct {
			Second float64
		}
		Samples   int
		Count     int
		Histogram *Histogram
	}{
		Time: struct {
			Cumulative string
			HMean      string
			Avg        string
			P50        string
			P75        string
			P95        string
			P99        string
			P999       string
			Long5p     string
			Short5p    string
			Max        string
			Min        string
			Range      string
			StdDev     string
		}{
			Cumulative: m.Time.Cumulative.String(),
			HMean:      m.Time.HMean.String(),
			Avg:        m.Time.Avg.String(),
			P50:        m.Time.P50.String(),
			P75:        m.Time.P75.String(),
			P95:        m.Time.P95.String(),
			P99:        m.Time.P99.String(),
			P999:       m.Time.P999.String(),
			Long5p:     m.Time.Long5p.String(),
			Short5p:    m.Time.Short5p.String(),
			Max:        m.Time.Max.String(),
			Min:        m.Time.Min.String(),
			Range:      m.Time.Range.String(),
			StdDev:     m.Time.StdDev.String(),
		},
		Rate: struct{ Second float64 }{
			Second: m.Rate.Second,
		},
		Histogram: m.Histogram,
		Samples:   m.Samples,
		Count:     m.Count,
	})
}
