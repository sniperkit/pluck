package storage

import (
	"net/http/cookiejar"
	"net/url"
	"sync"
)

// InMemoryStorage is the default storage backend of colly.
// InMemoryStorage keeps cookies and visited urls in memory
// without persisting data on the disk.
type InMemoryStorage struct {
	visitedURLs map[uint64]bool
	lock        *sync.RWMutex
	jar         *cookiejar.Jar
}

// Init initializes InMemoryStorage
func (s *InMemoryStorage) Init() error {
	if s.visitedURLs == nil {
		s.visitedURLs = make(map[uint64]bool)
	}
	if s.lock == nil {
		s.lock = &sync.RWMutex{}
	}
	if s.jar == nil {
		var err error
		s.jar, err = cookiejar.New(nil)
		return err
	}
	return nil
}

// Visited implements Storage.Visited()
func (s *InMemoryStorage) Visited(requestID uint64) error {
	s.lock.Lock()
	s.visitedURLs[requestID] = true
	s.lock.Unlock()
	return nil
}

// IsVisited implements Storage.IsVisited()
func (s *InMemoryStorage) IsVisited(requestID uint64) (bool, error) {
	s.lock.RLock()
	visited := s.visitedURLs[requestID]
	s.lock.RUnlock()
	return visited, nil
}

// Cookies implements Storage.Cookies()
func (s *InMemoryStorage) Cookies(u *url.URL) string {
	return StringifyCookies(s.jar.Cookies(u))
}

// SetCookies implements Storage.SetCookies()
func (s *InMemoryStorage) SetCookies(u *url.URL, cookies string) {
	s.jar.SetCookies(u, UnstringifyCookies(cookies))
}

// Close implements Storage.Close()
func (s *InMemoryStorage) Close() error {
	return nil
}
