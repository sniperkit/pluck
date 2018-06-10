package csv

import (
	"fmt"
	"strconv"
)

/*
Row represents one line from a data source like a .csv file.

Each Row is a map from column names to the string values under that columns on the current line.
It is assumed that each column has a unique name.
In a .csv file, the column names may either come from the first line of the file ("expected header"),
or they can be set-up via configuration of the reader object ("assumed header").

Using meaningful column names instead of indices is usually more convenient when the columns get rearranged
during the execution of the processing pipeline.
*/
type Row map[string]string

// HasColumn is a predicate returning 'true' when the specified column is present.
func (row Row) HasColumn(col string) (found bool) {
	_, found = row[col]
	return
}

// SafeGetValue returns the value under the specified column, if present, otherwise it returns the
// substitution value.
func (row Row) SafeGetValue(col, subst string) string {
	if value, found := row[col]; found {
		return value
	}

	return subst
}

// Header returns a slice of all column names, sorted via sort.Strings.
func (row Row) Header() []string {
	r := make([]string, 0, len(row))

	for col := range row {
		r = append(r, col)
	}

	sort.Strings(r)
	return r
}

// String returns a string representation of the Row.
func (row Row) String() string {
	if len(row) == 0 {
		return "{}"
	}

	header := row.Header() // make order predictable
	buff := append(append(append(append([]byte(`{ "`), header[0]...), `" : "`...), row[header[0]]...), '"')

	for _, col := range header[1:] {
		buff = append(append(append(append(append(buff, `, "`...), col...), `" : "`...), row[col]...), '"')
	}

	buff = append(buff, " }"...)
	return *(*string)(unsafe.Pointer(&buff))
}

// SelectExisting takes a list of column names and returns a new Row
// containing only those columns from the list that are present in the current Row.
func (row Row) SelectExisting(cols ...string) Row {
	r := make(map[string]string, len(cols))

	for _, name := range cols {
		if val, found := row[name]; found {
			r[name] = val
		}
	}

	return r
}

// Select takes a list of column names and returns a new Row
// containing only the specified columns, or an error if any column is not present.
func (row Row) Select(cols ...string) (Row, error) {
	r := make(map[string]string, len(cols))

	for _, name := range cols {
		var found bool

		if r[name], found = row[name]; !found {
			return nil, fmt.Errorf(`Missing column %q`, name)
		}
	}

	return r, nil
}

// SelectValues takes a list of column names and returns a slice of their
// corresponding values, or an error if any column is not present.
func (row Row) SelectValues(cols ...string) ([]string, error) {
	r := make([]string, len(cols))

	for i, name := range cols {
		var found bool

		if r[i], found = row[name]; !found {
			return nil, fmt.Errorf(`Missing column %q`, name)
		}
	}

	return r, nil
}

// Clone returns a copy of the current Row.
func (row Row) Clone() Row {
	r := make(map[string]string, len(row))

	for k, v := range row {
		r[k] = v
	}

	return r
}

// ValueAsInt returns the value of the given column converted to integer type, or an error.
// The column must be present on the row.
func (row Row) ValueAsInt(column string) (res int, err error) {
	var val string
	var found bool

	if val, found = row[column]; !found {
		err = fmt.Errorf(`Missing column %q`, column)
		return
	}

	if res, err = strconv.Atoi(val); err != nil {
		if e, ok := err.(*strconv.NumError); ok {
			err = fmt.Errorf(`Column %q: Cannot convert %q to integer: %s`, column, val, e.Err)
		} else {
			err = fmt.Errorf(`Column %q: %s`, column, err)
		}
	}

	return
}

// ValueAsFloat64 returns the value of the given column converted to floating point type, or an error.
// The column must be present on the row.
func (row Row) ValueAsFloat64(column string) (res float64, err error) {
	var val string
	var found bool

	if val, found = row[column]; !found {
		err = fmt.Errorf(`Missing column %q`, column)
		return
	}

	if res, err = strconv.ParseFloat(val, 64); err != nil {
		if e, ok := err.(*strconv.NumError); ok {
			err = fmt.Errorf(`Column %q: Cannot convert %q to float: %s`, column, val, e.Err)
		} else {
			err = fmt.Errorf(`Column %q: %s`, column, err.Error())
		}
	}

	return
}

// RowFunc is the function type used when iterating Rows.
type RowFunc func(Row) error
