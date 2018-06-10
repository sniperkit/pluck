package tablib

import (
	"bytes"
	"encoding/binary"
	"io"
	"io/ioutil"
	"os"
	"sync"
	// helpers
	// humanize "github.com/sniperkit/colly/plugins/data/transform/humanize"
)

const defaultBufferCap = 16 * 1024

var (
	AllowedExportFormats []string = []string{"yaml", "json", "xlsx", "csv", "tsv", "txt", "ascii", "html", "markdown", "sql", "mysql", "postgres", "toml"}
)

// ExportFormat represents a export data format
type ExportFormat string

const (
	YAML     ExportFormat = "yaml"     // YAML format
	JSON     ExportFormat = "json"     // JSON format
	TOML     ExportFormat = "toml"     // TOML format
	XML      ExportFormat = "xml"      // XML format
	XLSX     ExportFormat = "xlsx"     // XLSX format
	CSV      ExportFormat = "csv"      // CSV format
	TSV      ExportFormat = "tsv"      // TSV format
	TXT      ExportFormat = "txt"      // TXT format
	ASCII    ExportFormat = "ascii"    // ASCII format
	HTML     ExportFormat = "html"     // HTML format
	MARKDOWN ExportFormat = "markdown" // MARKDOWN format
	MYSQL    ExportFormat = "mysql"    // MYSQL format
	POSTGRES ExportFormat = "postgres" // POSTGRES format
)

// Export represents an exportable dataset, it cannot be manipulated at this point
// and it can just be converted to a string, []byte or written to a io.Writer.
// The exportable struct just holds a bytes.Buffer that is used by the tablib library
// to write export formats content. Real work is delegated to bytes.Buffer.
type Export struct {
	buffer *bytes.Buffer
	lock   *sync.RWMutex
	wg     *sync.WaitGroup
}

// newBuffer returns a new bytes.Buffer instance already initialized
// with an underlying bytes array of the capacity equal to defaultBufferCap.
func newBuffer() *bytes.Buffer {
	return newBufferWithCap(defaultBufferCap)
}

// newBufferWithCap returns a new bytes.Buffer instance already initialized
// with an underlying bytes array of the given capacity.
func newBufferWithCap(initialCap int) *bytes.Buffer {
	initialBuf := make([]byte, 0, initialCap)
	return bytes.NewBuffer(initialBuf)
}

// newExport creates a new instance of Export from a bytes.Buffer.
func newExport(buffer *bytes.Buffer) *Export {
	return &Export{
		buffer: buffer,
		lock:   &sync.RWMutex{},
		wg:     &sync.WaitGroup{},
	}
}

// newExport creates a new instance of Export from a byte array.
func newExportFromBytes(buf []byte) *Export {
	return &Export{
		buffer: bytes.NewBuffer(buf),
		lock:   &sync.RWMutex{},
		wg:     &sync.WaitGroup{},
	}
}

// newExportFromString creates a new instance of Export from a string.
func newExportFromString(str string) *Export {
	buff := newBufferWithCap(len(str))
	buff.WriteString(str)
	return newExport(buff)
}

// Size returns the size of the exported dataset as a byte array.
func (e *Export) Size(outputType string) (length int) {
	e.lock.RLock()
	defer e.lock.RUnlock()

	switch outputType {

	case "string":
		length = len(e.buffer.String())
	case "bytes":
		fallthrough
	default:
		length = binary.Size(e.buffer.Bytes())
	}

	return length
}

// Bytes returns the contents of the exported dataset as a byte array.
func (e *Export) Bytes() []byte {
	e.lock.RLock()
	defer e.lock.RUnlock()

	return e.buffer.Bytes()
}

// String returns the contents of the exported dataset as a string.
func (e *Export) String() string {
	e.lock.RLock()
	defer e.lock.RUnlock()

	return e.buffer.String()
}

// Interface returns the contents of the exported dataset as an interface.
func (e *Export) Interface() interface{} {
	e.lock.RLock()
	defer e.lock.RUnlock()

	return e.buffer.String()
}

// WriteTo writes the exported dataset to w.
func (e *Export) WriteTo(w io.Writer) (int64, error) {
	e.lock.RLock()
	defer e.lock.RUnlock()

	return e.buffer.WriteTo(w)
}

// WriteFile writes the databook or dataset content to a file named by filename.
// If the file does not exist, WriteFile creates it with permissions perm;
// otherwise WriteFile truncates it before writing.
func (e *Export) WriteFile(filename string, perm os.FileMode) error {
	e.lock.RLock()
	defer e.lock.RUnlock()

	return ioutil.WriteFile(filename, e.Bytes(), perm)
}
