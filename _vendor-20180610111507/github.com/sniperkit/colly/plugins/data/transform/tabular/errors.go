package tablib

import (
	"errors"
)

var (

	// ErrInvalidDimensions is returned when trying to append/insert too much
	// or not enough values to a row or column
	ErrInvalidDimensions = errors.New("tablib: Invalid dimension")

	// ErrInvalidColumnIndex is returned when trying to insert a column at an
	// invalid index
	ErrInvalidColumnIndex = errors.New("tablib: Invalid column index")

	// ErrInvalidRowIndex is returned when trying to insert a row at an
	// invalid index
	ErrInvalidRowIndex = errors.New("tablib: Invalid row index")

	// ErrInvalidDataset is returned when trying to validate a Dataset against
	// the constraints that have been set on its columns.
	ErrInvalidDataset = errors.New("tablib: Invalid dataset")

	// ErrInvalidTag is returned when trying to add a tag which is not a string.
	ErrInvalidTag = errors.New("tablib: A tag must be a string")

	// ErrExportFormatNotSupported is thrown when the export format requested is not supported by tablib
	ErrExportFormatNotSupported = errors.New("Export format not supported yet, please choose one of")

	// ErrCompressFormatNotSupported is thrown when the compression format requested is not supported by tablib
	ErrCompressFormatNotSupported = errors.New("Compression format not supported yet, please choose one of")

	// ErrCompressionFormatNotSet is thrown when the compression format to delete is not set for the current dataset
	ErrCompressionFormatNotSet = errors.New("Compression format not exists for the current dataset.")

	// ErrUnmatchedKeys is thrown when unmarshalling a format failed on some keys
	ErrUnmatchedKeys = errors.New("Unmatched keys in your dataset.")

	// ErrUnmarshallingJsonWithGson is thrown when unmarshalling JSON format failed
	ErrUnmarshallingJsonWithGson = errors.New("Error while unmarhsalling a json source with `gson` package.")

	// ErrUnmarshallingJsonWithMxj is thrown when unmarshalling JSON format failed
	ErrUnmarshallingJsonWithMxj = errors.New("Error while unmarhsalling a json source with `mxj` package.")
)
