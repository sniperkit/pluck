// readjson.go - JSON file reading utilities
// Copyright © 2017 C. L. Banning, All rights reserved.
// See LICENSE file for information.
// Utility function taken from tamgroup/rwjson package. ©2016 TAM Group, Inc.

package checkjson

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

// ReadJSONFile returns an array of the JSON objects in 'file'. The file can have
// comments outside of the JSON objects as well as comments embedded in the
// JSON objects if preceeded by the number, '#', symbol.
//	File "test.json":
//		This file contains some test data for ReadJSONFile ...
//		{
//		  "author": "B. Dylan",
//		  "title" : "Ballad of a Thin Man"  # one of my favorites
//		}
//
//	Code:
//		j, _ := ReadJSONFile("test.json")
//		fmt.Println(string(j[0])) // prints: {"author":"B. Dylan","title":"Ballad of a Thin Man"}
//
func ReadJSONFile(file string) ([][]byte, error) {
	fd, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("err, opening file %s: %s", file, err.Error())
	}
	defer fd.Close()

	fi, err := fd.Stat()
	if err != nil {
		return nil, fmt.Errorf("err, Stat for %s: %s", file, err.Error())
	}

	// consume the whole file
	content := make([]byte, fi.Size())
	if i, err := fd.Read(content); err != nil || int64(i) != fi.Size() {
		return nil, fmt.Errorf("err, reading %s: %s", file, err.Error())
	}

	buf := bytes.NewBuffer(content)
	a := make([][]byte, 0)
	n := 1
	for {
		b, err := getJSONObject(buf)
		if err != nil {
			return a, fmt.Errorf("object #: %d - %s", n, err.Error())
		}
		if len(b) == 0 {
			break
		}
		a = append(a, b)
		n++
	}
	return a, nil
}

// ReadJSONReader returns the next JSON object from an io.Reader; it returns io.EOF
// if the Reader terminates. (See ReadJSONFile for notes on handling of embedded comments in
// JSON object.)
func ReadJSONReader(r io.Reader) ([]byte, error) {
	// We need to verify that the io.Reader implements io.ByteReader -
	// bufio.Buffer does, but os.File doesn't.  If it doesn't we wrap
	// it in a local io.ByteReader/Reader.
	buf, ok := r.(io.ByteReader)
	if !ok {
		buf = myByteReader(r) // see code at EOF
	}
	return getJSONObject(buf)
}

/* For buf created by file reads, have to handle Ctrl-characters ... strip them out
   these are the ones that GO handles directly, while some are unlikely, just handle them all!
	\a   U+0007 alert or bell
	\b   U+0008 backspace
	\f   U+000C form feed
	\n   U+000A line feed or newline
	\r   U+000D carriage return
	\t   U+0009 horizontal tab
	\v   U+000b vertical tab
*/
func getJSONObject(buf io.ByteReader) ([]byte, error) {
	var braces bool
	var braceCnt int
	var literal bool
	var comment bool
	var err error
	var lastB byte

	result := make([]byte, 0)
	b := make([]byte, 1)

	for {
		b[0], err = buf.ReadByte()
		if err != nil {
			// the only error returned is io.EOF
			break
		}
		// see if we're outside a JSON object
		if !braces && b[0] != '{' {
			continue
		}
		// see if we're scanning a comment
		if comment && b[0] != '\n' {
			continue
		}
		switch b[0] {
		case '#':
			// rest of line is a comment?
			if !literal {
				comment = true
				continue
			}
		case '\n':
			if comment {
				comment = false
				continue
			}
		case '{':
			if !literal {
				braceCnt++
				if !braces {
					braces = true
				}
			}
		case '}':
			if !literal {
				braceCnt--
			}
		case '"':
			if !literal {
				literal = true
			} else if lastB != '\\' {
				literal = false
			}
		}
		if !literal {
			if i := bytes.IndexAny(b, " \a\b\f\n\r\t\v"); i >= 0 {
				continue
			}
		}
		result = append(result, b[0])
		if braceCnt == 0 {
			return result, nil
		}
		lastB = b[0]
	}
	if braceCnt != 0 {
		return result, fmt.Errorf("EOF with unmatched braces: %s", result)
	}

	return result, nil // io.EOF
}

// ================ local io.ByteReader wrapper for an io.Reader ...
// Source: unabasedly appropriated from github.com/clbanning/mxj.

type byteReader struct {
	r io.Reader
	b []byte
}

func myByteReader(r io.Reader) io.ByteReader {
	b := make([]byte, 1)
	return &byteReader{r, b}
}

// need for io.Reader - but we don't use it ...
func (b *byteReader) Read(p []byte) (int, error) {
	return b.r.Read(p)
}

func (b *byteReader) ReadByte() (byte, error) {
	_, err := b.r.Read(b.b)
	if len(b.b) > 0 {
		return b.b[0], err
	}
	var c byte
	return c, err
}
