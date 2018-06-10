package resource

import (
	emodel "github.com/sniperkit/colly/plugins/data/aggregate/web/service/elasticfeed/model"
	"github.com/sniperkit/colly/plugins/data/aggregate/web/service/service/stream"
)

type ResourceManager struct {
	engine emodel.Elasticfeed
}

func (this *ResourceManager) Init() {}

func (this *ResourceManager) GetStreamService() *stream.StreamService {
	return this.GetEngine().GetServiceManager().GetStreamService()
}

func (this *ResourceManager) GetEngine() emodel.Elasticfeed {
	return this.engine
}

func NewResourceManager(engine emodel.Elasticfeed) emodel.ResourceManager {
	return &ResourceManager{engine}
}
