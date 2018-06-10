package service

import (
	"github.com/sniperkit/colly/plugins/data/aggregate/web/service/service/predict"
	"github.com/sniperkit/colly/plugins/data/aggregate/web/service/service/store"
	"github.com/sniperkit/colly/plugins/data/aggregate/web/service/service/stream"
	"github.com/sniperkit/colly/plugins/data/aggregate/web/service/service/system"

	"github.com/sniperkit/colly/plugins/data/aggregate/web/service/elasticfeed/model"
)

type Service struct{}

type ServiceManager struct {
	engine model.Elasticfeed

	store   *store.DbService
	stream  *stream.StreamService
	system  *system.SystemService
	predict *predict.PredictService
}

func (this *ServiceManager) Init() {
	this.system.Init()
	this.stream.Init()
	this.store.Init()
	this.predict.Init()
}

func (this *ServiceManager) GetDbService() *store.DbService {
	return this.store
}

func (this *ServiceManager) GetStreamService() *stream.StreamService {
	return this.stream
}

func (this *ServiceManager) GetSystemService() *system.SystemService {
	return this.system
}

func (this *ServiceManager) GetPredictService() *predict.PredictService {
	return this.predict
}

func NewServiceManager(engine model.Elasticfeed) *ServiceManager {

	store := store.NewDbService()
	stream := stream.NewStreamService()
	system := system.NewSystemService()
	predict := predict.NewPredictService()

	return &ServiceManager{engine, store, stream, system, predict}
}
