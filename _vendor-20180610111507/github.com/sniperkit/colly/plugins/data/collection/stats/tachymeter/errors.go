package tachymeter

import (
	"errors"
	"strings"
)

var (
	// Tachymeter instance related errors
	errTachymeterFailed   = errors.New("Tachymeter could not be started.")
	errTachymeterDataset  = errors.New("Tachymeter dataset unknown. Available: " + strings.Join(allowed_export_datasets, ",") + ".")
	errTachymeterEncoding = errors.New("Tachymeter datasets cannot be exported in this format. Available: " + strings.Join(allowed_export_formats, ",") + ".")

	// Metrics related errors
	errMetricsFailed = errors.New("Metrics have failed.")
	errMetricsFormat = errors.New("Metrics datasets cannot be exported in this format. Available: " + strings.Join(allowed_export_formats, ",") + ".")

	// Exportable related errors
	errExportFormatType = errors.New("Export format type is invalid. Available: " + strings.Join(allowed_export_formats, ",") + ".")
)
