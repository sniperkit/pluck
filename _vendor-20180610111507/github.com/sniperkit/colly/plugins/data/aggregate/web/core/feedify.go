package feedify

import (
	// Golang packages
	"fmt"
	"strconv"

	// Beego framework packages
	"github.com/astaxie/beego"

	// feedify packages
	"github.com/sniperkit/colly/plugins/data/aggregate/web/core/config"
	_ "github.com/sniperkit/colly/plugins/data/aggregate/web/core/graph/adapter"
	_ "github.com/sniperkit/colly/plugins/data/aggregate/web/core/stream/adapter/message"
)

func GetConfigKey(key string) string {
	return config.GetConfigKey(key)
}

func Banner() {
	fmt.Printf("Starting app '%s' on port '%s'\n", config.GetConfigKey("appname"), config.GetConfigKey("feedify::port"))
}

func SetStaticPath(url string, path string) *beego.App {
	return beego.SetStaticPath(url, path)
}

func Error(v ...interface{}) {
	beego.Error(v...)
}

func Run() {
	Banner()

	beego.HttpPort, _ = strconv.Atoi(config.GetConfigKey("feedify::port"))
	beego.Run()
}
