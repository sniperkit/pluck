# Colly - Export format-agnostic tabular dataset

## Summary

In this example, we will build a crawler able to format and pre-process data extracted by the colly collector, and output them with a format-agnostic tabular dataset callback.

Export formats supported:

* JSON (Sets + Books)
* YAML (Sets + Books)
* XLSX (Sets + Books)
* XML (Sets + Books)
* TSV (Sets)
* CSV (Sets)
* ASCII + Markdown (Sets)
* MySQL (Sets)
* Postgres (Sets)

Loading formats supported:

* JSON (Sets + Books)
* YAML (Sets + Books)
* XML (Sets)
* CSV (Sets)
* TSV (Sets)

## Overview

### e.Dataset
A Dataset is a table of tabular data. It must have a header row. Datasets can be exported to JSON, YAML, TOML, CSV, TSV, and XML. They can be filtered, sorted and validated against constraint on columns.

### e.Databook
A Databook is a set of Datasets. The most common form of a Databook is an Excel file with multiple spreadsheets. Databooks can be exported to JSON, YAML, TOML and XML.

### e.Exportable
An exportable is a struct that holds a buffer representing the Databook or Dataset after it has been formated to any of the supported export formats.
At this point the Datbook or Dataset cannot be modified anymore, but it can be returned as a `string`, a `[]byte` or written to a `io.Writer` or a file.

