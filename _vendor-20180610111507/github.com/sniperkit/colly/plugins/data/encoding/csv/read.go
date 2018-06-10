package csv

import (
	"encoding/csv"
	"errors"
	"io"
	"os"
	"sync"
)

// Reader is iterable csv reader. The iteration is performed in its Iterate() method, which
// may only be invoked once per each instance of the Reader.
type Reader struct {
	source                       maker
	delimiter, comment           rune
	numFields                    int
	lazyQuotes, trimLeadingSpace bool
	header                       map[string]int
	headerFromFirstRow           bool
	lock                         *sync.RWMutex
	wg                           *sync.WaitGroup
}

type maker = func() (io.Reader, func(), error)

// FromReader constructs a new csv reader from the given io.Reader, with default settings.
func FromReader(input io.Reader) *Reader {
	return makeReader(func() (io.Reader, func(), error) {
		return input, func() {}, nil
	})
}

// FromReadCloser constructs a new csv reader from the given io.ReadCloser, with default settings.
func FromReadCloser(input io.ReadCloser) *Reader {
	return makeReader(func() (io.Reader, func(), error) {
		return input, func() { input.Close() }, nil
	})
}

// FromFile constructs a new csv reader bound to the specified file, with default settings.
func FromFile(name string) *Reader {
	return makeReader(func() (io.Reader, func(), error) {
		file, err := os.Open(name)

		if err != nil {
			return nil, nil, err
		}

		return file, func() { file.Close() }, nil
	})
}

func makeReader(fn maker) *Reader {
	return &Reader{
		source:             fn,
		delimiter:          ',',
		headerFromFirstRow: true,
		lock:               &sync.RWMutex{},
		wg:                 &sync.WaitGroup{},
	}
}

// Delimiter sets the symbol to be used as a field delimiter.
func (r *Reader) Delimiter(c rune) *Reader {
	r.lock.RWLock()
	r.delimiter = c
	r.lock.RWUnlock()
	return r
}

// CommentChar sets the symbol that starts a comment.
func (r *Reader) CommentChar(c rune) *Reader {
	r.lock.RWLock()
	r.comment = c
	r.lock.RWUnlock()
	return r
}

// LazyQuotes specifies that a quote may appear in an unquoted field and a
// non-doubled quote may appear in a quoted field of the input.
func (r *Reader) LazyQuotes() *Reader {
	r.lock.RWLock()
	r.lazyQuotes = true
	r.lock.RWUnlock()
	return r
}

// TrimLeadingSpace specifies that the leading white space in a field should be ignored.
func (r *Reader) TrimLeadingSpace() *Reader {
	r.lock.RWLock()
	r.trimLeadingSpace = true
	r.lock.RWUnlock()
	return r
}

// AssumeHeader sets the header for those input sources that do not have their column
// names specified on the first row. The header specification is a map
// from the assigned column names to their corresponding column indices.
func (r *Reader) AssumeHeader(spec map[string]int) *Reader {
	r.lock.RWLock()
	if len(spec) == 0 {
		// panic("Empty header spec")
		panic(errEmptyHeaderSpec)
	}

	for name, col := range spec {
		if col < 0 {
			panic("Header spec: Negative index for column " + name)
		}
	}

	r.header = spec
	r.headerFromFirstRow = false
	r.lock.RWUnlock()
	return r
}

// ExpectHeader sets the header for input sources that have their column
// names specified on the first row. The row gets verified
// against this specification when the reading starts.
// The header specification is a map from the expected column names to their corresponding
// column indices. A negative value for an index means that the real value of the index
// will be found by searching the first row for the specified column name.
func (r *Reader) ExpectHeader(spec map[string]int) *Reader {
	if len(spec) == 0 {
		// panic("Empty header spec")
		panic(errEmptyHeaderSpec)
	}

	r.header = make(map[string]int, len(spec))

	for name, col := range spec {
		r.header[name] = col
	}

	r.headerFromFirstRow = true
	return r
}

