// Package tachymeter yields summarized data describing a series of timed events.
package tachymeter

import (
	"bytes"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const PACKAGE_VERSION = "0.0.1"

const (
	TACHYMETER_DEFAULT_SAMPLE_CAPACITY   int  = 4
	TACHYMETER_DEFAULT_HISTOGRAM_BUCKETS int  = 10
	TACHYMETER_DEFAULT_SAFE_MODE         bool = true
	TACHYMETER_EXPORT_MAX_ROWS           int  = 1000000
	TACHYMETER_EXPORT_MAX_COLS           int  = 50
	TACHYMETER_MAX_DATABOOKS             int  = 5
	TACHYMETER_MAX_DATASETS              int  = 5
)

var (
	allowed_export_outputs                 []string = []string{"file", "io", "buffer", "string", "bytes"}
	allowed_export_datasets                []string = []string{"snapshots", "timeslice|time_slice", "ranks|ranking"}
	allowed_export_formats                 []string = []string{"yaml|yml", "csv", "json", "xml", "tsv", "xlsx", "ascii"}
	default_export_tachymeter_metrics_cols []string = []string{
		"SnapshotAt", "ExportedAt", "Wall", "Cumulative", "HMean", "Avg.", "p50", "p75", "p95", "p99", "p999", "Long 5%", "Short 5%", "Max", "Min", "Range", "Rate/sec.",
	}
	default_export_tachymeter_timeslices_cols []string = []string{"CreatedAt"}
	default_export_tachymeter_ranks_cols      []string = []string{"EventAt"}
)

/*
	Refs:
	- https://github.com/jamiealquiza/tachymeter/tree/master/example/tachymeter-graphing
	-
*/

// Config holds tachymeter initialization parameters. Size defines the sample capacity.
// Note: Tachymeter is thread safe.
type Config struct {
	SampleSize int  `default:'50'`
	HBins      int  `default:'10'`   // Histogram buckets. change for HBins
	SafeMode   bool `default:'true'` // Deprecated. Flag held on to as to not break existing users.
	Export     *Export
}

// Tachymeter holds event durations
// and counts.
type Tachymeter struct {
	sync.Mutex
	hBins    int
	safeMode bool
	size     uint64
	count    uint64
	wallTime time.Duration
	times    timeSlice
	ranks    timeRank
	buffer   *bytes.Buffer
}

// New initializes a new Tachymeter.
func New(c *Config) *Tachymeter {
	if c == nil {
		c = &Config{
			HBins:      10,
			SampleSize: 50,
			SafeMode:   true,
		}
	}

	var hSize int
	if c.HBins >= 0 {
		hSize = c.HBins
	} else {
		hSize = 10
	}

	return &Tachymeter{
		size:     uint64(c.SampleSize),
		ranks:    make(timeRank, c.SampleSize),
		hBins:    hSize,
		safeMode: c.SafeMode,
	}
}

// Setter for the sample size
func (m *Tachymeter) SetConfig(c *Config) *Tachymeter {
	m.Lock()
	defer m.Unlock()

	if c == nil {
		c = &Config{
			HBins:      TACHYMETER_DEFAULT_HISTOGRAM_BUCKETS,
			SampleSize: TACHYMETER_DEFAULT_SAMPLE_CAPACITY,
			SafeMode:   TACHYMETER_DEFAULT_SAFE_MODE,
		}
	}

	var hSize int
	if c.HBins >= 0 {
		hSize = c.HBins
	} else {
		hSize = TACHYMETER_DEFAULT_HISTOGRAM_BUCKETS
	}

	m.size = uint64(c.SampleSize)
	m.ranks = make(timeRank, c.SampleSize)
	m.hBins = hSize
	m.safeMode = c.SafeMode

	return m
}

// By default, tachymeter calcualtes rate based on the number of events
// possible per-second according to average event duration.
// This model doesn't work in asynchronous or parallelized scenarios since events
// may be overlapping in time. For example, with many Goroutines writing durations
// to a shared tachymeter in parallel, the global rate must be determined by using
// the total event count over the total wall time elapsed.
func (m *Tachymeter) WallTime(wallTime time.Duration) time.Duration {
	m.Lock()
	defer m.Unlock()
	return m.wallTime
}

// SetWallTime optionally sets an elapsed wall time duration.
// This affects rate output by using total events counted over time.
// This is useful for concurrent/parallelized events that overlap
// in wall time and are writing to a shared Tachymeter instance.
func (m *Tachymeter) SetWallTime(t time.Duration) {
	m.Lock()
	defer m.Unlock()
	m.wallTime = t
}

// SetWallTime optionally sets an elapsed wall time duration.
// This affects rate output by using total events counted over time.
// This is useful for concurrent/parallelized events that overlap
// in wall time and are writing to a shared Tachymeter instance.
func (m *Tachymeter) WittWallTime(t time.Duration) *Tachymeter {
	m.Lock()
	defer m.Unlock()
	m.wallTime = t
	return m
}

// AddTimeWithLabel adds a time.Duration to Tachymeter with a specific label.
func (m *Tachymeter) AddTimeWithLabel(label string, t time.Duration) {
	//	m.Times[(atomic.AddUint64(&m.Count, 1)-1)%m.Size] = t
	m.ranks[(atomic.AddUint64(&m.count, 1)-1)%m.size] = ranking{duration: t, label: label}
}

// AddTime adds a time.Duration to Tachymeter.
func (m *Tachymeter) AddTime(t time.Duration) {
	// m.times[(atomic.AddUint64(&m.count, 1)-1)%m.size] = t
	m.ranks[(atomic.AddUint64(&m.count, 1)-1)%m.size] = ranking{duration: t}
}

// Reset resets a Tachymeter instance for reuse.
func (m *Tachymeter) Reset() {
	// This lock is obviously not needed for  the m.Count update, rather to prevent a
	// Tachymeter reset while Calc is being called.
	m.Lock()
	atomic.StoreUint64(&m.count, 0)
	m.Unlock()
}

// Getter for the sample size
func (m *Tachymeter) Size() uint64 {
	m.Lock()
	defer m.Unlock()
	return m.size
}

// Setter for the sample size
func (m *Tachymeter) SetSize(size uint64) {
	m.Lock()
	defer m.Unlock()
	m.size = size
}

// Setter for the sample size
func (m *Tachymeter) WithSize(size uint64) *Tachymeter {
	m.Lock()
	defer m.Unlock()
	m.size = size
	return m
}

// Setter for the number of histogram buckets
func (m *Tachymeter) Buckets() int {
	m.Lock()
	defer m.Unlock()
	return m.hBins
}

// Setter for the number of histogram buckets
func (m *Tachymeter) WithHistBuckets(hBins int) *Tachymeter {
	m.Lock()
	defer m.Unlock()
	m.hBins = hBins
	return m
}

// Setter for the number of histogram buckets
func (m *Tachymeter) SetHistBuckets(hBins int) {
	m.Lock()
	defer m.Unlock()
	m.hBins = hBins
}

// String returns a formatted Metrics string.
func (m *Tachymeter) Convert(dataset string, format string) (output string, err error) {

	// select/prepare dataset
	var rawResults []string
	dataset = strings.ToLower(dataset)
	switch dataset {
	case "ranks", "ranking":
		// rawResults = m.ranks // []ranking: label, startedAt, endedAt, duration, err
	case "timeslice", "time_slice":
		// rawResults = m.times // []time.Duration
	default:
		rawResults = []string{""}
		err = errTachymeterDataset
		return
	}

	// export format
	format = strings.ToLower(format)
	switch format {
	case "postgres", "postgresql":
		output = fmt.Sprintf("%v", rawResults)
	case "ascii":
		output = fmt.Sprintf("%v", rawResults)
	case "mysql":
		output = fmt.Sprintf("%v", rawResults)
	case "xlsx":
		output = fmt.Sprintf("%v", rawResults)
	case "tsv":
		output = fmt.Sprintf("%v", rawResults)
	case "csv":
		output = fmt.Sprintf("%v", rawResults)
	case "xml":
		output = fmt.Sprintf("%v", rawResults)
	case "yaml":
		output = fmt.Sprintf("%v", rawResults)
	case "json":
		output = fmt.Sprintf("%v", rawResults)
	default:
		err = errTachymeterEncoding
		return
	}
	return
}

// Bytes(), String(), WriteTo(io.Writer), WriteFile(filename string, perm os.FileMode)

// Getter to check the count of samples registered
func (m *Tachymeter) Count() uint64 {
	m.Lock()
	defer m.Unlock()

	return m.count
}

// Scale scales the input x with the input-min a0,
// input-max a1, output-min b0, and output-max b1.
func scale(x, a0, a1, b0, b1 float64) float64 {
	a, b := x-a0, a1-a0
	var c float64
	if a == 0 {
		c = 0
	} else {
		c = a / b
	}
	return c*(b1-b0) + b0
}
