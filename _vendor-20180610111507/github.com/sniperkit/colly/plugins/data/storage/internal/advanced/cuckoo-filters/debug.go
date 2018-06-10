package cuckoofilters

import (
	"errors"
)

var (
	// ErrTooFull is returned when a filter is too full and needs to be resized.
	errTooFull = errors.New("cuckoo filter too full")

	// errNearestPowerOfTwo should never happen normally
	errNearestPowerOfTwo = errors.New("Nearest power of two error.")

	// errNearestPowerOfTwoStr should never happen normally
	errNearestPowerOfTwoStr string = "Nearest power of two error."
)