// SelectColumns specifies the names of the columns to read from the input source.
// The header specification is built by searching the first row of the input
// for the names specified and recording the indices of those columns. It is an error
// if any column name is not found.
func (r *Reader) SelectColumns(names ...string) *Reader {
	if len(names) == 0 {
		// panic("Empty header spec")
		panic(errEmptyHeaderSpec)
	}

	r.header = make(map[string]int, len(names))

	for _, name := range names {
		if _, found := r.header[name]; found {
			panic("Header spec: Duplicate column name: " + name)
		}

		r.header[name] = -1
	}

	r.headerFromFirstRow = true
	return r
}

// NumFields sets the expected number of fields on each row of the input.
// It is an error if any row does not have this exact number of fields.
func (r *Reader) NumFields(n int) *Reader {
	r.numFields = n
	return r
}

// NumFieldsAuto specifies that the number of fields on each row must match that of
// the first row of the input.
func (r *Reader) NumFieldsAuto() *Reader {
	return r.NumFields(0)
}

// NumFieldsAny specifies that each row of the input may have different number
// of fields. Rows shorter than the maximum column index in the header specification will be padded
// with empty fields.
func (r *Reader) NumFieldsAny() *Reader {
	return r.NumFields(-1)
}

// Iterate reads the input row by row, converts each line to the Row type, and calls
// the supplied RowFunc.
func (r *Reader) Iterate(fn RowFunc) error {
	// source
	input, close, err := r.source()

	if err != nil {
		return err
	}

	defer close()

	// csv.Reader
	reader := csv.NewReader(input)
	reader.Comma = r.delimiter
	reader.Comment = r.comment
	reader.LazyQuotes = r.lazyQuotes
	reader.TrimLeadingSpace = r.trimLeadingSpace
	reader.FieldsPerRecord = r.numFields

	// header
	var header map[string]int

	lineNo := uint64(1)

	if r.headerFromFirstRow {
		if header, err = r.makeHeader(reader); err != nil {
			return mapError(err, lineNo)
		}

		lineNo++
	} else {
		header = r.header
	}

	// iteration
	var line []string

	for line, err = reader.Read(); err == nil; line, err = reader.Read() {
		row := make(map[string]string, len(header))

		for name, index := range header {
			if index < len(line) {
				row[name] = line[index]
			} else if r.numFields < 0 { // padding allowed
				row[name] = ""
			} else {
				return &DataSourceError{
					Line: lineNo,
					Err:  fmt.Errorf("Column not found: %q (%d)", name, index),
				}
			}
		}

		if err = fn(row); err != nil {
			break
		}

		lineNo++
	}

	// map error
	if err != io.EOF {
		return mapError(err, lineNo)
	}

	return nil
}

// build header spec from the first row of the input file
func (r *Reader) makeHeader(reader *csv.Reader) (map[string]int, error) {
	line, err := reader.Read()

	if err != nil {
		return nil, err
	}

	if len(line) == 0 {
		return nil, errEmptyHeader
	}

	if len(r.header) == 0 { // the header is not specified - get it from the first line
		header := make(map[string]int, len(line))

		for i, name := range line {
			header[name] = i
		}

		return header, nil
	}

	// check and update the specified header
	header := make(map[string]int, len(r.header))

	// fix column indices
	for i, name := range line {
		if index, found := r.header[name]; found {
			if index == -1 || index == i {
				header[name] = i
			} else {
				return nil, fmt.Errorf(`Misplaced column %q: expected at pos. %d, but found at pos. %d`, name, index, i)
			}
		}
	}

	// check if all columns are found
	if len(header) < len(r.header) {
		// compose the list of the missing columns
		var list []string

		for name := range r.header {
			if _, found := header[name]; !found {
				list = append(list, name)
			}
		}

		// return error message
		if len(list) > 1 {
			return nil, errors.New("Columns not found: " + strings.Join(list, ", "))
		}

		return nil, errors.New("Column not found: " + list[0])
	}

	// all done
	return header, nil
}
