package cui

import (
	"fmt"
	"sort"
	"strconv"
	"sync"
	"time"
)

var stats Statistics

func init() {
	stats = Statistics{
		lock: sync.RWMutex{},
		// counters
		numberOfRequestsByStatusCode:  make(map[int]int),
		numberOfRequestsByContentType: make(map[string]int),
		// top lists
		listOfResponsesContentTypes:   make(map[string]int),
		listOfResponsesStatusCodes:    make(map[string]int),
		listOfResponsesFiltersMatches: make(map[string]int),
	}
}

func UpdateStatistics(w WorkResult) {
	go stats.Add(w)
}

type Statistics struct {
	lock sync.RWMutex

	rawResults  []WorkResult
	snapShots   []Snapshot
	logMessages []string

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

func (s *Statistics) Add(w WorkResult) Snapshot {
	// update the raw results
	s.lock.Lock()
	defer s.lock.Unlock()

	s.rawResults = append(s.rawResults, w)

	// initialize start and end time
	if s.numberOfRequests == 0 {
		s.startTime = w.GetStartTime()
		s.endTime = w.GetEndTime()
	}

	// start time
	if w.GetStartTime().Before(s.startTime) {
		s.startTime = w.GetStartTime()
	}

	// end time
	if w.GetEndTime().After(s.endTime) {
		s.endTime = w.GetEndTime()
	}

	// update the total number of requests
	s.numberOfRequests = len(s.rawResults)

	// is successful
	if w.GetStatusCode() > 199 && w.GetStatusCode() < 400 {
		s.numberOfSuccessfulRequests += 1
	} else {
		s.numberOfUnsuccessfulRequests += 1
	}

	// number of workers
	s.numberOfWorkers = w.GetNumberOfWorkers()

	// number of requests by status code
	s.numberOfRequestsByStatusCode[w.GetStatusCode()] += 1

	// number of requests by content type
	s.numberOfRequestsByContentType[w.GetContentType()] += 1

	ct := w.GetContentType()
	if ct == "" {
		ct = "unknown"
	}
	s.listOfResponsesContentTypes[ct] += 1

	sc := strconv.Itoa(w.GetStatusCode())
	if sc == "" {
		sc = "unknown"
	}
	s.listOfResponsesStatusCodes[sc] += 1

	/*
		fm := w.GetFiltersMatches()
		if fm == "" {
			fm = "unknown"
		}
		 s.listOfResponsesFiltersMatches[w.GetFiltersMatches()] += 1
	*/

	// update the total duration
	responseTime := w.GetEndTime().Sub(w.GetStartTime())
	s.totalResponseTime += responseTime

	// size
	s.totalSizeInBytes += w.GetSize()
	averageSizeInBytes := s.totalSizeInBytes / s.numberOfRequests

	// average response time
	averageResponseTime := time.Duration(s.totalResponseTime.Nanoseconds() / int64(s.numberOfRequests))

	// number of requests per second
	requestsPerSecond := float64(s.numberOfRequests) / s.endTime.Sub(s.startTime).Seconds()

	// log messages
	s.logMessages = append(s.logMessages, w.String())

	// create a snapshot
	snapShot := Snapshot{

		// times
		timestamp:           w.GetEndTime(),
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

	messages, err := getLatestLogMessages(s.logMessages, count)
	if err != nil {
		panic(err)
	}

	return messages
}

func getLatestLogMessages(messages []string, count int) ([]string, error) {
	if count < 0 {
		return nil, fmt.Errorf("The count cannot be negative")
	}

	numberOfMessages := len(messages)
	if count == numberOfMessages {
		return messages, nil
	}

	if count < numberOfMessages {
		return messages[numberOfMessages-count:], nil
	}

	if count > numberOfMessages {
		fillLines := make([]string, count-numberOfMessages)
		return append(fillLines, messages...), nil
	}
	panic("Unreachable")
}

func sortTopList(input map[string]int) (sorted []string) {
	n := map[int][]string{}
	var a []int
	for k, v := range input {
		n[v] = append(n[v], k)
	}
	for k := range n {
		a = append(a, k)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(a)))
	for _, k := range a {
		for _, s := range n[k] {
			sorted = append(sorted, s)
			// fmt.Printf("%s, %d\n", s, k)
		}
	}
	return
}

func (s *Statistics) TopList(count int) (res []string) { return }

// Get top list of content types responses
/*
	notes:
	- sort string slices with 'sortedKeys' https://play.golang.org/p/x4CoUsJ5tK
	- https://github.com/indraniel/go-learn/blob/master/09-sort-map-keys-by-values.go
	- https://stackoverflow.com/questions/18695346/how-to-sort-a-mapstringint-by-its-values
*/

func (s *Statistics) TopContentTypes(count int) []string {
	s.lock.RLock()
	defer s.lock.RUnlock()
	types, err := getTopContentTypes(s.listOfResponsesContentTypes, count)
	if err != nil {
		panic(err)
	}
	return types
}

func getTopContentTypes(types map[string]int, count int) ([]string, error) {
	if count < 0 {
		return nil, fmt.Errorf("The count cannot be negative")
	}
	topList := sortTopList(types)
	numberOfContentTypes := len(topList)
	if count == numberOfContentTypes {
		return topList, nil
	}
	if count < numberOfContentTypes {
		return topList[numberOfContentTypes-count:], nil
	}
	if count > numberOfContentTypes {
		fillLines := make([]string, count-numberOfContentTypes)
		return append(fillLines, topList...), nil
	}
	panic("Unreachable")
}

// Get top list of status code responses
func (s *Statistics) TopStatusCodes(count int) []string {
	s.lock.RLock()
	defer s.lock.RUnlock()
	codes, err := getTopStatusCodes(s.listOfResponsesStatusCodes, count)
	if err != nil {
		panic(err)
	}
	return codes
}

func getTopStatusCodes(codes map[string]int, count int) ([]string, error) {
	if count < 0 {
		return nil, fmt.Errorf("The count cannot be negative")
	}
	topList := sortTopList(codes)
	numberOfStatusCodes := len(topList)
	if count == numberOfStatusCodes {
		return topList, nil
	}
	if count < numberOfStatusCodes {
		return topList[numberOfStatusCodes-count:], nil
	}
	if count > numberOfStatusCodes {
		fillLines := make([]string, count-numberOfStatusCodes)
		return append(fillLines, topList...), nil
	}
	panic("Unreachable")
}

// Get top list of status code responses
func (s *Statistics) TopFiltersMatches(count int) []string {
	s.lock.RLock()
	defer s.lock.RUnlock()
	filters, err := getTopFiltersMatches(s.listOfResponsesFiltersMatches, count)
	if err != nil {
		panic(err)
	}
	return filters
}

func getTopFiltersMatches(filters map[string]int, count int) ([]string, error) {
	if count < 0 {
		return nil, fmt.Errorf("The count cannot be negative")
	}
	topList := sortTopList(filters)
	numberOfFiltersMatches := len(topList)
	if count == numberOfFiltersMatches {
		return topList, nil
	}
	if count < numberOfFiltersMatches {
		return topList[numberOfFiltersMatches-count:], nil
	}
	if count > numberOfFiltersMatches {
		fillLines := make([]string, count-numberOfFiltersMatches)
		return append(fillLines, topList...), nil
	}
	panic("Unreachable")
}
