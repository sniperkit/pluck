package mxj

/*
import (
	"sort"
	"strconv"
	"fmt"

	json "github.com/sniperkit/xutil/plugin/format/json"
	csv "github.com/sniperkit/xutil/plugin/format/csv/concurrent-writer"
	"github.com/sniperkit/xutil/plugin/struct/csv"
)

func sortedKeys(m map[string]interface{}) []string {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func printRow(w *csv.Writer, keys []string, d map[string]interface{}) error {
	var record []string
	for _, k := range keys {
		switch f := d[k].(type) {
		case string:
			record = append(record, f)
		case float64:
			record = append(record, strconv.FormatFloat(f, 'f', -1, 64))
		case bool:
			if f {
				record = append(record, "true")
			} else {
				record = append(record, "false")
			}
		default:
			log.Fatalf("Unsupported type %T. Aborting.\n", f)
		}
	}
	return w.Write(record)
}

*/
