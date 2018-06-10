package metric

import (
	"strconv"
	"sync"
	"time"

	cmap "github.com/sniperkit/colly/plugins/data/structure/map/multi" // Concurrent multi-map
)

var (
	stats                          Statistics
	statsConcurrentMap             *cmap.ConcurrentMap             // concurrent map
	statsShardedConcurrentMap      *cmap.ShardedConcurrentMap      // concurrent map with shards. (Default: 16)
	statsConcurrentMultiMap        *cmap.ConcurrentMultiMap        // concurrent multi-map
	statsShardedConcurrentMultiMap *cmap.ShardedConcurrentMultiMap // concurrent multi-map with shards
)

func InitStatsCollector() {

	stats = Statistics{
		lock: sync.RWMutex{},
		// concurrent maps
		cLogMessages: cmap.NewConcurrentMap(),
		cLogFilters:  cmap.NewConcurrentMap(),
		// counters
		numberOfRequestsByStatusCode:  make(map[int]int),
		numberOfRequestsByContentType: make(map[string]int),
		// top lists
		listOfResponsesContentTypes:   make(map[string]int),
		listOfResponsesStatusCodes:    make(map[string]int),
		listOfResponsesFiltersMatches: make(map[string]int),
	}
}

func UpdateStatistics(r Response) {
	go stats.Add(r)
}

type Statistics struct {
	lock sync.RWMutex

	rawResults  []Response
	snapShots   []Snapshot
	logMessages []string
	logFilters  []string

	// Test
	cLogMessages *cmap.ConcurrentMap
	cLogFilters  *cmap.ConcurrentMap

	startTime time.Time
	endTime   time.Time

	totalResponseTime time.Duration

	numberOfWorkers               int
	numberOfRequests              int
	numberOfSuccessfulRequests    int
	numberOfUnsuccessfulRequests  int
	numberOfRequestsByStatusCode  map[int]int
	numberOfRequestsByContentType map[string]int

	// lists
	listOfResponsesContentTypes   map[string]int
	listOfResponsesStatusCodes    map[string]int
	listOfResponsesFiltersMatches map[string]int

	totalSizeInBytes int
}

func NewStatsCollector() *Statistics {
	stats := &Statistics{
		lock: sync.RWMutex{},
		// concurrent maps
		cLogMessages: cmap.NewConcurrentMap(),
		cLogFilters:  cmap.NewConcurrentMap(),
		// counters
		numberOfRequestsByStatusCode:  make(map[int]int),
		numberOfRequestsByContentType: make(map[string]int),
		// top lists
		listOfResponsesContentTypes:   make(map[string]int),
		listOfResponsesStatusCodes:    make(map[string]int),
		listOfResponsesFiltersMatches: make(map[string]int),
	}
	return stats
}

func (s *Statistics) UpdateStatistics(r Response) {
	go s.Add(r)
}

func (s *Statistics) Add(r Response) Snapshot {
	// update the raw results
	s.lock.Lock()
	defer s.lock.Unlock()

	s.rawResults = append(s.rawResults, r)

	// initialize start and end time
	if s.numberOfRequests == 0 {
		s.startTime = r.GetStartTime()
		s.endTime = r.GetEndTime()
	}

	// start time
	if r.GetStartTime().Before(s.startTime) {
		s.startTime = r.GetStartTime()
	}

	// end time
	if r.GetEndTime().After(s.endTime) {
		s.endTime = r.GetEndTime()
	}

	// update the total number of requests
	s.numberOfRequests = len(s.rawResults)

	// is successful
	if r.GetStatusCode() > 199 && r.GetStatusCode() < 400 {
		s.numberOfSuccessfulRequests += 1
	} else {
		s.numberOfUnsuccessfulRequests += 1
	}

	// number of workers
	s.numberOfWorkers = r.GetNumberOfWorkers()

	// number of requests by status code
	s.numberOfRequestsByStatusCode[r.GetStatusCode()] += 1

	// number of requests by content type
	s.numberOfRequestsByContentType[r.GetContentType()] += 1

	ct := r.GetContentType()
	if ct == "" {
		ct = "unknown"
	}
	s.listOfResponsesContentTypes[ct] += 1

	sc := strconv.Itoa(r.GetStatusCode())
	if sc == "" {
		sc = "unknown"
	}
	s.listOfResponsesStatusCodes[sc] += 1

	// update the total duration
	responseTime := r.GetEndTime().Sub(r.GetStartTime())
	s.totalResponseTime += responseTime

	// size
	s.totalSizeInBytes += r.GetSize()
	averageSizeInBytes := s.totalSizeInBytes / s.numberOfRequests

	// average response time
	averageResponseTime := time.Duration(s.totalResponseTime.Nanoseconds() / int64(s.numberOfRequests))

	// number of requests per second
	requestsPerSecond := float64(s.numberOfRequests) / s.endTime.Sub(s.startTime).Seconds()

	// log messages
	s.logMessages = append(s.logMessages, r.String())

	// filtering messages
	// s.logFilters = append(s.logFilters, r.String())

	// create a snapshot
	snapShot := Snapshot{

		// times
		timestamp:           r.GetEndTime(),
		averageResponseTime: averageResponseTime,

		// counters
		numberOfWorkers:               s.numberOfWorkers,
		totalNumberOfRequests:         s.numberOfRequests,
		numberOfSuccessfulRequests:    s.numberOfSuccessfulRequests,
		numberOfUnsuccessfulRequests:  s.numberOfUnsuccessfulRequests,
		numberOfRequestsPerSecond:     requestsPerSecond,
		numberOfRequestsByStatusCode:  s.numberOfRequestsByStatusCode,
		numberOfRequestsByContentType: s.numberOfRequestsByContentType,

		// size
		totalSizeInBytes:   s.totalSizeInBytes,
		averageSizeInBytes: averageSizeInBytes,
	}

	// pp.Println(snapShot)
	s.snapShots = append(s.snapShots, snapShot)
	return snapShot
}

func (s *Statistics) LastSnapshot() Snapshot {
	s.lock.RLock()
	defer s.lock.RUnlock()

	lastSnapshotIndex := len(s.snapShots) - 1
	if lastSnapshotIndex < 0 {
		return Snapshot{}
	}
	return s.snapShots[lastSnapshotIndex]
}

func (s *Statistics) LastLogMessages(count int) []string {
	s.lock.RLock()
	defer s.lock.RUnlock()

	messages, err := GetLatestLogMessages(s.logMessages, count)
	if err != nil {
		panic(err)
	}
	return messages
}
