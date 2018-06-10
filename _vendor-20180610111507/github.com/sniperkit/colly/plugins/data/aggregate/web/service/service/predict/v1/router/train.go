package router

import (
	"github.com/feedlabs/feedify"
	"github.com/sniperkit/colly/plugins/data/aggregate/web/service/service/predict/v1/controller"
)

func InitTrainRouters() {
	feedify.Router("/v1/predict/train", &controller.DefaultController{}, "get:Get")
}
