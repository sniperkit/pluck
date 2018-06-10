package main

import (
	"fmt"
	"strings"

	// colly - core
	colly "github.com/sniperkit/colly/pkg"
	debug "github.com/sniperkit/colly/pkg/debug"

	// colly - plugins
	pp "github.com/sniperkit/colly/plugins/app/debug/pp"
)

// app vars
var (
	// appVersion specifies the version of the app. If the executable is built with Makefile, the appVersion will use the actual repo's short tag version
	appVersion = "0.0.1-alpha"
	// appDebug specifies if the app debug/verbose some development event logged
	appDebug = true
)

// web target
var (
	targetRootURL = "https://golanglibs.com/sitemap.txt"
)

// collector vars
var (
	// collectorDebug sets collector's debugger
	collectorDebug = false
	// collectorDebugger stores the collector's log event listener
	collectorDebugger *debug.LogDebugger = &debug.LogDebugger{}
	// collectorJsonParser
	collectorJsonParser = "mxj"
	// collectorTabEnabled sets some debugging information
	collectorTabEnabled = true
	// collectorDatasetOutputPrefixPath specifies the prefix path for all saved dumps
	collectorDatasetOutputPrefixPath = "./shared/dataset"
	// collectorDatasetOutputBasename specifies the template to use to write the dataset dump
	collectorDatasetOutputBasename = "colly_github_%d"
	// collectorDatasetOutputFormat sets the ouput format of the dataset extracted by the collector
	// `OnTAB` event Export formats supported:
	//  - JSON (Sets + Books)
	//  - YAML (Sets + Books)
	//  - XLSX (Sets + Books)
	//  - XML (Sets + Books)
	//  - TSV (Sets)
	//  - CSV (Sets)
	//  - ASCII + Markdown (Sets)
	//  - MySQL (Sets)
	//  - Postgres (Sets)
	// `OnTAB` event loading formats supported:
	//  - JSON (Sets + Books)
	//  - YAML (Sets + Books)
	//  - XML (Sets)
	//  - CSV (Sets)
	//  - TSV (Sets)
	collectorDatasetOutputFormat = "tabular-grid"
	//  collectorSubDatasetColumns specifies the columns to filter from the json content
	collectorSubDatasetColumns = []string{"id", "name", "full_name", "description", "language", "stargazers_count", "forks_count"}
)

// AppendDynamicColumn to the tabular dataset
func addFreq(row []interface{}) interface{} {
	if row == nil {
		return 0
	}
	return len(row[2].(string))
}

// PrettyPrint structs
func prettyPrint(msg ...interface{}) {
	pp.Println(msg...)
}

func init() {
	// Ensure that the output format is set in lower case.
	collectorDatasetOutputFormat = strings.ToLower(collectorDatasetOutputFormat)
}

func main() {

	// Instantiate default collector
	c := colly.NewCollector(
		// Visit only domains: api.github.com
		colly.AllowedDomains("api.github.com"),
		colly.AllowTabular(collectorTabEnabled),
	)

	// UseJsonParser. Available: mxj, gjson, json (default)
	c.UseJsonParser = collectorJsonParser

	if collectorDebug {
		c.SetDebugger(collectorDebugger)
	}

	// On every a element which has tabular format data call callback
	// Notes:
	// - `OnTAB` callback event are enabled only if the `AllowTabular` attribute is set to true.
	// - `OnTAB` use a fork of the package `github.com/agrison/go-tablib`
	// - `OnTAB` query specifications are available in 'SPECS.md'
	c.OnTAB("0:0", func(e *colly.TABElement) {

		// Debug the dataset slice
		if appDebug {
			fmt.Println("Valid=", e.Dataset.Valid(), "Height=", e.Dataset.Height(), "Width=", e.Dataset.Width())
			pp.Printf("Headers: \n %s \n\n", e.Dataset.Headers())
		}

		// Select sub-dataset
		ds, err := e.Dataset.Select(0, 0, "column_1")
		if err != nil {
			fmt.Println("error:", err)
		}

		// Update dataset
		// Add a dynamic column, by passing a function which has access to the current row, and must return a value:
		ds.AppendDynamicColumn("changefreq", addFreq)
		ds.AppendDynamicColumn("priority", addFreq)
		// ds.AppendDynamicColumn("freq", addFreq)

		// Export dataset
		// ds.EXPORT_FORMAT().String() 					--> returns the contents of the exported dataset as a string.
		// ds.EXPORT_FORMAT().Bytes() 					--> returns the contents of the exported dataset as a byte array.
		// ds.EXPORT_FORMAT().WriteTo(writer) 			--> writes the exported dataset to w.
		// ds.EXPORT_FORMAT().WriteFile(filename, perm) --> writes the databook or dataset content to a file named by filename.
		var output string
		switch collectorDatasetOutputFormat {
		// YAML
		case "yaml":
			if export, err := ds.YAML(); err == nil {
				output = export.String()
			} else {
				fmt.Println("error:", err)
			}

		// JSON
		case "json":
			if export, err := ds.JSON(); err == nil {
				output = export.String()
			} else {
				fmt.Println("error:", err)
			}

		// TSV
		case "tsv":
			if export, err := ds.TSV(); err == nil {
				output = export.String()
			} else {
				fmt.Println("error:", err)
			}

		// CSV
		case "csv":
			if export, err := ds.CSV(); err == nil {
				output = export.String()
			} else {
				fmt.Println("error:", err)
			}

		// Markdown
		case "markdown", "tabular-markdown":
			output = ds.Markdown().String()

		// HTML
		case "html":
			output = ds.HTML().String()

		// MySQL
		case "mysql":
			output = ds.MySQL("github_starred").String()

		// Postgres
		case "postgresql":
			output = ds.Postgres("github_starred").String()

		// ASCII - TabularGrid
		case "grid-default", "ascii-grid", "tabular-grid":
			output = ds.Tabular("grid" /* tablib.TabularGrid */).String()

		// ASCII - TabularSimple
		case "grid-simple", "ascii-simple", "tabular-simple":
			output = ds.Tabular("simple" /* tablib.TabularSimple */).String()

		// ASCII - TabularSiTabularCondensedmple
		case "grid-condensed", "ascii-condensed", "tabular-condensed":
			output = ds.Tabular("condensed" /* tablib.TabularCondensed */).String()

		}

		// output final export
		fmt.Println(output)

	})

	// On every a element which has json content-type or file extension call callback
	// Notes:
	// - If `AllowTabular` is true, OnJSON is overrided by the OnTAB callback event
	// - OnJSON use a fork of the package `github.com/antchfx/jsonquery`
	c.OnJSON("//description", func(e *colly.JSONElement) {
		if appDebug {
			fmt.Printf("Values found: %s\n", e.Text)
		}
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		if appDebug {
			fmt.Println("Visiting", r.URL.String())
		}
	})

	// Start scraping on https://golanglibs.org pages
	c.Visit(targetRootURL)
}
