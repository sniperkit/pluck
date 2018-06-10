package main

import (
	config "github.com/sniperkit/colly/plugins/data/extract/text/pluck/config"
)

var PLUCKER_CONFIG_UNITS = []*config.Config{
	{
		Name:        "row",                            // Name sets the key in the returned map, after completion
		Activators:  []string{"row[", "\"", ":", ","}, // Activators must be found in order, before capturing commences
		Deactivator: "]",                              // Deactivator restarts capturing
		Limit:       -1,                               // Limit specifies the number of times capturing can occur
		// Finisher:    "]",                        // Finisher trigger the end of capturing this pluck
		// Permanent:   1,                          // Permanent set the number of activators that stay permanently (counted from left to right)
		// Sanitize:    true,                       // Sanitize enables the html stripping
		// Maximum:     1,                          // Maximum set the number of characters for a capture
	},
	{
		Name:        "rows",
		Activators:  []string{"row[", "\"", ":", ","},
		Deactivator: "]",
		Limit:       -1,
	},
	{
		Name:        "col",
		Activators:  []string{"row[", "\"", ":", ","},
		Deactivator: "]",
		Limit:       -1,
	},
	{
		Name:        "cols",
		Activators:  []string{"row[", "\"", ":", ","},
		Deactivator: "]",
		Limit:       -1,
	},
	/*
		{
			Activators:  []string{"Section 2", "a", "href", `"`},
			Permanent:   1,
			Deactivator: `"`,
			Finisher:    "Section 3",
			Limit:       -1,
			Name:        "0",
			Maximum:     6,
		},
	*/

}
