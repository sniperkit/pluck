package csv

import (
	"encoding/csv"
)

// DataSource is the iterator type used throughout this library. The iterator
// calls the given RowFunc once per each row. The iteration continues until
// either the data source is exhausted or the supplied RowFunc returns a non-nil error, in
// which case the error is returned back to the caller of the iterator. A special case of io.EOF simply
// stops the iteration and the iterator function returns nil error.
type DataSource func(RowFunc) error

// TakeRows converts a slice of Rows to a DataSource.
func TakeRows(rows []Row) DataSource {
	return func(fn RowFunc) error {
		return iterate(rows, fn)
	}
}

// the core iteration
func iterate(rows []Row, fn RowFunc) (err error) {
	var row Row
	var i int

	for i, row = range rows {
		if err = fn(row.Clone()); err != nil {
			break
		}
	}

	switch err {
	case nil:
		// nothing to do
	case io.EOF:
		err = nil // end of iteration
	default:
		// wrap error
		err = &DataSourceError{
			Line: uint64(i),
			Err:  err,
		}
	}

	return
}

// Take converts anything with Iterate() method to a DataSource.
func Take(src interface {
	Iterate(fn RowFunc) error
}) DataSource {
	return src.Iterate
}

// Transform is the most generic operation on a Row. It takes a function that
// maps a Row to another Row and an error. Any error returned from that function
// stops the iteration, otherwise the returned Row, if not empty, gets passed
// down to the next stage of the processing pipeline.
func (src DataSource) Transform(trans func(Row) (Row, error)) DataSource {
	return func(fn RowFunc) error {
		return src(func(row Row) (err error) {
			if row, err = trans(row); err == nil && len(row) > 0 {
				err = fn(row)
			}

			return
		})
	}
}

// Filter takes a predicate which, when applied to a Row, decides if that Row
// should be passed down for further processing. The predicate should return 'true' to pass the Row.
func (src DataSource) Filter(pred func(Row) bool) DataSource {
	return func(fn RowFunc) error {
		return src(func(row Row) (err error) {
			if pred(row) {
				err = fn(row)
			}

			return
		})
	}
}

// Map takes a function which gets applied to each Row when the source is iterated over. The function
// may return a modified input Row, or an entirely new Row.
func (src DataSource) Map(mf func(Row) Row) DataSource {
	return func(fn RowFunc) error {
		return src(func(row Row) error {
			return fn(mf(row))
		})
	}
}

// Validate takes a function which checks every Row upon iteration and returns an error
// if the validation fails. The iteration stops at the first error encountered.
func (src DataSource) Validate(vf func(Row) error) DataSource {
	return func(fn RowFunc) error {
		return src(func(row Row) (err error) {
			if err = vf(row); err == nil {
				err = fn(row)
			}

			return
		})
	}
}

// Top specifies the number of Rows to pass down the pipeline before stopping the iteration.
func (src DataSource) Top(n uint64) DataSource {
	return func(fn RowFunc) error {
		counter := n

		return src(func(row Row) error {
			if counter == 0 {
				return io.EOF
			}

			counter--
			return fn(row)
		})
	}
}

// Drop specifies the number of Rows to ignore before passing the remaining rows down the pipeline.
func (src DataSource) Drop(n uint64) DataSource {
	return func(fn RowFunc) error {
		counter := n

		return src(func(row Row) error {
			if counter == 0 {
				return fn(row)
			}

			counter--
			return nil
		})
	}
}

// TakeWhile takes a predicate which gets applied to each Row upon iteration.
// The iteration stops when the predicate returns 'false'.
func (src DataSource) TakeWhile(pred func(Row) bool) DataSource {
	return func(fn RowFunc) error {
		var done bool

		return src(func(row Row) error {
			if done = (done || !pred(row)); done {
				return io.EOF
			}

			return fn(row)
		})
	}
}

// DropWhile ignores all the Rows for as long as the specified predicate is true;
// afterwards all the remaining Rows are passed down the pipeline.
func (src DataSource) DropWhile(pred func(Row) bool) DataSource {
	return func(fn RowFunc) error {
		var yield bool

		return src(func(row Row) (err error) {
			if yield = (yield || !pred(row)); yield {
				err = fn(row)
			}

			return
		})
	}
}

// DropColumns removes the specifies columns from each row.
func (src DataSource) DropColumns(columns ...string) DataSource {
	if len(columns) == 0 {
		panic("No columns specified in DropColumns()")
	}

	return func(fn RowFunc) error {
		return src(func(row Row) error {
			for _, col := range columns {
				delete(row, col)
			}

			return fn(row)
		})
	}
}

// SelectColumns leaves only the specified columns on each row. It is an error
// if any such column does not exist.
func (src DataSource) SelectColumns(columns ...string) DataSource {
	if len(columns) == 0 {
		panic("No columns specified in SelectColumns()")
	}

	return func(fn RowFunc) error {
		return src(func(row Row) (err error) {
			if row, err = row.Select(columns...); err == nil {
				err = fn(row)
			}

			return
		})
	}
}

// IndexOn iterates the input source, building index on the specified columns.
// Columns are taken from the specified list from left to the right.
func (src DataSource) IndexOn(columns ...string) (*Index, error) {
	return createIndex(src, columns)
}

// UniqueIndexOn iterates the input source, building unique index on the specified columns.
// Columns are taken from the specified list from left to the right.
func (src DataSource) UniqueIndexOn(columns ...string) (*Index, error) {
	return createUniqueIndex(src, columns)
}

// Join returns a DataSource which is a join between the current DataSource and the specified
// Index. The specified columns are matched against those from the index, in the order of specification.
// Empty 'columns' list yields a join on the columns from the Index (aka "natural join") which all must
// exist in the current DataSource.
// Each row in the resulting table contains all the columns from both the current table and the index.
// This is a lazy operation, the actual join is performed only when the resulting table is iterated over.
func (src DataSource) Join(index *Index, columns ...string) DataSource {
	if len(columns) == 0 {
		columns = index.impl.columns
	} else if len(columns) > len(index.impl.columns) {
		panic("Too many source columns in Join()")
	}

	return func(fn RowFunc) error {
		return src(func(row Row) (err error) {
			var values []string

			if values, err = row.SelectValues(columns...); err == nil {
				n := len(index.impl.rows)

				for i := index.impl.first(values); i < n && !index.impl.cmp(i, values, false); i++ {
					if err = fn(mergeRows(index.impl.rows[i], row)); err != nil {
						break
					}
				}
			}

			return
		})
	}
}

func mergeRows(left, right Row) Row {
	r := make(map[string]string, len(left)+len(right))

	for k, v := range left {
		r[k] = v
	}

	for k, v := range right {
		r[k] = v
	}

	return r
}

// Except returns a table containing all the rows not in the specified Index, unchanged. The specified
// columns are matched against those from the index, in the order of specification. If no columns
// are specified then the columns list is taken from the index.
func (src DataSource) Except(index *Index, columns ...string) DataSource {
	if len(columns) == 0 {
		columns = index.impl.columns
	} else if len(columns) > len(index.impl.columns) {
		panic("Too many source columns in Except()")
	}

	return func(fn RowFunc) error {
		return src(func(row Row) (err error) {
			var values []string

			if values, err = row.SelectValues(columns...); err == nil {
				if !index.impl.has(values) {
					err = fn(row)
				}
			}

			return
		})
	}
}
