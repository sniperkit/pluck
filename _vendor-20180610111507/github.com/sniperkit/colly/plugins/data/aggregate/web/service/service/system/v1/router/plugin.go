package router

import (
	"github.com/feedlabs/feedify"
	"github.com/sniperkit/colly/plugins/data/aggregate/web/service/service/system/v1/controller"
)

func InitPluginRouters() {
	feedify.Router("/v1/system/plugin", &controller.PluginController{}, "get:GetList;post:Post")
	feedify.Router("/v1/system/plugin/:pluginId:string", &controller.PluginController{}, "get:Get;delete:Delete;put:Put")
	feedify.Router("/v1/system/plugin/:pluginId:string/upload", &controller.PluginController{}, "put:PutFile")
}
