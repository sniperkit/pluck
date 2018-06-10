package main

import (
	colly "github.com/sniperkit/colly/pkg"
	proxy "github.com/sniperkit/colly/pkg/proxy"
)

// collector - proxy
var (
	cp colly.ProxyFunc  // collector's default proxy function
	pl *proxy.ProxyList // collector's multi-protocol proxy object
)
