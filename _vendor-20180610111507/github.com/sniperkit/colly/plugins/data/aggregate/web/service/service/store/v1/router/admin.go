package router

import (
	"github.com/feedlabs/feedify"
	"github.com/sniperkit/colly/plugins/data/aggregate/web/service/service/store/v1/controller"
)

func InitAdminRouters() {
	feedify.Router("/v1/admin", &controller.AdminController{}, "get:GetList;post:Post")
	feedify.Router("/v1/admin/:adminId:string", &controller.AdminController{}, "get:Get;delete:Delete;put:Put")
}
