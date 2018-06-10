package metric

import (
	"fmt"
	"sort"
	"time"
)

type Snapshot struct {
	// time
	timestamp           time.Time
	timeSinceStart      time.Duration
	averageResponseTime time.Duration

	// counters
	numberOfWorkers              int
	totalNumberOfRequests        int
	numberOfSuccessfulRequests   int
	numberOfUnsuccessfulRequests int
	numberOfRequestsPerSecond    float64

	numberOfRequestsByStatusCode  map[int]int
	numberOfRequestsByContentType map[string]int

	// lists
	listOfResponsesContentTypes   map[string]int
	listOfResponsesStatusCodes    map[string]int
	listOfResponsesFiltersMatches map[string]int

	// size
	totalSizeInBytes   int
	averageSizeInBytes int
}

func (s Snapshot) Timestamp() time.Time {
	return s.timestamp
}

func (s Snapshot) NumberOfWorkers() int {
	return s.numberOfWorkers
}

func (s Snapshot) NumberOfErrors() int {
	return s.numberOfUnsuccessfulRequests
}

func (s Snapshot) TotalNumberOfRequests() int {
	return s.totalNumberOfRequests
}

func (s Snapshot) TotalSizeInBytes() int {
	return s.totalSizeInBytes
}

func (s Snapshot) AverageSizeInBytes() int {
	return s.averageSizeInBytes
}

func (s Snapshot) AverageResponseTime() time.Duration {
	return s.averageResponseTime
}

func (s Snapshot) RequestsPerSecond() float64 {
	return s.numberOfRequestsPerSecond
}

func GetLatestLogMessages(messages []string, count int) ([]string, error) {
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
