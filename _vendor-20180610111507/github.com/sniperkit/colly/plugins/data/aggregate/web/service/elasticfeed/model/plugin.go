package model

import (
	pmodel "github.com/sniperkit/colly/plugins/data/aggregate/web/service/plugin/model"
)

type PluginManager interface {
	LoadPipeline(name string) (pmodel.Pipeline, error)
}
