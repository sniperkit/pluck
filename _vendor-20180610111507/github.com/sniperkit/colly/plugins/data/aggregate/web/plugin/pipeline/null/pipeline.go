package null

import (
	"github.com/sniperkit/colly/plugins/data/aggregate/web/service/common"
	"github.com/sniperkit/colly/plugins/data/aggregate/web/service/workflow"
	//	"github.com/sniperkit/colly/plugins/data/aggregate/web/service/plugin/model"
)

type config struct {
	common.ElasticfeedConfig `mapstructure:",squash"`

	tpl *workflow.ConfigTemplate
}

type Pipeline struct {
	config config
}

func (p *Pipeline) Prepare(raws ...interface{}) ([]string, error) {
	return nil, nil
}

func (p *Pipeline) Run(data interface{}) (interface{}, error) {
	return nil, nil
}

func (p *Pipeline) Cancel() {
}
