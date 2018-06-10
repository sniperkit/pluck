package advanced

import (
	"errors"
	"fmt"
	"os"
)

var (
	errCurrentProxyUnset = errors.New("currentProxy is unset")
	errEmptyPoolProxies  = errors.New("Empty pool of proxies")
)

func fatalf(fmtStr string, args interface{}) {
	fmt.Fprintf(os.Stderr, fmtStr, args)
	os.Exit(-1)
}
