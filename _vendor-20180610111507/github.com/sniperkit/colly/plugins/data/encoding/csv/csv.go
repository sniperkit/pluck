/*
Copyright (c) 2016,2017,2018, Maxim Konakov
All rights reserved.

Redistribution and use in source and binary forms, with or without modification,
are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice,
   this list of conditions and the following disclaimer.
2. Redistributions in binary form must reproduce the above copyright notice,
   this list of conditions and the following disclaimer in the documentation
   and/or other materials provided with the distribution.
3. Neither the name of the copyright holder nor the names of its contributors
   may be used to endorse or promote products derived from this software without
   specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED.
IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT,
INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY
OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING
NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE,
EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

// Package csvplus extends the standard Go encoding/csv package with fluent
// interface, lazy stream processing operations, indices and joins.
package csv

import (
	"bytes"
	"encoding/csv"
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"unsafe"
)

// compare the specified columns from the two rows
func equalRows(columns []string, r1, r2 Row) bool {
	for _, col := range columns {
		if r1[col] != r2[col] {
			return false
		}
	}

	return true
}

// check if all the column names from the specified list are unique
func allColumnsUnique(columns []string) bool {
	set := make(map[string]struct{}, len(columns))

	for _, col := range columns {
		if _, found := set[col]; found {
			return false
		}

		set[col] = struct{}{}
	}

	return true
}

// All is a predicate combinator that takes any number of other predicates and
// produces a new predicate which returns 'true' only if all the specified predicates
// return 'true' for the same input Row.
func All(funcs ...func(Row) bool) func(Row) bool {
	return func(row Row) bool {
		for _, pred := range funcs {
			if !pred(row) {
				return false
			}
		}

		return true
	}
}

// Any is a predicate combinator that takes any number of other predicates and
// produces a new predicate which returns 'true' if any the specified predicates
// returns 'true' for the same input Row.
func Any(funcs ...func(Row) bool) func(Row) bool {
	return func(row Row) bool {
		for _, pred := range funcs {
			if pred(row) {
				return true
			}
		}

		return false
	}
}

// Not produces a new predicate that reverts the return value from the given predicate.
func Not(pred func(Row) bool) func(Row) bool {
	return func(row Row) bool {
		return !pred(row)
	}
}

// Like produces a predicate that returns 'true' if its input Row matches all the corresponding
// values from the specified 'match' Row.
func Like(match Row) func(Row) bool {
	if len(match) == 0 {
		panic(errEmptyMatchFuncLike)
		// panic("Empty match function in Like() predicate")
	}

	return func(row Row) bool {
		for key, val := range match {
			if v, found := row[key]; !found || v != val {
				return false
			}
		}

		return true
	}
}
