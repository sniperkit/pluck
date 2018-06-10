package router

import (
	"github.com/feedlabs/feedify"
	"github.com/sniperkit/colly/plugins/data/aggregate/web/service/service/store/v1/controller"
)

func InitEntryWorkflows() {
	feedify.Router("/v1/application/:applicationId:string/feed/:feedId:int/workflow", &controller.WorkflowController{}, "get:GetList;post:Post")
	feedify.Router("/v1/application/:applicationId:string/feed/:feedId:int/workflow/:feedWorkflowId:int", &controller.WorkflowController{}, "get:Get;delete:Delete;put:Put")
}
