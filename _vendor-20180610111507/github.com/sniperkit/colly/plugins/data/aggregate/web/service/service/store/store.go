package store

import (
	"github.com/sniperkit/colly/plugins/data/aggregate/web/service/service/store/v1/controller"
	"github.com/sniperkit/colly/plugins/data/aggregate/web/service/service/store/v1/router"
)

type DbService struct{}

func (this *DbService) Init() {
	router.InitRouters()
	controller.InitService()
}

func NewDbService() *DbService {
	return &DbService{}
}
