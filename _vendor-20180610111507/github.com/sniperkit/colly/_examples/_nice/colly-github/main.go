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
	appDebug = false

	//-- End
)

// github api vars
var (

	// githubAPIAccount sets the github user name to request for its starred repositories list
	githubAPIAccount = "roscopecoltran"

	// githubAPIPaginationPage
	githubAPIPaginationPage = 1

	// githubAPIPaginationPerPage
	githubAPIPaginationPerPage = 10

	// githubAPIPaginationDirection
	githubAPIPaginationDirection = "desc"

	// githubAPIPaginationSort
	githubAPIPaginationSort = "updated"

	// githubAPIPaginationParams
	githubAPIPaginationParams = []string{
		fmt.Sprintf("page=%d", githubAPIPaginationPage),
		fmt.Sprintf("per_page=%d", githubAPIPaginationPerPage),
		fmt.Sprintf("direction=%s", githubAPIPaginationDirection),
		fmt.Sprintf("sort=%s", githubAPIPaginationSort),
	}

	// githubAPIEndpointURL
	githubAPIEndpointURL = fmt.Sprintf("https://api.github.com/users/%s/starred?%s", githubAPIAccount, strings.Join(githubAPIPaginationParams, "&"))

	//-- End
)

// collector vars
var (

	// collectorDebug sets collector's debugger
	collectorDebug = false

	// collectorDebugger stores the collector's log event listener
	collectorDebugger *debug.LogDebugger = &debug.LogDebugger{}

	// - collectorJsonParser sets the json parser package to unmarshal JSON responses.
	//   - Available parsers:
	//     - `JSON` default golang "encoding/json" package. Important: This parser does not extract/flatten nested object headers
	//     - `MXJ` decodes / encodes JSON or XML to/from map[string]interface{}; extracts values with dot-notation paths and wildcards.
	//     - `GJSON` (Not Ready Yet), decodes JSON document; performs one line retrieval, dot notation paths, iteration, and parsing json lines.
	collectorJsonParser = "mxj"

	// collectorTabEnabled specifies if the collector load and marshall content-types that are tabular compatible
	// - `OnTAB` supported loading formats:
	//   - JSON (Sets + Books)
	//   - YAML (Sets + Books)
	//   - TOML (Sets + Books)
	//   - XML (Sets)
	//   - CSV (Sets)
	//   - TSV (Sets)
	// - IMPORTANT:
	//   - input must be marshallable as a slice interfaces ([]interface).
	//   - map[string]interface will be converted to []map[string]interface or []interface
	collectorTabEnabled = true

	// collectorDatasetOutputPrefixPath specifies the prefix path for all saved dumps
	collectorDatasetOutputPrefixPath = "./shared/dataset"

	// collectorDatasetOutputBasename specifies the template to use to write the dataset dump
	collectorDatasetOutputBasename = "colly_github_%d"

	// collectorDatasetOutputFormat sets the ouput format of the dataset extracted by the collector
	// `OnTAB` event export/print supported formats:
	//  - JSON (Sets + Books)
	//  - YAML (Sets + Books)
	//  - XLSX (Sets + Books)
	//  - XML (Sets + Books)
	//  - TSV (Sets)
	//  - CSV (Sets)
	//  - ASCII + Markdown (Sets)
	//  - MySQL (Sets)
	//  - Postgres (Sets)
	collectorDatasetOutputFormat = "tabular-grid"

	//  collectorSubDatasetColumns specifies the columns to filter from the json content
	collectorSubDatasetColumns = []string{"id", "name", "full_name", "description", "language", "stargazers_count", "forks_count"}
)

// init() function is executed when the executable is started, before function main()
func init() {
	// Ensure that the output format is set in lower case.
	collectorDatasetOutputFormat = strings.ToLower(collectorDatasetOutputFormat)
}

// descriptionLength function is a callback used to append dynamic column to a `OnTab` event dataset
func descriptionLength(row []interface{}) interface{} {
	if row == nil {
		return 0
	}
	//if appDebug {
	//	fmt.Printf("\n----- Calculated the description length for row:\n")
	//	prettyPrint(row)
	//}
	if len(row) < 2 {
		return 0
	}
	return len(row[2].(string))
}

// prettyPrint wraps debug message with `github.com/k0kubun/pp` package functions.
func prettyPrint(output ...interface{}) {
	pp.Println(output...)
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
		ds, err := e.Dataset.Select(0, 0, "id", "full_name", "description", "language", "stargazers_count", "owner_login", "owner_id")
		if err != nil {
			fmt.Println("error:", err)
		}

		// Update dataset
		// Add a dynamic column, by passing a function which has access to the current row, and must return a value:
		ds.AppendDynamicColumn("description_length", descriptionLength)

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

	// Start scraping on https://api.github.com
	c.Visit(githubAPIEndpointURL)
}
