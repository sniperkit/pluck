package csv

import (
	"bytes"
	"encoding/csv"
	"io"
	"os"
)

// call the given function with the file stream open for writing
func writeFile(name string, fn func(io.Writer) error) (err error) {
	var file *os.File

	if file, err = os.Create(name); err != nil {
		return
	}

	defer func() {
		if p := recover(); p != nil {
			file.Close()
			os.Remove(name)
			panic(p)
		}

		if e := file.Close(); e != nil && err == nil {
			err = e
		}

		if err != nil {
			os.Remove(name)
		}
	}()

	err = fn(file)
	return
}

// ToCsv iterates the data source and writes the selected columns in .csv format to the given io.Writer.
// The data are written in the "canonical" form with the header on the first line and with all the lines
// having the same number of fields, using default settings for the underlying csv.Writer.
func (src DataSource) ToCsv(out io.Writer, columns ...string) (err error) {
	if len(columns) == 0 {
		panic("Empty column list in ToCsv() function")
	}

	w := csv.NewWriter(out)

	// header
	if err = w.Write(columns); err == nil {
		// rows
		err = src(func(row Row) (e error) {
			var values []string

			if values, e = row.SelectValues(columns...); e == nil {
				e = w.Write(values)
			}

			return
		})
	}

	if err == nil {
		w.Flush()
		err = w.Error()
	}

	return
}

// ToCsvFile iterates the data source and writes the selected columns in .csv format to the given file.
// The data are written in the "canonical" form with the header on the first line and with all the lines
// having the same number of fields, using default settings for the underlying csv.Writer.
func (src DataSource) ToCsvFile(name string, columns ...string) error {
	return writeFile(name, func(file io.Writer) error {
		return src.ToCsv(file, columns...)
	})
}

// ToJSON iterates over the data source and writes all Rows to the given io.Writer in JSON format.
func (src DataSource) ToJSON(out io.Writer) (err error) {
	var buff bytes.Buffer

	buff.WriteByte('[')

	count := uint64(0)
	enc := json.NewEncoder(&buff)

	enc.SetIndent("", "")
	enc.SetEscapeHTML(false)

	err = src(func(row Row) (e error) {
		if count++; count != 1 {
			buff.WriteByte(',')
		}

		if e = enc.Encode(row); e == nil && buff.Len() > 10000 {
			_, e = buff.WriteTo(out)
		}

		return
	})

	if err == nil {
		buff.WriteByte(']')
		_, err = buff.WriteTo(out)
	}

	return
}

// ToJSONFile iterates over the data source and writes all Rows to the given file in JSON format.
func (src DataSource) ToJSONFile(name string) error {
	return writeFile(name, src.ToJSON)
}

// ToRows iterates the DataSource storing the result in a slice of Rows.
func (src DataSource) ToRows() (rows []Row, err error) {
	err = src(func(row Row) error {
		rows = append(rows, row)
		return nil
	})

	return
}
