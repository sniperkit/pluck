// Per: https://github.com/sniperkit/xutil/plugin/format/convert/mxj/issues/24
// Per: https://github.com/sniperkit/xutil/plugin/format/convert/mxj/issues/25

package main

import (
	"fmt"
	"io"
	"os"

	"github.com/sniperkit/colly/plugins/data/transform/mxj/x2j"
)

func main() {
	for {
		_, _, err := x2j.XmlReaderToJsonWriter(os.Stdin, os.Stdout)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}
