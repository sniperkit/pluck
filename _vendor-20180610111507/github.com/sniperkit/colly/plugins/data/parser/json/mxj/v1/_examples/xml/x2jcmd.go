// Per: https://github.com/sniperkit/colly/plugins/data/transform/mxj/issues/24
// Per: https://github.com/sniperkit/colly/plugins/data/transform/mxj/issues/25

package main

import (
	"fmt"
	"github.com/sniperkit/colly/plugins/data/transform/mxj/x2j"
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
