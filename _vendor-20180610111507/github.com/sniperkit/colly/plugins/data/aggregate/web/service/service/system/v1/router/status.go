package router

import (
	"github.com/feedlabs/feedify"
	"github.com/sniperkit/colly/plugins/data/aggregate/web/service/service/system/v1/controller"
)

func InitRouters() {
	feedify.Router("/v1/system/status", &controller.StatusController{}, "get:Get")
}
