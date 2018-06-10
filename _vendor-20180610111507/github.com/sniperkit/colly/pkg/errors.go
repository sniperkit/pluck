package colly

import (
	"errors"
)

var (
	// Cache dir not found
	ErrCacheDirNotFound = errors.New("Cache directory not found")

	// Cache dir not created
	ErrCacheDirFailed = errors.New("Cache directory could not be created")

	// ErrForbiddenDomain is the error thrown if visiting a domain which is not allowed in AllowedDomains
	ErrForbiddenDomain = errors.New("Forbidden domain")

	// ErrMissingURL is the error type for missing URL errors
	ErrMissingURL = errors.New("Missing URL")

	// ErrMaxDepth is the error type for exceeding max depth
	ErrMaxDepth = errors.New("Max depth limit reached")

	// ErrForbiddenURL is the error thrown if visiting a URL which is not allowed by URLFilters
	ErrForbiddenURL = errors.New("ForbiddenURL")

	// ErrNoURLFiltersMatch is the error thrown if visiting a URL which is not allowed by URLFilters
	ErrNoURLFiltersMatch = errors.New("No URLFilters match")

	// ErrAlreadyVisited is the error type for already visited URLs
	ErrAlreadyVisited = errors.New("URL already visited")

	// ErrRobotsTxtBlocked is the error type for robots.txt errors
	ErrRobotsTxtBlocked = errors.New("URL blocked by robots.txt")

	// ErrNoCookieJar is the error type for missing cookie jar
	ErrNoCookieJar = errors.New("Cookie jar is not available")

	// ErrNoPattern is the error type for LimitRules without patterns
	ErrNoPattern = errors.New("No pattern defined in LimitRule")

	// ErrNotValidTabularFormat
	ErrNotValidTabularFormat = errors.New("Not valid tabular format")

	// ErrInvalidTabularQuery
	ErrInvalidTabularQuery = errors.New("Invalid tabular query. Get a row or multiple rows: `0`, `0,1` or slice a Dataset: `0:5`")

	// ErrTabularInvalidQuery
	ErrTabularInvalidQuery = errors.New("Invalid tabular query.")

	// ErrTabularRowSelectionNotImplemented
	ErrTabularRowSelectionNotImplemented = errors.New("Tabular row selection is not implemented yet.")

	// ErrTabularMixedSelectionNotImplemented
	ErrTabularMixedSelectionNotImplemented = errors.New("Tabular with custom row and col selection is not implemented yet.")
)
