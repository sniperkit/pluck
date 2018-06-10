package tachymeter

import (
	"bytes"
	"fmt"
	"math"
	"strings"
)

// Histogram is a map["low-high duration"]count of events that
// fall within the low-high time duration range.
type Histogram []map[string]uint64

// Dump prints a formatted histogram output to console
// scaled to a width of s.
func (h *Histogram) Dump(s int) {
	fmt.Println(h.String(s))
}

// String returns a formatted Metrics string scaled
// to a width of s.
func (h *Histogram) String(s int) string {
	if h == nil {
		return ""
	}

	var min, max uint64 = math.MaxUint64, 0
	// Get the histogram min/max counts.
	for _, bin := range *h {
		for _, v := range bin {
			if v > max {
				max = v
			}
			if v < min {
				min = v
			}
		}
	}

	// Handle cases of no or
	// a single bin.
	switch len(*h) {
	case 0:
		return ""
	case 1:
		min = 0
	}

	var b bytes.Buffer

	// Build histogram string.
	for _, bin := range *h {
		for k, v := range bin {
			// Get the bar length.
			blen := scale(float64(v), float64(min), float64(max), 1, float64(s))
			line := fmt.Sprintf("%20s %s\n", k, strings.Repeat("-", int(blen)))
			b.WriteString(line)
		}
	}

	return b.String()
}
