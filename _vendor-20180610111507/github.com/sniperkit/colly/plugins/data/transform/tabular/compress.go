package tablib

import (
	"bytes"
	"sync"
	// parallel_gzip "github.com/klauspost/pgzip"
	// compress "github.com/klauspost/compress"
)

// public vars
var IsFaultTolerant bool = true

// private vaes
var allowedCompressFormats []string = []string{"zip", "gz", "tar.gz", "tar", "rar", "lz"}

// CompressFormat represents a compression data format
type CompressFormat string

const (
	ZIP    CompressFormat = "zip"    // ZIP format
	GZ     CompressFormat = "gzip"   // GZ format
	TAR_GZ CompressFormat = "tar.gz" // ETCD format
	TAR    CompressFormat = "tar"    // TAR format
	RAR    CompressFormat = "rar"    // RAR format
	LZ     CompressFormat = "lz"     // LZ format
)

// Export represents an exportable dataset, it cannot be manipulated at this point
// and it can just be converted to a string, []byte or written to a io.Writer.
// The exportable struct just holds a bytes.Buffer that is used by the tablib library
// to write export formats content. Real work is delegated to bytes.Buffer.
type Compressable struct {
	Formats       map[string]*CompressFormatConfig `json:"formats" yaml:"formats" toml:"formats" xml:"formats" ini:"formats"`
	FaultTolerant bool                             `default:"true" json:"fault_tolerant" yaml:"fault_tolerant" toml:"fault_tolerant" xml:"faultTolerant" ini:"faultTolerant"`
	Errors        []string                         `json:"-" yaml:"-" toml:"-" xml:"-" ini:"-"`
	buffer        *bytes.Buffer
	lock          *sync.RWMutex
	wg            *sync.WaitGroup
}

func NewCompressable() *Compressable {
	c := &Compressable{
		Formats:       make(map[string]*CompressFormatConfig, len(allowedCompressFormats)),
		FaultTolerant: IsFaultTolerant,
		lock:          &sync.RWMutex{},
		wg:            &sync.WaitGroup{},
	}
	return c
}

type CompressFormatConfig struct {
	Enabled       bool   `default:"true" json:"enabled" yaml:"enabled" toml:"enabled" xml:"enabled" ini:"enabled"`
	Format        string `required:"true" json:"engine" yaml:"engine" toml:"engine" xml:"engine" ini:"engine"`
	Level         int    `json:"level" yaml:"level" toml:"level" xml:"level" ini:"level"`
	SplitAt       int    `default:"false" json:"split_at" yaml:"split_at" toml:"split_at" xml:"split_at" ini:"split_at"`
	SplitPrefix   string `default:"part_" json:"split_prefix" yaml:"split_prefix" toml:"split_prefix" xml:"splitPrefix" ini:"splitPrefix"`
	SplitSuffix   string `default:"_%d" json:"split_suffix" yaml:"split_suffix" toml:"split_suffix" xml:"splitSuffix" ini:"splitSuffix"`
	FileExtension string `default:"true" json:"file_extension" yaml:"file_extension" toml:"file_extension" xml:"fileExtension" ini:"fileExtension"`
}

func (c *Compressable) Format(format string, status bool, level int, splitAt int, splitPrefix string, splitSuffix string, faultTolerant bool) *Compressable {
	if ok := inArray(format, allowedCompressFormats); !ok {
		if !c.FaultTolerant {
			panic(ErrCompressFormatNotSupported)
		}
		c.Errors = append(c.Errors, ErrCompressFormatNotSupported.Error())
		return c
	}

	cfg := &CompressFormatConfig{
		Enabled:     status,
		Format:      format,
		Level:       level,
		SplitAt:     splitAt,
		SplitPrefix: splitPrefix,
		SplitSuffix: splitSuffix,
	}

	c.lock.Lock()
	c.Formats[cfg.Format] = cfg
	c.lock.Unlock()

	return c
}

func (c *Compressable) FormatWithConfig(cfg *CompressFormatConfig) *Compressable {
	if ok := inArray(cfg.Format, allowedCompressFormats); !ok {
		if !c.FaultTolerant {
			panic(ErrCompressFormatNotSupported)
		}
		c.Errors = append(c.Errors, ErrCompressFormatNotSupported.Error())
		return c
	}

	format := &CompressFormatConfig{
		Enabled:     cfg.Enabled,
		Format:      cfg.Format,
		Level:       cfg.Level,
		SplitAt:     cfg.SplitAt,
		SplitPrefix: cfg.SplitPrefix,
		SplitSuffix: cfg.SplitSuffix,
	}

	c.lock.Lock()
	c.Formats[cfg.Format] = format
	c.lock.Unlock()

	return c
}

func (c *Compressable) Delete(format string) *Compressable {
	if ok := inArray(format, allowedCompressFormats); !ok {
		if !c.FaultTolerant {
			panic(ErrCompressFormatNotSupported)
		}
		c.Errors = append(c.Errors, ErrCompressFormatNotSupported.Error())
		return c
	}

	if c.Formats[format] != nil {
		delete(c.Formats, format)
	} else {
		if !c.FaultTolerant {
			panic(ErrCompressFormatNotSupported)
		}
		c.Errors = append(c.Errors, ErrCompressionFormatNotSet.Error())
	}
	return c
}
