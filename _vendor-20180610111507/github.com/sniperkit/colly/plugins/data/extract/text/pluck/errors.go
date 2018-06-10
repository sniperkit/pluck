package plucker

import (
	"errors"
)

const (
	errMarshallingResults string = "result marshalling failed"
	errFailedToOpen       string = "problem opening "
)

var (
	errEmptyResults = errors.New("No results to encode")
)
