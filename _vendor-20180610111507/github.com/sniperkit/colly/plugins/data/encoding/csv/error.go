package csv

import (
	"errors"
	"fmt"
)

var (
	errEmptyHeaderSpec    = errors.New("Empty header spec")
	errEmptyHeader        = errors.New("Empty header")
	errEmptyMatchFuncLike = errors.New("Empty match function in Like() predicate")
)

// DataSourceError is the type of the error returned from Reader.Iterate method.
type DataSourceError struct {
	Line uint64 // counting from 1
	Err  error
}

// Error returns a human-readable error message string.
func (e *DataSourceError) Error() string {
	return fmt.Sprintf(`Row %d: %s`, e.Line, e.Err)
}

// annotate error with line number
func mapError(err error, lineNo uint64) error {
	switch e := err.(type) {
	case *csv.ParseError:
		return &DataSourceError{
			Line: lineNo,
			Err:  e.Err,
		}
	case *os.PathError:
		return &DataSourceError{
			Line: lineNo,
			Err:  errors.New(e.Op + ": " + e.Err.Error()),
		}
	default:
		return &DataSourceError{
			Line: lineNo,
			Err:  err,
		}
	}
}
