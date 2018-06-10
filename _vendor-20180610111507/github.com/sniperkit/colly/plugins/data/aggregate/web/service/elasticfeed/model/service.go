package model

import (
	"github.com/sniperkit/colly/plugins/data/aggregate/web/service/service/stream"
)

type ServiceManager interface {
	GetStreamService() *stream.StreamService

	Init()
}
