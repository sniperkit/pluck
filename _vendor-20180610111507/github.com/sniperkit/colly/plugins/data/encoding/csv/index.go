package csv

import (
	"errors"
	"os"
)

// Index is a sorted collection of Rows with O(log(n)) complexity of search
// on the indexed columns. Iteration over the Index yields a sequence of Rows sorted on the index.
type Index struct {
	impl indexImpl
}

// Iterate iterates over all rows of the index. The rows are sorted by the values of the columns
// specified when the Index was created.
func (index *Index) Iterate(fn RowFunc) error {
	return iterate(index.impl.rows, fn)
}

// Find returns a DataSource of all Rows from the Index that match the specified values
// in the indexed columns, left to the right. The number of specified values may be less than
// the number of the indexed columns.
func (index *Index) Find(values ...string) DataSource {
	return TakeRows(index.impl.find(values))
}

// SubIndex returns an Index containing only the rows where the values of the
// indexed columns match the supplied values, left to the right. The number of specified values
// must be less than the number of indexed columns.
func (index *Index) SubIndex(values ...string) *Index {
	if len(values) >= len(index.impl.columns) {
		panic("Too many values in SubIndex()")
	}

	return &Index{indexImpl{
		rows:    index.impl.find(values),
		columns: index.impl.columns[len(values):],
	}}
}

// ResolveDuplicates calls the specified function once per each pack of duplicates with the same key.
// The specified function must not modify its parameter and is expected to do one of the following:
//
// - Select and return one row from the input list. The row will be used as the only row with its key;
//
// - Return an empty row. The entire set of rows will be ignored;
//
// - Return an error which will be passed back to the caller of ResolveDuplicates().
func (index *Index) ResolveDuplicates(resolve func(rows []Row) (Row, error)) error {
	return index.impl.dedup(resolve)
}

// WriteTo writes the index to the specified file.
func (index *Index) WriteTo(fileName string) (err error) {
	var file *os.File

	if file, err = os.Create(fileName); err != nil {
		return
	}

	defer func() {
		if e := file.Close(); e != nil || err != nil {
			os.Remove(fileName)

			if err == nil {
				err = e
			}
		}
	}()

	enc := gob.NewEncoder(file)

	if err = enc.Encode(index.impl.columns); err == nil {
		err = enc.Encode(index.impl.rows)
	}

	return
}

// LoadIndex reads the index from the specified file.
func LoadIndex(fileName string) (*Index, error) {
	var file *os.File
	var err error

	if file, err = os.Open(fileName); err != nil {
		return nil, err
	}

	defer file.Close()

	index := new(Index)
	dec := gob.NewDecoder(file)

	if err = dec.Decode(&index.impl.columns); err != nil {
		return nil, err
	}

	if err = dec.Decode(&index.impl.rows); err != nil {
		return nil, err
	}

	return index, nil
}

func createIndex(src DataSource, columns []string) (*Index, error) {
	switch len(columns) {
	case 0:
		panic("Empty column list in CreateIndex()")
	case 1:
		// do nothing
	default:
		if !allColumnsUnique(columns) {
			panic("Duplicate column name(s) in CreateIndex()")
		}
	}

	index := &Index{indexImpl{columns: columns}}

	// copy Rows with validation
	if err := src(func(row Row) error {
		for _, col := range columns {
			if !row.HasColumn(col) {
				return fmt.Errorf(`Missing column %q while creating an index`, col)
			}
		}

		index.impl.rows = append(index.impl.rows, row)
		return nil
	}); err != nil {
		return nil, err
	}

	// sort
	sort.Sort(&index.impl)
	return index, nil
}

func createUniqueIndex(src DataSource, columns []string) (index *Index, err error) {
	// create index
	if index, err = createIndex(src, columns); err != nil || len(index.impl.rows) < 2 {
		return
	}

	// check for duplicates by linear search; not the best idea.
	rows := index.impl.rows

	for i := 1; i < len(rows); i++ {
		if equalRows(columns, rows[i-1], rows[i]) {
			return nil, errors.New("Duplicate value while creating unique index: " + rows[i].SelectExisting(columns...).String())
		}
	}

	return
}

// index implementation
type indexImpl struct {
	rows    []Row
	columns []string
}

// functions required by sort.Sort()
func (index *indexImpl) Len() int      { return len(index.rows) }
func (index *indexImpl) Swap(i, j int) { index.rows[i], index.rows[j] = index.rows[j], index.rows[i] }

func (index *indexImpl) Less(i, j int) bool {
	left, right := index.rows[i], index.rows[j]

	for _, col := range index.columns {
		switch strings.Compare(left[col], right[col]) {
		case -1:
			return true
		case 1:
			return false
		}
	}

	return false
}

// deduplication
func (index *indexImpl) dedup(resolve func(rows []Row) (Row, error)) (err error) {
	// find first duplicate
	var lower int

	for lower = 1; lower < len(index.rows); lower++ {
		if equalRows(index.columns, index.rows[lower-1], index.rows[lower]) {
			break
		}
	}

	if lower == len(index.rows) {
		return
	}

	dest := lower - 1

	// loop: find index of the first row with another key, resolve, then find next duplicate
	for lower < len(index.rows) {
		// the duplicate is in [lower-1, upper[ range
		values, _ := index.rows[lower].SelectValues(index.columns...)

		upper := lower + sort.Search(len(index.rows)-lower, func(i int) bool {
			return index.cmp(lower+i, values, false)
		})

		// resolve
		var row Row

		if row, err = resolve(index.rows[lower-1 : upper]); err != nil {
			return
		}

		lower = upper + 1

		// store the chosen row if not 'empty'
		if len(row) >= len(index.columns) {
			index.rows[dest] = row
			dest++
		}

		// find next duplicate
		for lower < len(index.rows) {
			if equalRows(index.columns, index.rows[lower-1], index.rows[lower]) {
				break
			}

			index.rows[dest] = index.rows[lower-1]
			lower++
			dest++
		}
	}

	if err == nil {
		index.rows = index.rows[:dest]
	}

	return
}

// search on the index
func (index *indexImpl) find(values []string) []Row {
	// check boundaries
	if len(values) == 0 {
		return index.rows
	}

	if len(values) > len(index.columns) {
		panic("Too many columns in indexImpl.find()")
	}

	// get bounds
	upper := sort.Search(len(index.rows), func(i int) bool {
		return index.cmp(i, values, false)
	})

	lower := sort.Search(upper, func(i int) bool {
		return index.cmp(i, values, true)
	})

	// done
	return index.rows[lower:upper]
}

func (index *indexImpl) first(values []string) int {
	return sort.Search(len(index.rows), func(i int) bool {
		return index.cmp(i, values, true)
	})
}

func (index *indexImpl) has(values []string) bool {
	// find the lowest index
	i := index.first(values)

	// check if the row at the lowest index matches the values
	return i < len(index.rows) && !index.cmp(i, values, false)
}

func (index *indexImpl) cmp(i int, values []string, eq bool) bool {
	row := index.rows[i]

	for j, val := range values {
		switch strings.Compare(row[index.columns[j]], val) {
		case 1:
			return true
		case -1:
			return false
		}
	}

	return eq
}
