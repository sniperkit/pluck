// Per: https://github.com/sniperkit/xutil/plugin/format/convert/mxj/pkg/issues/24
// Per: https://github.com/sniperkit/xutil/plugin/format/convert/mxj/pkg/issues/25

package main

import (
	"fmt"
	"github.com/sniperkit/xutil/plugin/format/convert/mxj/pkg/x2j"
	"io"
	"os"
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
