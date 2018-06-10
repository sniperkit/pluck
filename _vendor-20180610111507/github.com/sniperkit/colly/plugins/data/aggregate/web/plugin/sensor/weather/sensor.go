package weather

import (
	"github.com/sniperkit/colly/plugins/data/aggregate/web/service/common"
	"github.com/sniperkit/colly/plugins/data/aggregate/web/service/workflow"
	//	"github.com/sniperkit/colly/plugins/data/aggregate/web/service/plugin/model"
)

type config struct {
	common.ElasticfeedConfig `mapstructure:",squash"`

	tpl *workflow.ConfigTemplate
}

type Sensor struct {
	config config
}

func (p *Sensor) Prepare(raws ...interface{}) ([]string, error) {
	return nil, nil
}

func (p *Sensor) Run(data interface{}) (interface{}, error) {
	return data, nil
}

func (p *Sensor) Cancel() {
}
